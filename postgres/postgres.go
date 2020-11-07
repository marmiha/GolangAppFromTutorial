package postgres

import (
	"github.com/go-pg/pg/v10"
	"github.com/go-pg/pg/v10/orm"
	"todo/domain"
)

func New(opts *pg.Options) *pg.DB {
	db := pg.Connect(opts)
	return db
}

func CreateSchema(db *pg.DB, options *orm.CreateTableOptions) error {
	// Our database models/structs.
	models := []interface{} {
		(*domain.User)(nil),	// Add multiple in the list.
	}

	// Create each of the models.
	for _, model := range models {
		err := db.Model(model).CreateTable(options)
		if err != nil {
			return err
		}
	}
	return nil
}
