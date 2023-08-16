package tests

import (
	"fmt"
	// "encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
	// "github.com/shopspring/decimal"

	// "github.com/kayprogrammer/bidout-auction-v7/utils"
	// "github.com/kayprogrammer/bidout-auction-v7/models"
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
			LastName: "LastNameChanged",
			FileType: &fileType,
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

func TestAuctioneer(t *testing.T) {
	app := fiber.New()
	db := Setup(t, app)
	BASEURL := "/api/v7/auctioneer"

	// Run Auctioneer Endpoint Tests
	getProfile(t, app, db, BASEURL)
	updateProfile(t, app, db, BASEURL)

	// Drop Tables and Close Connectiom
	DropTables(db)
	CloseTestDatabase(db)
}
