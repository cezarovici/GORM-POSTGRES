package database

import (
	"github.com/cezarovici/GORM-POSTGRES/models"
	"gorm.io/gorm"

	"fmt"
	"log"
	"os"

	"gorm.io/driver/postgres"
	"gorm.io/gorm/logger"
)

// DB represents a Database instance
var DB *gorm.DB

var dbEnvs = fmt.Sprintf(
	"host=db user=%s password=%s dbname=%s port=5432 sslmode=disable TimeZone=Asia/Shanghai",
	os.Getenv("DB_USER"),
	os.Getenv("DB_PASSWORD"),
	os.Getenv("DB_NAME"),
)

func ConnectDb() error {
	db, err := gorm.Open(postgres.Open(dbEnvs), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	})

	if err != nil {
		log.Fatal("Failed to connect to database. \n", err)
		return err
	}

	DB = db

	log.Println("connected")
	DB.Logger = logger.Default.LogMode(logger.Info)

	log.Println("running migrations")
	DB.AutoMigrate(&models.User{})

	return nil
}
