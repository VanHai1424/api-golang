package coursestorage

import (
	"gorm.io/gorm"
)

type MySQLStorage struct {
	db *gorm.DB
}

func NewMySQLStorage(db *gorm.DB) *MySQLStorage {
	return &MySQLStorage{db: db}
}
