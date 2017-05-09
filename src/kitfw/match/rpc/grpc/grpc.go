package grpc

import (
	"context"
	"fmt"
	"kitfw/commom/pb"
	ep "kitfw/match/rpc/endpoint"
	kitservice "kitfw/match/service"
	"net"

	stdopentracing "github.com/opentracing/opentracing-go"

	"kitfw/match/envconf"

	logger "kitfw/commom/log"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/metrics"
	"google.golang.org/grpc"
)

type grpcService struct {
	endpoint endpoint.Endpoint
	tracer   stdopentracing.Tracer
}

func NewGrpcService(tracer stdopentracing.Tracer, requestCount metrics.Counter, duration metrics.Histogram, endpointDuration metrics.Histogram) *grpcService {
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
	return &grpcService{endpoint: requestEndpoint, tracer: tracer}
}

func (g *grpcService) RunGrpcServer(errc chan error) {

	grpcAddr := fmt.Sprintf("%s:%d", envconf.EnvCfg.HOST, envconf.EnvCfg.GRPC_PORT)
	ln, err := net.Listen("tcp", grpcAddr)
	if err != nil {
		errc <- err
		return
	}
	ctx := context.Background()
	srv := MakeGRPCServer(ctx, g.endpoint, g.tracer, logger.GetDefaultLogger())
	s := grpc.NewServer()
	pb.RegisterKitfwServer(s, srv)
	logger.Info("transport", "gRPC", "grpcAddr", grpcAddr)
	errc <- s.Serve(ln)
}
