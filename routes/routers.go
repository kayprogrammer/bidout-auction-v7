package routes

import (
	"github.com/gofiber/fiber/v2"
	midw "github.com/kayprogrammer/bidout-auction-v7/authentication"
)

func SetupRoutes(app *fiber.App) {
	api := app.Group("/api/v7")

	// General Routes
	generalRouter := api.Group("/general")
	generalRouter.Get("/site-detail", GetSiteDetails)
	generalRouter.Post("/subscribe", Subscribe)
	generalRouter.Get("/reviews", GetReviews)
	


	// Auth Routes
	authRouter := api.Group("/auth")
	authRouter.Post("/register", Register)
	authRouter.Post("/verify-email", VerifyEmail)
	authRouter.Post("/resend-verification-email", ResendVerificationEmail)
	authRouter.Post("/send-password-reset-otp", SendPasswordResetOtp)
	authRouter.Post("/set-new-password", SetNewPassword)
	authRouter.Post("/login", midw.ClientMiddleware, Login)
	authRouter.Post("/refresh", Refresh)
	authRouter.Get("/logout", midw.AuthMiddleware, Logout)

	// Listings Routes
	listingsRouter := api.Group("/listings")
	listingsRouter.Get("", midw.ClientMiddleware, GetListings)
	listingsRouter.Get("/detail/:slug", GetListing)
	listingsRouter.Get("/watchlist", midw.ClientMiddleware, GetWatchlistListings)
	listingsRouter.Post("/watchlist", midw.ClientMiddleware, AddOrRemoveWatchlistListing)
	listingsRouter.Get("/categories", GetCategories)
	listingsRouter.Get("/categories/:slug", GetCategoryListings)
	listingsRouter.Get("/detail/:slug/bids", GetListingBids)
	listingsRouter.Post("/detail/:slug/bids", midw.AuthMiddleware, CreateBid)

	// Auctioneer Routes
	auctioneerRouter := api.Group("/auctioneer")
	auctioneerRouter.Get("", midw.AuthMiddleware, GetProfile)
	auctioneerRouter.Put("", midw.AuthMiddleware, UpdateProfile)
	auctioneerRouter.Get("/listings", midw.AuthMiddleware, GetAuctioneerListings)
	auctioneerRouter.Post("/listings", midw.AuthMiddleware, CreateListing)
	auctioneerRouter.Patch("/listings/:slug", midw.AuthMiddleware, UpdateListing)
	auctioneerRouter.Get("/listings/:slug/bids", midw.AuthMiddleware, GetAuctioneerListingBids)
}