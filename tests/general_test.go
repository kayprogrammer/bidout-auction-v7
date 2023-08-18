package tests

import (
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"

	"github.com/kayprogrammer/bidout-auction-v7/models"
	"github.com/kayprogrammer/bidout-auction-v7/utils"
)

func getSiteDetails(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	t.Run("Get Site Details", func(t *testing.T) {
		url := fmt.Sprintf("%s/site-detail", baseUrl)
		req := httptest.NewRequest("GET", url, nil)
		res, _ := app.Test(req)

		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Site Details Fetched!", body["message"])
		dataKeys := []string{"address", "email", "fb", "ig", "name", "phone", "tw", "wh"}
		assert.Equal(t, true, utils.KeysExistInMap(dataKeys, body["data"].(map[string]interface{})))
	})
}

func subscribe(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	t.Run("Subscribe", func(t *testing.T) {
		url := fmt.Sprintf("%s/subscribe", baseUrl)
		validEmail := "test_subscriber@email.com"
		emailData := models.Subscriber{Email: validEmail}

		res := ProcessTestBody(t, app, url, "POST", emailData)

		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Subscription successful!", body["message"])
		expectedData := make(map[string]interface{})
		expectedData["email"] = validEmail
		assert.Equal(t, expectedData, body["data"].(map[string]interface{}))
	})
}

func getReviews(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	t.Run("Get Reviews", func(t *testing.T) {
		reviewText := "This is a nice new platform"
		url := fmt.Sprintf("%s/reviews", baseUrl)

		// Create Review
		reviewer := CreateTestUser(db)
		review := models.Review{ReviewerId: reviewer.ID, Show: true, Text: reviewText}
		db.Create(&review)

		// Make request
		req := httptest.NewRequest("GET", url, nil)
		res, _ := app.Test(req)

		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Reviews fetched!", body["message"])

		var expectedData []map[string]interface{}
		expectedReviewerData := map[string]interface{}{
			"name":   reviewer.FullName(),
			"avatar": nil,
		}
		expectedReviewData := map[string]interface{}{
			"reviewer": expectedReviewerData,
			"text":     reviewText,
		}
		expectedData = append(expectedData, expectedReviewData)
		data, _ := json.Marshal(body["data"])
		expectedDataJson, _ := json.Marshal(expectedData)
		assert.Equal(t, expectedDataJson, data)
	})
}
func TestGeneral(t *testing.T) {
	app := fiber.New()
	db := Setup(t, app)
	BASEURL := "/api/v7/general"

	// Run General Endpoint Tests
	getSiteDetails(t, app, db, BASEURL)
	subscribe(t, app, db, BASEURL)
	getReviews(t, app, db, BASEURL)

	// Drop Tables and Close Connectiom
	DropTables(db)
	CloseTestDatabase(db)
}
