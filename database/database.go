package database

import (
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB

// use GORM with postgreSQL connection: https://gorm.io/docs/connecting_to_the_database.html#PostgreSQL
func InitDB() (err error) {
	dsn := os.Getenv("SQL_DSN")
	DB, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}

	return nil
}

func AutoMigrateTables(models ...interface{}) error {
	if err := DB.AutoMigrate(models...); err != nil {
		return err
	}

	return nil
}
