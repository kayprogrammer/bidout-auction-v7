package tests

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/kayprogrammer/bidout-auction-v7/models"
	"github.com/kayprogrammer/bidout-auction-v7/utils"
	"github.com/stretchr/testify/assert"
	"gorm.io/gorm"
)

var BASEURL = "/api/v7/general"

func getSiteDetails(t *testing.T, app *fiber.App, db *gorm.DB) {
	t.Run("getSiteDetails", func(t *testing.T) {
		url := fmt.Sprintf("%s/site-detail", BASEURL)
		req := httptest.NewRequest("GET", url, nil)
		res, _ := app.Test(req)

		// Assert Status code
		assert.Equal(t, res.StatusCode, 200)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, body["status"], "success")
		assert.Equal(t, body["message"], "Site Details Fetched!")
		dataKeys := []string{"address", "email", "fb", "ig", "name", "phone", "tw", "wh"}
		assert.Equal(t, utils.KeysExistInMap(dataKeys, body["data"].(map[string]interface{})), true)
	})
}

func subscribe(t *testing.T, app *fiber.App, db *gorm.DB) {
	t.Run("subscribe", func(t *testing.T) {
		url := fmt.Sprintf("%s/subscribe", BASEURL)
		validEmail := "test_subscriber@email.com"
		emailData := models.Subscriber{Email: validEmail}

		// Marshal the test data to JSON
		requestBytes, err := json.Marshal(emailData)
		requestBody := bytes.NewReader(requestBytes)
		assert.Nil(t, err)
		req := httptest.NewRequest("POST", url, requestBody)
		req.Header.Set("Content-Type", "application/json")
		res, _ := app.Test(req)

		// Assert Status code
		assert.Equal(t, res.StatusCode, 200)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, body["status"], "success")
		assert.Equal(t, body["message"], "Subscription successful!")
		expectedData := make(map[string]interface{})
		expectedData["email"] = validEmail
		assert.Equal(t, body["data"].(map[string]interface{}), expectedData)
	})
}

func getReviews(t *testing.T, app *fiber.App, db *gorm.DB) {
	t.Run("getReviews", func(t *testing.T) {
		reviewText := "This is a nice new platform"
		url := fmt.Sprintf("%s/reviews", BASEURL)

		// Create Review
		reviewer := CreateTestUser(db)
		review := models.Review{ReviewerId: reviewer.ID, Show: true, Text: reviewText}
		db.Create(&review)

		// Make request
		req := httptest.NewRequest("GET", url, nil)
		res, _ := app.Test(req)

		// Assert Status code
		assert.Equal(t, res.StatusCode, 200)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, body["status"], "success")
		assert.Equal(t, body["message"], "Reviews fetched!")

		var expectedData []map[string]interface{}
		expectedReviewerData := map[string]interface{}{
			"name": reviewer.FullName(),
			"avatar": nil,
		} 
		expectedReviewData := map[string]interface{}{
			"reviewer":  expectedReviewerData,
			"text":   reviewText,
		}
		expectedData = append(expectedData, expectedReviewData)
		data, _ := json.Marshal(body["data"])
		expectedDataJson, _ := json.Marshal(expectedData)
		assert.Equal(t, data, expectedDataJson)
	})
}
func TestGeneral(t *testing.T) {
	app := fiber.New()
	db := Setup(t, app)

	// Run General Endpoint Tests
	getSiteDetails(t, app, db)
	subscribe(t, app, db)
	getReviews(t, app, db)

	// Drop Tables and Close Connectiom
	DropTables(db)
	CloseTestDatabase(db)
}