package tests

import (
	"fmt"
	"log"
	"net/http/httptest"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/kayprogrammer/bidout-auction-v7/utils"
	"github.com/stretchr/testify/assert"
)

var BASEURL = "/api/v7/general"

func TestGetSiteDetails(t *testing.T) {
	app := fiber.New()
	db := Setup(t, app)

	url := fmt.Sprintf("%s/site-detail", BASEURL)
	req := httptest.NewRequest("GET", url, nil)
	res, _ := app.Test(req)

	// Assert Status code
	assert.Equal(t, res.StatusCode, 200)

	// Parse and assert body
	body := ParseResponseBody(t, res.Body).(map[string]interface{})
	log.Println(body)
	assert.Equal(t, body["status"], "success")
	assert.Equal(t, body["message"], "Site Details Fetched!")
	dataKeys := []string{"address", "email", "fb", "ig", "name", "phone", "tw", "wh"}
	assert.Equal(t, utils.KeysExistInMap(dataKeys, body["data"].(map[string]interface{})), true)
	
	// Drop Tables and Close Connectiom
	DropTables(db)
	CloseTestDatabase(db)
}
