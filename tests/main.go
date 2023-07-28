package tests

import (
	"fmt"
	"log"
	"testing"
	"time"

	"github.com/ory/dockertest/v3"
	_ "github.com/ory/dockertest/v3/docker"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"github.com/kayprogrammer/bidout-auction-v7/utils"

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
	return db
}

func SetupTestDatabase(t *testing.T) *gorm.DB {
	pool, err := dockertest.NewPool("")
	if err != nil {
		t.Fatalf("Failed to create Docker pool: %v", err)
	}

	// Generate a unique name for the test database
	testDBName := fmt.Sprintf("testdb_%d", time.Now().UnixNano())

	// Generate a random password for the test database
	testDBPassword := utils.GenerateRandomPassword()

	resource, err := pool.Run("postgres", "latest", []string{
		fmt.Sprintf("POSTGRES_PASSWORD=%s", testDBPassword),
		fmt.Sprintf("POSTGRES_DB=%s", testDBName),
	})
	if err != nil {
		t.Fatalf("Failed to start PostgreSQL container: %v", err)
	}

	t.Cleanup(func() {
		if err := pool.Purge(resource); err != nil {
			t.Errorf("Failed to clean up PostgreSQL container: %v", err)
		}
	})

	// Wait for the database to be ready
	dsn := "host=localhost port=" + resource.GetPort("5432/tcp") + " user=postgres dbname=" + testDBName + " password=" + testDBPassword + " sslmode=disable"
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

func Client(t *testing.T) {
	// Set up the test database
	db := SetupTestDatabase(t)
	defer closeTestDatabase(db)
}