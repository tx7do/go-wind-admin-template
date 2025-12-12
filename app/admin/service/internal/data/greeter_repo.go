package data

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	helloworldV1 "github.com/tx7do/go-wind-admin-template/api/gen/go/helloworld/service/v1"
)

type GreeterRepo struct {
	data *Data
	log  *log.Helper
}

// NewGreeterRepo .
func NewGreeterRepo(data *Data, logger log.Logger) *GreeterRepo {
	return &GreeterRepo{
		data: data,
		log:  log.NewHelper(logger),
	}
}

func (r *GreeterRepo) SayHello(_ context.Context, req *helloworldV1.HelloRequest) (*helloworldV1.HelloReply, error) {
	return &helloworldV1.HelloReply{Message: "Hello " + req.GetName()}, nil
}
