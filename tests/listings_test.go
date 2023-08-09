package tests

import (
	"fmt"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/kayprogrammer/bidout-auction-v7/utils"
	"github.com/kayprogrammer/bidout-auction-v7/models"
	"github.com/kayprogrammer/bidout-auction-v7/schemas"
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
		assert.Equal(t, true, (len(data) > 0))
	})
}

func getListing(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	// Drop and Create Tables since the previous test uses the create_listing it...
	DropTables(db)
	CreateTables(db)

	// Create Listing
	listing := CreateListing(db)

	t.Run("Get Listing", func(t *testing.T) {
		// Verify that a particular listing retrieval fails with an invalid slug
		url := fmt.Sprintf("%s/detail/invalid_slug", baseUrl)

		// Make request
		req := httptest.NewRequest("GET", url, nil)
		res, _ := app.Test(req)

		// Assert Status code
		assert.Equal(t, 404, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "Listing does not exist!", body["message"])

		// Verify that a particular listing is retrieved successfully
		slug := listing.Slug
		url = fmt.Sprintf("%s/detail/%s", baseUrl, *slug)

		// Make request
		req = httptest.NewRequest("GET", url, nil)
		res, _ = app.Test(req)

		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body = ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Listing details fetched", body["message"])

		// Parse and assert body
		dataKeys := []string{"listing", "related_listings"}
		assert.Equal(t, true, utils.KeysExistInMap(dataKeys, body["data"].(map[string]interface{})))
	})
}

func getWatchlistListings(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	// Drop and Create Tables since the previous test uses the create_listing it...
	DropTables(db)
	CreateTables(db)

	// Create Listing
	listing := CreateListing(db)

	// Create Watchlist
	watchlist := models.Watchlist{UserId: &listing.AuctioneerId, ListingId: listing.ID}
	db.Create(&watchlist)

	t.Run("Get Watchlist Listings", func(t *testing.T) {
		url := fmt.Sprintf("%s/watchlist", baseUrl)

		// Make request
		req := httptest.NewRequest("GET", url, nil)
		jwt := CreateJwt(db, listing.AuctioneerId)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwt.Access))
		res, _ := app.Test(req)

		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Watchlist Listings fetched", body["message"])

		data, _ := json.Marshal(body["data"])
		assert.Equal(t, true, (len(data) > 0))
	})
}

func createOrRemoveUserWatchlistsListing(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	// Drop and Create Tables since the previous test uses the user table it...
	DropTables(db)
	CreateTables(db)

	t.Run("Create Or Remove User Watchlists Listing", func(t *testing.T) {
		listing := CreateListing(db)

		url := fmt.Sprintf("%s/watchlist", baseUrl)
		addRemoveWatchlistData := schemas.AddOrRemoveWatchlistSchema{
			Slug: "invalid_listing_slug", // Invalid listing slug
		}

		res := ProcessTestBody(t, app, url, "POST", addRemoveWatchlistData)

		// # Test for invalid slug
		// Assert Status code
		assert.Equal(t, 404, res.StatusCode)
		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "Listing does not exist!", body["message"])


		// Verify that the watchlist was created successfully
		addRemoveWatchlistData.Slug = *listing.Slug
		jwt := CreateJwt(db, listing.AuctioneerId)
		res = ProcessTestBody(t, app, url, "POST", addRemoveWatchlistData, jwt.Access)
		// Assert response
		assert.Equal(t, 201, res.StatusCode)
		// Parse and assert body
		body = ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Listing added to user watchlist", body["message"])
	})
}

func TestListing(t *testing.T) {
	app := fiber.New()
	db := Setup(t, app)
	BASEURL := "/api/v7/listings"

	// Run Listings Endpoint Tests
	getListings(t, app, db, BASEURL)
	getListing(t, app, db, BASEURL)
	getWatchlistListings(t, app, db, BASEURL)
	createOrRemoveUserWatchlistsListing(t, app, db, BASEURL)

	// Drop Tables and Close Connectiom
	DropTables(db)
	CloseTestDatabase(db)
}
