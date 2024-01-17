package postgres

import (
	"github.com/imrenagicom/demo-app/internal/config"

	"github.com/jmoiron/sqlx"
)

func NewSQLx(c config.SQL) *sqlx.DB {
	db, err := sqlx.Open("postgres", c.DataSourceName())
	if err != nil {
		panic(err)
	}
	db.SetMaxOpenConns(c.MaxOpenConn)
	db.SetMaxIdleConns(c.MaxIdleConn)
	return db
}
