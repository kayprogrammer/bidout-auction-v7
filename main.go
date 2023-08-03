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
// @description Type 'Bearer jwt_string' to correctly set the API Key
// @Security GuestUserId
// @securityDefinitions.apikey GuestUserId 
// @in header 
// @name GuestUserId 
// @description For guest watchlists. Get ID from '/api/v7/listings/watchlist' POST endpoint
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
	}))
	log.Fatal(app.Listen(":8000"))
}