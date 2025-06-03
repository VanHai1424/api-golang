package migration

import (
	coursemodel "crawdata/module/course/model"
	"log"

	"gorm.io/gorm"
)

func Migrate(db *gorm.DB) error {
	err := db.AutoMigrate(&coursemodel.Course{})
	if err != nil {
		return err
	}
	log.Println("Migration completed.")
	return nil
}

func InsertSampleData(db *gorm.DB) error {
	err := SeedData(db)
	if err != nil {
		return err
	}
	log.Println("Sample data inserted successfully.")
	return nil
}
