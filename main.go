package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/kayprogrammer/bidout-auction-v7/database"
	"github.com/kayprogrammer/bidout-auction-v7/routes"
	"github.com/kayprogrammer/bidout-auction-v7/initials"
)

func main() {
	database.ConnectDb()
	db := database.Database.Db
	initials.CreateInitialData(db)
	app := fiber.New()
	routes.SetupRoutes(app)
	log.Fatal(app.Listen(":8000"))
}