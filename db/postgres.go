package db

import "github.com/go-pg/pg/v10"

func New(opts *pg.Options) *pg.DB {
	db := pg.Connect(opts)
	return db
}
