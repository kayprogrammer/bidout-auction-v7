package main

import (
	"log"
	"github.com/gofiber/swagger"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/kayprogrammer/bidout-auction-v7/database"
	"github.com/kayprogrammer/bidout-auction-v7/routes"
	"github.com/kayprogrammer/bidout-auction-v7/initials"
	"github.com/kayprogrammer/bidout-auction-v7/config"

	_ "github.com/kayprogrammer/bidout-auction-v7/docs"
)

// @title Bidout Auction API
// @version 7.0
// @description A simple bidding API built with Fiber
// @Accept json
// @Produce json
// @BasePath  /api/v7
// @Security BearerAuth
// @securityDefinitions.apikey BearerAuth 
// @in header 
// @name Authorization 
// @description Type 'Bearer jwt_string' to correctly set the API Key
// @Security GuestUserAuth
// @securityDefinitions.apikey GuestUserAuth 
// @in header 
// @name GuestUserId 
// @description For guest watchlists. Get ID (uuid) from '/api/v7/listings/watchlist' POST endpoint
func main() {
	cfg := config.GetConfig()
	database.ConnectDb()
	db := database.Database.Db
	initials.CreateInitialData(db)

	app := fiber.New()

	// Set up the database middleware
	app.Use(database.DatabaseMiddleware)

	// CORS config
	app.Use(cors.New(cors.Config{
		AllowOrigins: cfg.CORSAllowedOrigins,
		AllowHeaders:  "Origin, Content-Type, Accept, Authorization, Guestuserid, Access-Control-Allow-Origin, Content-Disposition",
		AllowCredentials: true,
		AllowMethods: "GET, POST, PUT, PATCH, DELETE, OPTIONS",
	}))

	// Register routes
	routes.SetupRoutes(app)
	app.Get("/*", swagger.HandlerDefault) // default

	log.Fatal(app.Listen(":8000"))
}