package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kayprogrammer/bidout-auction-v7/models"
	"github.com/kayprogrammer/bidout-auction-v7/schemas"
	"github.com/kayprogrammer/bidout-auction-v7/utils"
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
	// Get listings
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

// @Summary Retrieve listing's detail
// @Description This endpoint retrieves detail of a listing.
// @Tags Listings
// @Param slug path string true  "Listing Slug"
// @Success 200 {object} schemas.ListingDetailResponseSchema
// @Router /api/v7/listings/detail/{slug} [get]
func GetListing(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	listing := models.Listing{}
	slug := c.Params("slug")

	// Get listing
	db.Preload(clause.Associations).Find(&listing,"slug = ?", slug)
	if listing.ID == uuid.Nil {
		return c.Status(404).JSON(utils.ErrorResponse{Message: "Listing does not exist!"}.Init())
	}
	listing = listing.Init(db)
	relatedListings := []models.Listing{}
	db.Preload(clause.Associations).Order("created_at DESC").Where("category_id = ?", listing.CategoryId).Not("id = ?", listing.ID).Find(&relatedListings)

	response := schemas.ListingDetailResponseSchema{
		ResponseSchema: schemas.ResponseSchema{Message: "Listing details fetched"}.Init(),
		Data:           schemas.ListingDetailResponseDataSchema{
			Listing:			listing,
			RelatedListings:	relatedListings,
		},
	}
	return c.Status(200).JSON(response)
}

