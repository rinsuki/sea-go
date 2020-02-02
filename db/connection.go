package db

import (
	"os"

	"github.com/jinzhu/gorm"
)

var connection *gorm.DB

func GetConnection() *gorm.DB {
	if connection != nil {
		return connection
	}

	db, err := gorm.Open("postgres", os.Getenv("DATABASE_URL"))
	if err != nil {
		panic(err)
	}
	db.LogMode(true)
	connection = db
	return connection
}
