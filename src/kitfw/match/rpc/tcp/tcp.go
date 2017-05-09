package tcp

import (
	"fmt"
	"kitfw/match/envconf"
	ep "kitfw/match/rpc/endpoint"
	kitservice "kitfw/match/service"
	"net"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
	stdopentracing "github.com/opentracing/opentracing-go"
)

type tcpService struct {
	endpoint endpoint.Endpoint
}

func NewTcpService(tracer stdopentracing.Tracer, requestCount metrics.Counter, duration metrics.Histogram, endpointDuration metrics.Histogram) *tcpService {
	// Business domain.
	var service kitservice.Service
	{
		service = kitservice.NewBasicService()
		service = kitservice.ServiceLoggingMiddleware()(service)
		service = kitservice.ServiceInstrumentingMiddleware(requestCount, duration)(service)
	}

	// Endpoint domain.
	var requestEndpoint endpoint.Endpoint
	{
		requestEndpoint = ep.MakeProcessEndpoint(service, tracer)
		requestEndpoint = ep.EndpointLoggingMiddleware(tracer)(requestEndpoint)
		requestEndpoint = ep.EndpointInstrumentingMiddleware(endpointDuration)(requestEndpoint)
		requestEndpoint = ep.TraceInternalService(tracer, "Process")(requestEndpoint)
	}
	return &tcpService{endpoint: requestEndpoint}
}

func (g *tcpService) RunTcpServer(errc chan error) {

	grpcAddr := fmt.Sprintf("%s:%d", envconf.EnvCfg.HOST, envconf.EnvCfg.TCP_PORT)
	ln, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		errc <- err
		return
	}
	_ = ln
	errc <- RunServer(ln, g.endpoint)
}
