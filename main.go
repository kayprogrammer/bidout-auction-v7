package main

import (
	"log"
	"github.com/gofiber/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/kayprogrammer/bidout-auction-v7/database"
	"github.com/kayprogrammer/bidout-auction-v7/routes"
	"github.com/kayprogrammer/bidout-auction-v7/initials"

	_ "github.com/kayprogrammer/bidout-auction-v7/docs"
)

// @title Bidout Auction API
// @version 7.0
// @description A simple bidding API built with Fiber
// @Accept json
// @Produce json
// @Security BearerAuth
// @securityDefinitions.apikey BearerAuth 
// @in header 
// @name Authorization 
// @description "Type 'Bearer jwt_string' to correctly set the API Key"
func main() {
	database.ConnectDb()
	db := database.Database.Db
	initials.CreateInitialData(db)

	app := fiber.New()

	// Set up the database middleware
	app.Use(database.DatabaseMiddleware)

	// Register routes
	routes.SetupRoutes(app)
	app.Get("/*", swagger.HandlerDefault) // default

	app.Get("/*", swagger.New(swagger.Config{ // custom
		URL: "http://example.com/doc.json",
		DeepLinking: false,
		// Expand ("list") or Collapse ("none") tag groups by default
		DocExpansion: "none",
		// Prefill OAuth ClientId on Authorize popup
		OAuth: &swagger.OAuthConfig{
			AppName:  "OAuth Provider",
			ClientId: "21bb4edc-05a7-4afc-86f1-2e151e4ba6e2",
		},
		// Ability to change OAuth2 redirect uri location
		OAuth2RedirectUrl: "http://127.0.0.1:8000/swagger/oauth2-redirect.html",
	}))
	log.Fatal(app.Listen(":8000"))
}