package phpproxy

import (
	"encoding/binary"
	"fmt"
	"net"
	"time"
)

//phpproxy connection pool
type PHPProxyConnPool struct {
	path     string
	poolSize int
	pool     chan net.Conn
}

var (
	proxyPool *PHPProxyConnPool
)

type TemporaryError interface {
	Temporary() bool
}

func InitPHPProxyConPool(path string, poolsize int) error {

	proxyPool = &PHPProxyConnPool{
		path:     path,
		poolSize: poolsize,
		pool:     make(chan net.Conn, poolsize),
	}

	for index := 0; index < poolsize; index++ {
		conn, err := net.Dial("unix", path)
		if err != nil {
			return err
		}
		proxyPool.pool <- conn
		time.Sleep(10 * time.Millisecond)
	}
	/*go func() {
		err := <-errc
		Close()
		errc <- err
	}()
	*/
	return nil
}

func Close() {
	if proxyPool == nil || proxyPool.path == "" {
		return
	}
	proxyPool.path = ""
	for index := 0; index < proxyPool.poolSize; index++ {
		conn := <-proxyPool.pool
		conn.Close()
	}
	close(proxyPool.pool)
	proxyPool = nil
	fmt.Println("phpproxy pool closed")
}
func Call(msg []byte) (ret []byte, err error) {

	if proxyPool == nil || proxyPool.path == "" {
		return nil, fmt.Errorf("call phpproxy err!connection pool not init!")
	}

	//wait for a conn
	conn := <-proxyPool.pool

	//back to the pool
	defer func() {

		//if has error，create a new connection,this may cause a deadlock
		if err != nil {
			var newcoon net.Conn
			conn.Close()
			newcoon, err = net.Dial("unix", proxyPool.path)
			if err != nil {
				fmt.Println("phpproxy connection error!", err)
				return
			}
			conn = newcoon
		}

		proxyPool.pool <- conn
	}()

	//write message length
	var length uint32 = uint32(len(msg))
	err = binary.Write(conn, binary.BigEndian, length)
	if nil != err {
		return nil, fmt.Errorf("phpproxy SendMessage err!write response length err,%v", err)
	}

	//write message flage
	var flage uint32 = 0
	err = binary.Write(conn, binary.BigEndian, flage)
	if nil != err {
		return nil, fmt.Errorf("phpproxy SendMessage err!write response flage err,%v", err)
	}

	//write response
	err = binary.Write(conn, binary.BigEndian, msg)
	if err != nil {
		return nil, fmt.Errorf("phpproxy SendMessage err!send response err,%v", err)
	}

	//read message length
	var readlength uint32
	err = binary.Read(conn, binary.BigEndian, &readlength)
	if nil != err {
		return nil, fmt.Errorf("phpproxy read response err!failed when read request length, %v", err)
	}

	//read flage
	var readflage uint32
	err = binary.Read(conn, binary.BigEndian, &readflage)
	if nil != err {
		return nil, fmt.Errorf("phpproxy read response err!failed when read flage, %v", err)
	}

	//read message body
	ret = make([]byte, readlength)
	err = binary.Read(conn, binary.BigEndian, ret)
	if nil != err {
		return nil, fmt.Errorf("phpproxy read response err!failed when read request data, %v", err)
	}

	return
}
