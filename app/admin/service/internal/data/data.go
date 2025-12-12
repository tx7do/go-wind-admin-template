package data

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/go-redis/v9"

	conf "github.com/tx7do/kratos-bootstrap/api/gen/go/conf/v1"
	redisClient "github.com/tx7do/kratos-bootstrap/cache/redis"
)

// Data .
type Data struct {
	log *log.Helper

	rdb *redis.Client
}

// NewData .
func NewData(
	logger log.Logger,
	rdb *redis.Client,
) (*Data, func(), error) {
	l := log.NewHelper(log.With(logger, "module", "data/admin-service"))

	d := &Data{
		log: l,

		rdb: rdb,
	}

	return d, func() {
		l.Info("closing the data resources")

		if d.rdb != nil {
			if err := d.rdb.Close(); err != nil {
				l.Error(err)
			}
		}
	}, nil
}

// NewRedisClient 创建Redis客户端
func NewRedisClient(cfg *conf.Bootstrap, logger log.Logger) *redis.Client {
	l := log.NewHelper(log.With(logger, "module", "redis/data/admin-service"))
	return redisClient.NewClient(cfg.Data, l)
}
