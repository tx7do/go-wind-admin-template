package server

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/logging"
	"github.com/go-kratos/kratos/v2/transport/http"

	swaggerUI "github.com/tx7do/kratos-swagger-ui"

	"github.com/tx7do/kratos-bootstrap/bootstrap"
	"github.com/tx7do/kratos-bootstrap/rpc"

	"github.com/tx7do/go-wind-admin-template/app/admin/service/cmd/server/assets"
	"github.com/tx7do/go-wind-admin-template/app/admin/service/internal/service"

	adminV1 "github.com/tx7do/go-wind-admin-template/api/gen/go/admin/service/v1"
)

// NewMiddleware 创建中间件
func newRestMiddleware(
	logger log.Logger,
) []middleware.Middleware {
	var ms []middleware.Middleware
	ms = append(ms, logging.Server(logger))

	return ms
}

// NewRestServer new an REST server.
func NewRestServer(
	ctx *bootstrap.Context,
	greeterService *service.GreeterService,
) *http.Server {
	cfg := ctx.GetConfig()

	if cfg == nil || cfg.Server == nil || cfg.Server.Rest == nil {
		return nil
	}

	srv, err := rpc.CreateRestServer(cfg,
		newRestMiddleware(ctx.GetLogger())...,
	)
	if err != nil {
		panic(err)
	}

	if cfg.GetServer().GetRest().GetEnableSwagger() {
		swaggerUI.RegisterSwaggerUIServerWithOption(
			srv,
			swaggerUI.WithTitle("GoWind Admin"),
			swaggerUI.WithMemoryData(assets.OpenApiData, "yaml"),
		)
	}

	adminV1.RegisterGreeterHTTPServer(srv, greeterService)

	return srv
}
