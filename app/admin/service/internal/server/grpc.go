package server

import (
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/transport/grpc"

	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"github.com/tx7do/kratos-bootstrap/rpc"

	"github.com/tx7do/go-wind-admin-template/app/admin/service/internal/service"

	helloworldV1 "github.com/tx7do/go-wind-admin-template/api/gen/go/helloworld/service/v1"
)

// NewGrpcServer new a gRPC server.
func NewGrpcServer(
	ctx *bootstrap.Context,
	greeterService *service.GreeterService,
) *grpc.Server {
	cfg := ctx.GetConfig()

	if cfg == nil || cfg.Server == nil || cfg.Server.Grpc == nil {
		return nil
	}

	srv, err := rpc.CreateGrpcServer(cfg, logging.Server(ctx.GetLogger()))
	if err != nil {
		panic(err)
	}

	helloworldV1.RegisterGreeterServer(srv, greeterService)

	return srv
}
