package platform

import (
	"time"

	_ "github.com/go-sql-driver/mysql" // Load MySQL driver
	"github.com/jmoiron/sqlx"
)

func NewMySql(uri string) (*sqlx.DB, error) {
	db, err := sqlx.Connect("mysql", uri)

	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return db, nil
}
