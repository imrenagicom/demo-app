package util

import (
	"github.com/jmoiron/sqlx"
	"github.com/redis/go-redis/v9"
)

type Clients struct {
	DB    *sqlx.DB
	Redis redis.UniversalClient
}
