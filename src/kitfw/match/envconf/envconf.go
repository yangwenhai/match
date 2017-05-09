package envconf

import (
	"fmt"
	"os"

	"github.com/joeshaw/envdecode"
)

var EnvCfg EnvConfig
var EnvMatchCfg MachEvnConfig

func init() {
	err := envdecode.Decode(&EnvCfg)
	if err != nil {
		fmt.Println("envdecode err!", err)
		os.Exit(-1)
	}
	err = envdecode.Decode(&EnvMatchCfg)
	if err != nil {
		fmt.Println("envdecode err!", err)
		os.Exit(-1)
	}
}

type EnvConfig struct {
	DEBUG_PORT          uint16   `env:"DEBUG_PORT,required"`
	TCP_PORT            uint16   `env:"TCP_PORT,required"`
	GRPC_PORT           uint16   `env:"GRPC_PORT, required"`
	HOST                string   `env:"HOST, required"`
	SERVER_NAME         string   `env:"SERVER_NAME, required"`
	NODE_NAME           string   `env:"NODE_NAME, required"`
	LOG_LEVEL           int      `env:"LOG_LEVEL, default=1"`
	ZK_PREFIX           string   `env:"ZK_PREFIX"`
	ZK_REGISTRY_PATH    string   `env:"ZK_REGISTRY_PATH"`
	ZK_ADDRS            []string `env:"ZK_ADDRS"`
	ZIPKIN_ADDR         string   `env:"ZIPKIN_ADDR"`
	PHPPROXY_SOCKT_PATH string   `env:"PHPPROXY_SOCKT_PATH"`
	PHPPROXY_POOL_SIZE  int      `env:"PHPPROXY_POOL_SIZE"`
	REDIS_ADDRS         []string `env:"REDIS_ADDRS"`
	REDIS_POOL_SIZE     int      `env:"REDIS_POOL_SIZE"`
}

type MachEvnConfig struct {
	MATCH_TIME_OUT           int `env:"MATCH_TIME_OUT,required"`
	MATCH_ONE_MIN_SCORE      int `env:"MATCH_ONE_MIN_SCORE,required"`
	MATCH_ONE_MAX_SCORE      int `env:"MATCH_ONE_MAX_SCORE,required"`
	MATCH_ONE_SCORE_LIMIT    int `env:"MATCH_ONE_SCORE_LIMIT,required"`
	MATCH_ONE_SCORE_INTERVAL int `env:"MATCH_ONE_SCORE_INTERVAL,required"`
	MATCH_ONE_NET_THREAD     int `env:"MATCH_ONE_NET_THREAD,required"`
}
