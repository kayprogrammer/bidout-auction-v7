package tests

import (
	"log"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

func getListings(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	t.Run("Get Listings", func(t *testing.T) {
		url := baseUrl

		// Create Listing
		CreateListing(db)

		// Make request
		req := httptest.NewRequest("GET", url, nil)
		res, _ := app.Test(req)

		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Listings fetched", body["message"])

		data, _ := json.Marshal(body["data"])
		log.Println(data)
		assert.Equal(t, true, (len(data) > 0))
	})
}
func TestListing(t *testing.T) {
	app := fiber.New()
	db := Setup(t, app)
	BASEURL := "/api/v7/listings"

	// Run Listings Endpoint Tests
	getListings(t, app, db, BASEURL)

	// Drop Tables and Close Connectiom
	DropTables(db)
	CloseTestDatabase(db)
}
