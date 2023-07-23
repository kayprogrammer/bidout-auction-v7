package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"github.com/kayprogrammer/bidout-auction-v7/database"
	"github.com/kayprogrammer/bidout-auction-v7/routes"
	"github.com/kayprogrammer/bidout-auction-v7/utils"
)

func main() {
	godotenv.Load()
	database.ConnectDb()
	utils.CreateInitialData()
	app := fiber.New()
	routes.SetupRoutes(app)
	log.Fatal(app.Listen(":8000"))
}