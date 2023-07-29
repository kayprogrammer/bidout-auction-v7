package tests

import (
	"fmt"
	"log"
	"testing"
	"time"
	"io"
	"encoding/json"

	"github.com/gofiber/fiber/v2"
	"github.com/kayprogrammer/bidout-auction-v7/config"
	"github.com/kayprogrammer/bidout-auction-v7/models"
	"github.com/kayprogrammer/bidout-auction-v7/routes"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func waitForDBConnection(t *testing.T, dsn string) *gorm.DB {
	maxRetries := 10 // Number of retries to wait for the database to be ready
	var db *gorm.DB
	var err error

	for i := 0; i < maxRetries; i++ {
		t.Logf("Waiting for the database to be ready... Attempt %d", i+1)
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			break
		}
		time.Sleep(1 * time.Second)
	}

	if err != nil {
		t.Fatalf("Failed to connect to the test database: %v", err)
	}

	// Add UUID extension
	result := db.Exec("CREATE EXTENSION IF NOT EXISTS \"uuid-ossp\";")
    if result.Error != nil {
        log.Fatal("failed to create extension: " + result.Error.Error())
    }

	// Migrations
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

	return db
}

func SetupTestDatabase(t *testing.T) *gorm.DB {
	cfg := config.GetConfig()
	MockedDB(CREATE)
    dsn := fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s TimeZone=%s",
		cfg.PostgresServer,
		cfg.PostgresUser,
		cfg.PostgresPassword,
		"bidout_auction_v7_by_kayprogrammer_test_db",
		cfg.PostgresPort,
		"disable",
		"UTC",
	)
	return waitForDBConnection(t, dsn)
}

func closeTestDatabase(db *gorm.DB) {
    sqlDB, err := db.DB()
    if err != nil {
        log.Fatal("Failed to get database connection: " + err.Error())
    }
    if err := sqlDB.Close(); err != nil {
        log.Fatal("Failed to close database connection: " + err.Error())
    }
}

func Setup(t *testing.T, app *fiber.App) *gorm.DB {
	// Set up the test database
	db := SetupTestDatabase(t)
	defer closeTestDatabase(db)

	// Inject your test database into the Fiber app's context
	app.Use(func(c *fiber.Ctx) error {
		c.Locals("db", db)
		return c.Next()
	})
	routes.SetupRoutes(app)
	return db
}

func ParseResponseBody(t *testing.T, b io.ReadCloser) interface{} {
	body, _ := io.ReadAll(b)
	// Parse the response body as JSON
    responseBody := make(map[string]interface{})
    err := json.Unmarshal(body, &responseBody)
    if err != nil {
        t.Errorf("error parsing response body as JSON: %s", err)
    }
	return responseBody
}