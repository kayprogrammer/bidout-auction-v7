package main

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/kayprogrammer/simple-fiber-ecommerce/database"
	"github.com/kayprogrammer/simple-fiber-ecommerce/routes"
	"github.com/joho/godotenv"
)

func main() {
	godotenv.Load()
	app := fiber.New()
	routes.SetupRoutes(app)
	log.Fatal(app.Listen(":8000"))
}