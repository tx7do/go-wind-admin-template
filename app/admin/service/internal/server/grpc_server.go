package server

import (
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/transport/grpc"

	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"github.com/tx7do/kratos-bootstrap/rpc"

	"github.com/tx7do/go-wind-admin-template/app/admin/service/internal/service"

	helloworldV1 "github.com/tx7do/go-wind-admin-template/api/gen/go/helloworld/service/v1"
)

type GrpcMiddlewares []middleware.Middleware

func NewGrpcMiddleware(ctx *bootstrap.Context) GrpcMiddlewares {
	var ms []middleware.Middleware
	ms = append(ms, logging.Server(ctx.GetLogger()))
	return ms
}

// NewGrpcServer new a gRPC server.
func NewGrpcServer(
	ctx *bootstrap.Context,

	middlewares GrpcMiddlewares,

	greeterService *service.GreeterService,
	// register:param ── 新模块服务形参在此行后注册(make register 工具锚点,勿删)
) *grpc.Server {
	cfg := ctx.GetConfig()

	if cfg == nil || cfg.Server == nil || cfg.Server.Grpc == nil {
		return nil
	}

	srv, err := rpc.CreateGrpcServer(cfg, middlewares...)
	if err != nil {
		panic(err)
	}

	helloworldV1.RegisterGreeterServer(srv, greeterService)

	// register:route ── 新模块路由在此行后注册(make register 工具锚点,勿删)

	return srv
}
