package routes

import (
	"github.com/gofiber/fiber/v2"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v7")

	// General Routes
	generalRouter := api.Group("/general")
	generalRouter.Get("/site-detail", GetSiteDetails)
	generalRouter.Post("/subscribe", Subscribe)
	generalRouter.Get("/reviews", GetReviews)


	// Auth Routes
	// authRouter := api.Group("/auth")

	// // Listings Routes
	// listingsRouter := api.Group("/listings")

	// // Auctioneer Routes
	// auctioneerRouter := api.Group("/auctioneer")


}