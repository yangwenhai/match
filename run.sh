#!/bin/bash

#通过web访问服务运行时信息的端口
export DEBUG_PORT=8101
#供phpproxy访问的端口
export TCP_PORT=8102
#grpc访问的端口
export GRPC_PORT=8103
#本机ip
export HOST=192.168.88.59
#本服务的名字，在metrics里会用到
export SERVER_NAME=match
#节点名字，zookeeper里会用到(比如 battle-0  battle-1  battle-2)
export NODE_NAME=match-00
#日志级别（1 debug 2 info 3 warning 4 error）
export LOG_LEVEL=1
#zookeeper节点前缀
export ZK_PREFIX="/pirate"
#本服务注册zookeeper时的路径
export ZK_REGISTRY_PATH="/pirate/match"
#zookeeper地址(多个地址用分号隔开)
export ZK_ADDRS="115.182.251.193:8181"
#zipkin地址
export ZIPKIN_ADDR="http://192.168.88.59:9411/api/v1/spans"
#phpproxy的socket文件地址
export PHPPROXY_SOCKT_PATH="" #"/home/pirate/phpproxy/var/phpproxy.sock"
#phpproxy连接池大小
export PHPPROXY_POOL_SIZE=30
#redis地址
export REDIS_ADDRS="192.168.1.36:8379"
#redis连接池大小
export REDIS_POOL_SIZE=30

#匹配超时(秒)
export MATCH_TIME_OUT=10 
#匹配时的积分下限
export MATCH_ONE_MIN_SCORE=0
#匹配时的积分上限
export MATCH_ONE_MAX_SCORE=3000
#匹配时两个玩家的积分相差多少可以匹配
export MATCH_ONE_SCORE_LIMIT=100
#将上面的MATCH_ONE_MIN_SCORE ~ MATCH_ONE_MAX_SCORE 分成多个区间，每个区间多少积分
export MATCH_ONE_SCORE_INTERVAL=30
#建立多少个redis连接进行匹配
export MATCH_ONE_NET_THREAD=5
    
./match