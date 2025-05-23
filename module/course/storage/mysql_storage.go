package coursestorage

import (
	"database/sql"
	"time"
)

type MySQLStorage struct {
	db *sql.DB
}

var vnLoc *time.Location

func init() {
	var err error
	vnLoc, err = time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		panic(err)
	}
}

func NewMySQLStorage(db *sql.DB) *MySQLStorage {
	return &MySQLStorage{db: db}
}
