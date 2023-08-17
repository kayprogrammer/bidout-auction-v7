package tests

import (
	"fmt"
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"github.com/shopspring/decimal"

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
	DropSingleTable(db, models.Listing{})
	CreateSingleTable(db, models.Listing{})

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
	// Drop and Create Tables since the previous test uses the watchlist table it...
	DropSingleTable(db, models.Watchlist{})
	CreateSingleTable(db, models.Watchlist{})

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

func getCategories(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	// Since our previous test already makes use of a category, then category exists in our db

	t.Run("Get Categories", func(t *testing.T) {
		url := fmt.Sprintf("%s/categories", baseUrl)

		// Make request
		req := httptest.NewRequest("GET", url, nil)
		res, _ := app.Test(req)

		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Categories fetched", body["message"])

		data, _ := json.Marshal(body["data"])
		assert.Equal(t, true, (len(data) > 0))
	})
}

func getCategoryListings(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	listing := CreateListing(db)
	t.Run("Get Category Listings", func(t *testing.T) {
		url := fmt.Sprintf("%s/categories/invalid_category_slug", baseUrl)

		// Make request
		req := httptest.NewRequest("GET", url, nil)
		res, _ := app.Test(req)

		// Verify that listings by an invalid category slug fails
		// Assert Status code
		assert.Equal(t, 404, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "Invalid category!", body["message"])

		// Verify that listings by a valid category slug succeeds
		category := models.Category{}
		db.Find(&category,"id = ?", listing.CategoryId)
		url = fmt.Sprintf("%s/categories/%s", baseUrl, *category.Slug)
		// Make request
		req = httptest.NewRequest("GET", url, nil)
		res, _ = app.Test(req)

		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body = ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Category Listings fetched", body["message"])

		data, _ := json.Marshal(body["data"])
		assert.Equal(t, true, (len(data) > 0))
	})
}

func getListingBids(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	listing := CreateListing(db)
	anotherVerifiedUser := CreateAnotherTestVerifiedUser(db)

	t.Run("Get Listing Bids", func(t *testing.T) {
		url := fmt.Sprintf("%s/detail/invalid_listing_slug/bids", baseUrl)

		// Make request
		req := httptest.NewRequest("GET", url, nil)
		res, _ := app.Test(req)

		// Verify that bids by an invalid listing slug fails
		// Assert Status code
		assert.Equal(t, 404, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "Invalid listing!", body["message"])

		// Verify that bids by a valid listing slug succeeds
		bid := models.Bid{UserId: anotherVerifiedUser.ID, ListingId: listing.ID, Amount: decimal.NewFromFloat(2000.00)}
		db.Create(&bid)
		url = fmt.Sprintf("%s/detail/%s/bids", baseUrl, *listing.Slug)
		// Make request
		req = httptest.NewRequest("GET", url, nil)
		res, _ = app.Test(req)

		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body = ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Listing Bids fetched", body["message"])

		data, _ := json.Marshal(body["data"])
		assert.Equal(t, true, (len(data) > 0))
	})
}

func createBid(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	listing := CreateListing(db)
	anotherVerifiedUser := CreateAnotherTestVerifiedUser(db)

	t.Run("Create Bid", func(t *testing.T) {

		url := fmt.Sprintf("%s/detail/invalid_listing_slug/bids", baseUrl)
		createBidData := schemas.CreateBidSchema{
			Amount: 2000.00,
		}
		jwt := CreateJwt(db, listing.AuctioneerId)
		res := ProcessTestBody(t, app, url, "POST", createBidData, jwt.Access)

		// Test for invalid slug
		// Assert Status code
		assert.Equal(t, 404, res.StatusCode)
		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "Listing does not exist!", body["message"])

		url = fmt.Sprintf("%s/detail/%s/bids", baseUrl, *listing.Slug)
		// Test for invalid user
		res = ProcessTestBody(t, app, url, "POST", createBidData, jwt.Access)		
		assert.Equal(t, 403, res.StatusCode) // Assert Status code
		// Parse and assert body
		body = ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "You cannot bid your own product!", body["message"])

		// Test for failure for lesser bidding price
		createBidData.Amount = 200.00
		jwt = CreateJwt(db, anotherVerifiedUser.ID)
		res = ProcessTestBody(t, app, url, "POST", createBidData, jwt.Access)
		// Assert Status code		
		assert.Equal(t, 400, res.StatusCode)
		// Parse and assert body
		body = ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "Bid amount cannot be less than the bidding price!", body["message"])

		// Verify that the bid was created successfully
		createBidData.Amount = 2000.00
		res = ProcessTestBody(t, app, url, "POST", createBidData, jwt.Access)
		// Assert response
		assert.Equal(t, 201, res.StatusCode)
		// Parse and assert body
		body = ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Bid added to listing", body["message"])

		// You can also test for other error responses.....
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
	getCategories(t, app, db, BASEURL)
	getCategoryListings(t, app, db, BASEURL)
	getListingBids(t, app, db, BASEURL)
	createBid(t, app, db, BASEURL)

	// Drop Tables and Close Connectiom
	DropTables(db)
	CloseTestDatabase(db)
}
