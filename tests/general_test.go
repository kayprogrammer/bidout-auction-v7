package tests

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"github.com/gofiber/fiber/v2"
)

func TestHelloHandler(t *testing.T) {
	app := fiber.New()

	// Set up a test request to the "/" route
	req := httptest.NewRequest(http.MethodGet, "/", nil)

	// Create a response recorder to capture the response
	res := httptest.NewRecorder()

	// Pass the request and response recorder to the app handler
	app.ServeHTTP(res, req)

	// Check the response status code
	if res.Code != http.StatusOK {
		t.Errorf("Expected status code %d, but got %d", http.StatusOK, res.Code)
	}

	// Check the response body
	expectedBody := "Hello, World!"
	if res.Body.String() != expectedBody {
		t.Errorf("Expected response body '%s', but got '%s'", expectedBody, res.Body.String())
	}
}
