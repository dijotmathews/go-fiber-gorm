package database

import (
	"log"
	"os"

	"github.com/dijotmathews/go-fiber-gorm/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// DbInstance ...
type DbInstance struct {
	Db *gorm.DB
}

// Database ...
var Database DbInstance

func ConnectDb() {
	db, err := gorm.Open(sqlite.Open("api.db"), &gorm.Config{})

	if err != nil {
		log.Fatal("Failed to connect to the database \n", err.Error())
		os.Exit(2)
	}

	log.Println("connected to the database successfully")

	db.Logger = logger.Default.LogMode(logger.Info)
	log.Println("Running Migration")

	db.AutoMigrate(&models.User{}, &models.Product{}, &models.Order{})

	Database = DbInstance{Db: db}

}
