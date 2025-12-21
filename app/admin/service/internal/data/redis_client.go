package data

import (
	"github.com/redis/go-redis/v9"

	"github.com/tx7do/kratos-bootstrap/bootstrap"
	redisClient "github.com/tx7do/kratos-bootstrap/cache/redis"
)

// NewRedisClient 创建Redis客户端
func NewRedisClient(ctx *bootstrap.Context) *redis.Client {
	cfg := ctx.GetConfig()
	if cfg == nil {
		return nil
	}
	return redisClient.NewClient(cfg.Data, ctx.NewLoggerHelper("redis/data/admin-service"))
}
