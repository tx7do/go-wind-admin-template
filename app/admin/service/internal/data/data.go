package data

import (
	"github.com/go-kratos/kratos/v2/log"
	"github.com/redis/go-redis/v9"
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
