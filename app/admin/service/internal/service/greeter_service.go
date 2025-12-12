package service

import (
	"context"

	"github.com/tx7do/go-wind-admin-template/app/admin/service/internal/data"

	helloworldV1 "github.com/tx7do/go-wind-admin-template/api/gen/go/helloworld/service/v1"
)

// GreeterService is a greeter service.
type GreeterService struct {
	helloworldV1.UnimplementedGreeterServer

	repo *data.GreeterRepo
}

// NewGreeterService new a greeter service.
func NewGreeterService(repo *data.GreeterRepo) *GreeterService {
	return &GreeterService{
		repo: repo,
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
