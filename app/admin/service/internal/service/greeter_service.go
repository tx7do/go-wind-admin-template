package service

import (
	"context"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/tx7do/kratos-bootstrap/bootstrap"

	"github.com/tx7do/go-wind-admin-template/app/admin/service/internal/data"

	helloworldV1 "github.com/tx7do/go-wind-admin-template/api/gen/go/helloworld/service/v1"
)

// GreeterService is a greeter service.
type GreeterService struct {
	helloworldV1.UnimplementedGreeterServer

	log *log.Helper

	repo *data.GreeterRepo
}

// NewGreeterService new a greeter service.
func NewGreeterService(ctx *bootstrap.Context, repo *data.GreeterRepo) *GreeterService {
	return &GreeterService{
		repo: repo,
		log:  ctx.NewLoggerHelper("greeter/service/admin-service"),
	}
}

// SayHello implements helloworld.GreeterServer.
func (s *GreeterService) SayHello(ctx context.Context, req *helloworldV1.HelloRequest) (*helloworldV1.HelloReply, error) {
	resp, err := s.repo.SayHello(ctx, req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}
