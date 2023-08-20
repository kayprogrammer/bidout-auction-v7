package tests

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	"github.com/shopspring/decimal"
	
	"github.com/kayprogrammer/bidout-auction-v7/models"
	"github.com/kayprogrammer/bidout-auction-v7/schemas"
)

func getProfile(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	user := CreateTestVerifiedUser(db)
	access := CreateJwt(db, user.ID).Access
	t.Run("Get Profile", func(t *testing.T) {
		url := baseUrl

		// Make request
		req := httptest.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access))
		res, _ := app.Test(req)

		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "User details fetched!", body["message"])
	})
}

func updateProfile(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	user := CreateTestVerifiedUser(db)
	access := CreateJwt(db, user.ID).Access
	t.Run("Update Profile", func(t *testing.T) {
		url := baseUrl
		fileType := "image/jpeg"
		updateProfileData := schemas.UpdateProfileSchema{
			FirstName: "FirstNameChanged",
			LastName:  "LastNameChanged",
			FileType:  &fileType,
		}
		// Make request
		res := ProcessTestBody(t, app, url, "PUT", updateProfileData, access)
		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "User updated!", body["message"])
	})
}

func getAuctioneerListings(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	user := CreateTestVerifiedUser(db)
	access := CreateJwt(db, user.ID).Access
	t.Run("Get Auctioneer Listings", func(t *testing.T) {
		url := fmt.Sprintf("%s/listings", baseUrl)

		// Create Listing
		CreateListing(db)

		// Make request
		req := httptest.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access))
		res, _ := app.Test(req)

		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Auctioneer Listings fetched", body["message"])

		data, _ := json.Marshal(body["data"])
		assert.Equal(t, true, (len(data) > 0))
	})
}

func createListing(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	user := CreateTestVerifiedUser(db)
	access := CreateJwt(db, user.ID).Access
	category := models.Category{Name: "Test Category"}
	db.Create(&category)
	t.Run("Create Listing", func(t *testing.T) {
		url := fmt.Sprintf("%s/listings", baseUrl)
		fileType := "image/jpeg"
		createListingData := schemas.CreateListingSchema{
			Name:        "Test Listing",
			Desc:        "Test description",
			Category:    *category.Slug,
			Price:       1000.20,
			ClosingDate: "2250-01-02T15:04:05.000Z",
			FileType:    fileType,
		}
		// Make request
		res := ProcessTestBody(t, app, url, "POST", createListingData, access)
		// Assert Status code
		assert.Equal(t, 201, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Listing created successfully", body["message"])

		// Verify that create listing failed with invalid category
		createListingData.Category = "invalid-category"
		res = ProcessTestBody(t, app, url, "POST", createListingData, access)
		// Assert Status code
		assert.Equal(t, 422, res.StatusCode)

		// Parse and assert body
		body = ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "Invalid Entry", body["message"])
		expectedData := make(map[string]interface{})
		expectedData["category"] = "Invalid category!"
		assert.Equal(t, expectedData, body["data"].(map[string]interface{}))
	})
}

func updateListing(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	user := CreateTestVerifiedUser(db)
	access := CreateJwt(db, user.ID).Access
	listing := CreateListing(db)

	category := models.Category{}
	db.Take(&category, listing.CategoryId)

	t.Run("Update Listing", func(t *testing.T) {
		url := fmt.Sprintf("%s/listings/%s", baseUrl, *listing.Slug)
		name := "Test Listing"
		desc := "Test description"
		price := 2000.50
		updateListingData := schemas.UpdateListingSchema{
			Name:  &name,
			Desc:  &desc,
			Price: &price,
		}
		// Make request
		res := ProcessTestBody(t, app, url, "PATCH", updateListingData, access)
		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Listing updated successfully", body["message"])

		// Verify the request fails with invalid listing slug
		url = fmt.Sprintf("%s/listings/%s", baseUrl, "invalid_slug")
		// Make request
		res = ProcessTestBody(t, app, url, "PATCH", updateListingData, access)
		// Assert Status code
		assert.Equal(t, 404, res.StatusCode)

		// Parse and assert body
		body = ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "Invalid listing!", body["message"])

		// You can test for other error responses
	})
}

func getAuctioneerListingBids(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	user := CreateTestVerifiedUser(db)
	access := CreateJwt(db, user.ID).Access
	listing := CreateListing(db)
	anotherVerifiedUser := CreateAnotherTestVerifiedUser(db)

	t.Run("Get Auctioneer Listing Bids", func(t *testing.T) {
		url := fmt.Sprintf("%s/listings/invalid_listing_slug/bids", baseUrl)

		// Make request
		req := httptest.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access))
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
		url = fmt.Sprintf("%s/listings/%s/bids", baseUrl, *listing.Slug)
		// Make request
		req = httptest.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", access))
		res, _ = app.Test(req)

		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body = ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Listing Bids fetched", body["message"])

		data, _ := json.Marshal(body["data"])
		assert.Equal(t, true, (len(data) > 0))

		// You can test for other error responses
	})
}

func TestAuctioneer(t *testing.T) {
	app := fiber.New()
	db := Setup(t, app)
	BASEURL := "/api/v7/auctioneer"

	// Run Auctioneer Endpoint Tests
	getProfile(t, app, db, BASEURL)
	updateProfile(t, app, db, BASEURL)
	getAuctioneerListings(t, app, db, BASEURL)
	createListing(t, app, db, BASEURL)
	updateListing(t, app, db, BASEURL)
	getAuctioneerListingBids(t, app, db, BASEURL)

	// Drop Tables and Close Connectiom
	DropTables(db)
	CloseTestDatabase(db)
}
