package platform

import (
	"fmt"
	"time"

	_ "github.com/go-sql-driver/mysql" // Load MySQL driver
	"github.com/jmoiron/sqlx"
)

type MySql struct {
	*sqlx.DB
}

func NewMySql(uri string) (*MySql, error) {
	db, err := sqlx.Connect("mysql", uri)

	if err != nil {
		return nil, err
	}

	db.SetConnMaxLifetime(5 * time.Minute)
	db.SetMaxOpenConns(10)
	db.SetMaxIdleConns(10)

	return &MySql{db}, nil
}

func (db *MySql) ExecuteTx(fn func(*sqlx.Tx) error) (retErr error) {
	tx, err := db.Beginx()
	if err != nil {
		return fmt.Errorf("cannot begin transaction: %w", err)
	}

	defer func() {
		if p := recover(); p != nil {
			if err := tx.Rollback(); err != nil {
				retErr = fmt.Errorf("cannot rollback transaction from panic: %w", err)
				return
			}
			panic(p)
		} else if err != nil {
			if err := tx.Rollback(); err != nil {
				retErr = fmt.Errorf("cannot rollback transaction from error: %w", err)
				return
			}
			retErr = err
		} else {
			if err := tx.Commit(); err != nil {
				retErr = fmt.Errorf("cannot commit transaction: %w", err)
				return
			}
		}
	}()

	return fn(tx)
}
