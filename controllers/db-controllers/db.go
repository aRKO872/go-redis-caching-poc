package db_controllers

import (
	"log"

	"github.com/go-redis-caching-poc/config"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

var DB *gorm.DB
var err error
 
 
func DatabaseConnection() {
	DB, err = gorm.Open(postgres.Open(config.GetConfig().DATABASE_URL), &gorm.Config{})
	if err != nil {
		log.Fatal("db connection error: ", err)
	}
	log.Println("db connection successful")
}