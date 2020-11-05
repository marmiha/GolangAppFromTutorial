package postgres

import (
	"github.com/go-pg/pg"
)

func New(opts *pg.Options) *pg.DB {
	db := pg.Connect(opts)
	return db
}

