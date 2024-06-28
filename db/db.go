package db

import (
	"database/sql"
	"log"

	"github.com/go-sql-driver/mysql"
)

// NewMySQLStorage creates a new MySQL connection
func NewMySQLStorage(cfg mysql.Config) (*sql.DB, error) {
	db, err := sql.Open("mysql", cfg.FormatDSN())
	if err != nil {
		log.Fatal(err)
	}

	return db, nil
}
