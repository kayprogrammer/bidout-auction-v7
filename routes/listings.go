package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kayprogrammer/bidout-auction-v7/models"
	"github.com/kayprogrammer/bidout-auction-v7/schemas"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// @Summary Retrieve all listings
// @Description This endpoint retrieves all listings.
// @Tags Listings
// @Param quantity query int false  "Listings Quantity"
// @Success 200 {object} schemas.ListingsResponseSchemas
// @Router /api/v7/listings [get]
func GetListings(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	listings := []models.Listing{}
	quantity := c.QueryInt("quantity")
	// Get reviews
	db.Preload(clause.Associations).Find(&listings)

	// Initialize each listing object in the slice
	for i := range listings {
		listings[i] = listings[i].Init(db)
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

