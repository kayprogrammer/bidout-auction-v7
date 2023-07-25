package database

import (
	"log"
	"os"
	"fmt"

	"github.com/kayprogrammer/bidout-auction-v7/models"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type DbInstance struct {
	Db *gorm.DB
}

var Database DbInstance

func ConnectDb() {
	dsnTemplate := "host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s"
	dsn := fmt.Sprintf(
		dsnTemplate,
		os.Getenv("POSTGRES_SERVER"),
		os.Getenv("POSTGRES_USER"),
		os.Getenv("POSTGRES_PASSWORD"),
		os.Getenv("POSTGRES_DB"),
		os.Getenv("POSTGRES_PORT"),
		"disable",
		"UTC",
	)
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{
		SkipDefaultTransaction: true,
		PrepareStmt: true,
	})

	if err != nil {
		log.Fatal("Failed to connect to the database! \n", err.Error())
		os.Exit(2)
	}
	log.Println("Connected to the database successfully")
	db.Logger = logger.Default.LogMode(logger.Info)
	log.Println("Running Migrations")

	// Add UUID extension
	result := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
    if result.Error != nil {
        log.Fatal("failed to create extension: " + result.Error.Error())
    }

	// Add Migrations
	db.AutoMigrate(
		// base
		&models.File{}, 
		
		// general
		&models.SiteDetail{}, 
		&models.Subscriber{}, 
		&models.Review{}, 

		// accounts
		&models.User{}, 
		&models.Jwt{}, 
		&models.Otp{},

		// listings
		&models.Category{}, 
		&models.Listing{}, 
		&models.Bid{},
		&models.Watchlist{},
	)

	Database = DbInstance{Db: db}
}