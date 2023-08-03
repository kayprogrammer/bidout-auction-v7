package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kayprogrammer/bidout-auction-v7/models"
	"github.com/kayprogrammer/bidout-auction-v7/schemas"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"github.com/satori/go.uuid"
)

// @Summary Retrieve all listings
// @Description This endpoint retrieves all listings.
// @Tags Listings
// @Param quantity query int false  "Listings Quantity"
// @Success 200 {object} schemas.ListingsResponseSchemas
// @Router /api/v7/listings [get]
// @Security BearerAuth
// @Security GuestUserAuth
func GetListings(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	client := GetClient(c)
	listings := []models.Listing{}
	quantity := c.QueryInt("quantity")
	// Get reviews
	db.Preload(clause.Associations).Order("created_at DESC").Find(&listings)

	// Initialize each listing object in the slice
	for i := range listings {
		listings[i] = listings[i].Init(db)
		if client != nil {
			watchlist := models.Watchlist{}
			db.Where("user_id = ?", client.ID).Or("guest_user_id = ?", client.ID).Where("listing_id = ?", listings[i].ID).Find(&watchlist)
			if watchlist.ID != uuid.Nil {
				listings[i].Watchlist = true
			}
		}
	}
	if quantity > 0 {
		listings = listings[:quantity]
	}
	response := schemas.ListingsResponseSchemas{
		ResponseSchema: schemas.ResponseSchema{Message: "Listings fetched"}.Init(),
		Data:           listings,
	}
	return c.Status(200).JSON(response)
}

