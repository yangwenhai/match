package main

import (
	"fmt"
	"net/http"
	"net/http/pprof"
	"os"
	"os/signal"
	"syscall"

	stdopentracing "github.com/opentracing/opentracing-go"
	zipkin "github.com/openzipkin/zipkin-go-opentracing"
	stdprometheus "github.com/prometheus/client_golang/prometheus"

	logger "kitfw/commom/log"
	"kitfw/commom/phpproxy"
	"kitfw/commom/registry"
	"kitfw/commom/store"
	"kitfw/match/envconf"
	rpcgrpc "kitfw/match/rpc/grpc"
	rpctcp "kitfw/match/rpc/tcp"
	"kitfw/match/service/matchservice"

	"github.com/go-kit/kit/metrics"
	"github.com/go-kit/kit/metrics/prometheus"
)

func main() {
	//log
	logger.SetDefaultLogLevel(envconf.EnvCfg.LOG_LEVEL)
	logger.Info("msg", fmt.Sprintf("hello %s", envconf.EnvCfg.SERVER_NAME))
	defer logger.Info("msg", fmt.Sprintf("goodbye %s", envconf.EnvCfg.SERVER_NAME))

	// Mechanical domain.
	errc := make(chan error)

	// Interrupt handler.
	go func() {
		c := make(chan os.Signal, 1)
		signal.Notify(c, syscall.SIGINT, syscall.SIGTERM)
		errc <- fmt.Errorf("%s", <-c)
	}()

	// Debug listener.
	go func() {
		m := http.NewServeMux()
		m.Handle("/debug/pprof/", http.HandlerFunc(pprof.Index))
		m.Handle("/debug/pprof/cmdline", http.HandlerFunc(pprof.Cmdline))
		m.Handle("/debug/pprof/profile", http.HandlerFunc(pprof.Profile))
		m.Handle("/debug/pprof/symbol", http.HandlerFunc(pprof.Symbol))
		m.Handle("/debug/pprof/trace", http.HandlerFunc(pprof.Trace))
		m.Handle("/metrics", stdprometheus.Handler())
		addr := fmt.Sprintf("%s:%d", envconf.EnvCfg.HOST, envconf.EnvCfg.DEBUG_PORT)
		logger.Info("transport", "debug", "debugAddr", addr)
		errc <- http.ListenAndServe(addr, m)
	}()

	// Metrics domain.
	fieldKeys := []string{"method", "protoid", "logid", "error"}
	var requestCount metrics.Counter
	{
		// Business level metrics.
		requestCount = prometheus.NewCounterFrom(stdprometheus.CounterOpts{
			Namespace: envconf.EnvCfg.SERVER_NAME,
			Name:      "request_count",
			Help:      "Number of requests received.",
		}, fieldKeys)
	}
	var duration metrics.Histogram
	{
		// Transport level metrics.
		duration = prometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: envconf.EnvCfg.SERVER_NAME,
			Name:      "request_duration_ns",
			Help:      "Request duration in nanoseconds.",
		}, fieldKeys)
	}
	var endpointDuration metrics.Histogram
	{
		// Transport level metrics.
		endpointDuration = prometheus.NewSummaryFrom(stdprometheus.SummaryOpts{
			Namespace: envconf.EnvCfg.SERVER_NAME,
			Name:      "endpoint_request_duration_ns",
			Help:      "endpoint request duration in nanoseconds.",
		}, []string{"method", "logid", "success"})
	}

	// Tracing domain.
	var tracer stdopentracing.Tracer
	{
		if envconf.EnvCfg.ZIPKIN_ADDR != "" {
			logger.Info("tracer", "Zipkin", "zipkinAddr", envconf.EnvCfg.ZIPKIN_ADDR)
			// collector, err := zipkin.NewKafkaCollector(
			// 	strings.Split(*zipkinAddr, ","),
			// 	zipkin.KafkaLogger(logger),
			// )
			collector, err := zipkin.NewHTTPCollector(envconf.EnvCfg.ZIPKIN_ADDR)
			if err != nil {
				logger.Error("tracer", "Zipkin", "err", err)
				os.Exit(1)
			}
			tracer, err = zipkin.NewTracer(
				zipkin.NewRecorder(collector, false, "HOST:PORT", envconf.EnvCfg.NODE_NAME),
				zipkin.ClientServerSameSpan(true),
			)
			if err != nil {
				logger.Error("tracer", "Zipkin", "err", err)
				os.Exit(1)
			}
		} else {
			logger.Info("tracer", "none")
			tracer = stdopentracing.GlobalTracer() // no-op
		}
	}

	//grpc for capnp
	g := rpcgrpc.NewGrpcService(tracer, requestCount, duration, endpointDuration)
	go g.RunGrpcServer(errc)

	//tcp for amf
	tc := rpctcp.NewTcpService(tracer, requestCount, duration, endpointDuration)
	go tc.RunTcpServer(errc)

	// zookeeper registry
	if len(envconf.EnvCfg.ZK_ADDRS) > 0 {
		r, err := registry.NewZkService(envconf.EnvCfg.ZK_ADDRS)
		if err != nil {
			logger.Error("error", fmt.Sprintf("unexpected error creating zookeeper client:%v", err))
			os.Exit(1)
		}
		err = r.Register(envconf.EnvCfg.ZK_REGISTRY_PATH, envconf.EnvCfg.NODE_NAME, envconf.EnvCfg.HOST, envconf.EnvCfg.TCP_PORT)
		if err != nil {
			logger.Error("error", fmt.Sprintf("unexpected error creating zookeeper register:%v", err))
			os.Exit(1)
		}
		logger.Info("zkAddr", envconf.EnvCfg.ZK_ADDRS[0])
	}

	//init phpproxy connetcion pool
	if envconf.EnvCfg.PHPPROXY_SOCKT_PATH != "" && envconf.EnvCfg.PHPPROXY_POOL_SIZE > 0 {
		if err := phpproxy.InitPHPProxyConPool(envconf.EnvCfg.PHPPROXY_SOCKT_PATH, envconf.EnvCfg.PHPPROXY_POOL_SIZE); err != nil {
			logger.Error("error", fmt.Sprintf("InitPHPProxyConPool error,%v", err))
			os.Exit(1)
		}
	} else {
		logger.Warn("warn", "phpproxy sockt path empty")
	}

	//init redis
	if len(envconf.EnvCfg.REDIS_ADDRS) > 0 && envconf.EnvCfg.REDIS_POOL_SIZE > 0 {
		if err := store.InitRedisConnPool(envconf.EnvCfg.REDIS_ADDRS, envconf.EnvCfg.REDIS_POOL_SIZE); err != nil {
			logger.Error("error", fmt.Sprintf("InitRedisConnPool error,%v", err))
			os.Exit(1)
		}
	}

	//start match service
	if err := matchservice.StartMatchService(); err != nil {
		logger.Error("error", fmt.Sprintf("StartMatchService error,%v", err))
		os.Exit(1)
	}

	// //watch lcserver
	// if *zkAddr != "" {
	// 	if err := registry.InitWatchLcserver([]string{*zkAddr}, "/pirate/lcserver"); err != nil {
	// 		logger.Error("error", fmt.Sprintf("unexpected error creating zookeeper client:%v", err))
	// 		os.Exit(1)
	// 	}
	// }

	// Run!
	logger.Error("exit", <-errc)
}
