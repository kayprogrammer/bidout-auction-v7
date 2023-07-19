package main

import (
	"log"
	"github.com/joho/godotenv"
	"github.com/gofiber/fiber/v2"
	"github.com/kayprogrammer/bidout-auction-v7/database"
	"github.com/kayprogrammer/bidout-auction-v7/routes"
)

func main() {
	godotenv.Load()
	database.ConnectDb()
	app := fiber.New()
	routes.SetupRoutes(app)
	log.Fatal(app.Listen(":8000"))
}