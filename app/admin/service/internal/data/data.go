package data

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/go-redis/v9"

	"github.com/tx7do/kratos-bootstrap/bootstrap"
)

// Data .
type Data struct {
	log *log.Helper

	rdb *redis.Client
}

// NewData .
func NewData(
	ctx *bootstrap.Context,
	rdb *redis.Client,
) (*Data, func(), error) {
	d := &Data{
		log: ctx.NewLoggerHelper("data/admin-service"),

		rdb: rdb,
	}

	return d, func() {
		d.log.Info("closing the data resources")

		if d.rdb != nil {
			if err := d.rdb.Close(); err != nil {
				d.log.Error(err)
			}
		}
	}, nil
}
