package store

import (
	"fmt"
	"sync"
	"time"

	logger "kitfw/commom/log"

	"github.com/garyburd/redigo/redis"
	"stathat.com/c/consistent"
)

type RedisConnMaps struct {
	addrs    []string
	poolSize int
	connMaps map[string]*RedisConnPool
	consist  *consistent.Consistent
	sync.RWMutex
}

type RedisConnPool struct {
	addr     string
	connPool chan redis.Conn
}

var (
	MAX_POOL_SIZE  = 30
	redisConMaps   *RedisConnMaps
	connectTimeout = 3 * time.Second
	readTimeout    = 2 * time.Second
	writeTimeout   = 2 * time.Second
)

type TemporaryError interface {
	Temporary() bool
}

func InitRedisConnPool(addrs []string, poolsize int) error {

	poolMaps := &RedisConnMaps{
		addrs:    addrs,
		poolSize: poolsize,
		connMaps: make(map[string]*RedisConnPool),
		consist:  consistent.New(),
	}

	for _, addr := range addrs {
		pool, err := CreateRedisConPools(addr, poolsize)
		if err != nil {
			return err
		}
		poolMaps.connMaps[addr] = pool
		poolMaps.consist.Add(addr)
		logger.Info("redis-pool", addr)
	}
	redisConMaps = poolMaps

	go redisConMaps.checkRedisConnection()
	return nil
}

func CreateRedisConPools(addr string, poolsize int) (*RedisConnPool, error) {

	if addr == "" {
		return nil, fmt.Errorf("empty redis addr")
	}

	redisPool := &RedisConnPool{
		addr:     addr,
		connPool: make(chan redis.Conn, poolsize),
	}

	for index := 0; index < poolsize; index++ {
		conn, err := redis.Dial("tcp", addr)
		if err != nil {
			return nil, err
		}
		redisPool.connPool <- conn
		time.Sleep(5 * time.Millisecond)
	}

	return redisPool, nil
}

func RedisConnPoolOk() bool {
	if redisConMaps == nil || len(redisConMaps.addrs) == 0 {
		return false
	}
	return true
}
func Do(key string, command string, args ...interface{}) (reply interface{}, err error) {

	if redisConMaps == nil || len(redisConMaps.addrs) == 0 {
		return nil, fmt.Errorf("CallRedis error!redis connection not init")
	}
	pool, err := redisConMaps.getConnectByKey(key)
	if err != nil {
		return nil, err
	}
	conn := <-pool.connPool
	defer func() {
		pool.connPool <- conn
	}()
	return conn.Do(command, args...)
}

func DoScript(key string, script string, cmd ...interface{}) (interface{}, error) {
	if redisConMaps == nil || len(redisConMaps.addrs) == 0 {
		return nil, fmt.Errorf("CallRedis error!redis connection not init")
	}
	pool, err := redisConMaps.getConnectByKey(key)
	if err != nil {
		return nil, err
	}
	conn := <-pool.connPool
	defer func() {
		pool.connPool <- conn
	}()

	var s = redis.NewScript(len(cmd), script)
	return s.Do(conn, cmd...)
}

func (p *RedisConnMaps) getConnectByKey(key string) (*RedisConnPool, error) {
	p.RLock()
	defer p.RUnlock()
	// hash
	addr, err := p.consist.Get(key)
	if err != nil {
		return nil, err
	}
	// get pool
	pool, ok := p.connMaps[addr]
	if !ok {
		return nil, fmt.Errorf("getConnectByKey nil!key:%s addr:%s", key, addr)
	}
	return pool, nil
}

func (p *RedisConnMaps) checkRedisConnection() error {

	checkMaps := make(map[string]redis.Conn)

	for {

		time.Sleep(10 * time.Second)

		for _, addr := range p.addrs {
			if _, ok := checkMaps[addr]; !ok {
				conn, err := redis.DialTimeout("tcp", addr, connectTimeout, readTimeout, writeTimeout)
				if err != nil {
					continue
				}
				checkMaps[addr] = conn // checkMaps can only be called in this gourouting,so there is no need to lock
			}
		}

		for addr, conn := range checkMaps {
			_, err := conn.Do("PING")
			if err != nil {
				logger.Error("redis-conn-lost", addr, "err", err)
				delete(checkMaps, addr)
				p.Lock()
				delete(p.connMaps, addr)
				p.consist.Remove(addr)
				p.Unlock()
				continue
			}
			if err == nil {
				if _, ok := p.connMaps[addr]; !ok {
					if pool, err := CreateRedisConPools(addr, p.poolSize); err == nil {
						p.Lock()
						p.connMaps[addr] = pool
						p.consist.Add(addr)
						p.Unlock()
						logger.Error("redis-pool-add", addr)
					} else {
						logger.Error("addr", addr, "err", err)
					}
				}
			}
		}

	}
}
