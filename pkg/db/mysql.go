package db

import (
	"database/sql"
	"fmt"

	_ "github.com/go-sql-driver/mysql"
)

type Sql struct {
	Db       *sql.DB
	Host     string
	User     string
	Password string
	Port     string
	Database string
}

func (s *Sql) InitDB() error {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=true",
		s.User, s.Password, s.Host, s.Port, s.Database)

	db, err := sql.Open("mysql", dsn)
	if err != nil {
		return err
	}

	if err := db.Ping(); err != nil {
		return err
	}

	s.Db = db
	return nil
}
