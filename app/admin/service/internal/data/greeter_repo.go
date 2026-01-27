package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	helloworldV1 "github.com/tx7do/go-wind-admin-template/api/gen/go/helloworld/service/v1"
)

type GreeterRepo struct {
	log *log.Helper
}

// NewGreeterRepo .
func NewGreeterRepo(ctx *bootstrap.Context) *GreeterRepo {
	return &GreeterRepo{
		log: ctx.NewLoggerHelper("greeter/repo/admin-service"),
	}
}

func (r *GreeterRepo) SayHello(_ context.Context, req *helloworldV1.HelloRequest) (*helloworldV1.HelloReply, error) {
	return &helloworldV1.HelloReply{Message: "Hello " + req.GetName()}, nil
}
