package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kayprogrammer/bidout-auction-v7/models"
	"github.com/kayprogrammer/bidout-auction-v7/schemas"
	"github.com/kayprogrammer/bidout-auction-v7/utils"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// @Summary Retrieve all listings
// @Description This endpoint retrieves all listings.
// @Tags Listings
// @Param quantity query int false  "Listings Quantity"
// @Success 200 {object} schemas.ListingsResponseSchema
// @Router /listings [get]
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
			db.Where("(user_id = ? OR guest_user_id = ?) AND listing_id = ?", client.ID, client.ID, listings[i].ID).Take(&watchlist)
			if watchlist.ID != uuid.Nil {
				listings[i].Watchlist = true
			}
		}
	}
	if quantity > 0 {
		listings = listings[:quantity]
	}
	response := schemas.ListingsResponseSchema{
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
// @Router /listings/detail/{slug} [get]
func GetListing(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	listing := models.Listing{}
	slug := c.Params("slug")

	// Get listing
	db.Preload(clause.Associations).Find(&listing, "slug = ?", slug)
	if listing.ID == uuid.Nil {
		return c.Status(404).JSON(utils.ErrorResponse{Message: "Listing does not exist!"}.Init())
	}
	listing = listing.Init(db)
	relatedListings := []models.Listing{}
	db.Preload(clause.Associations).Order("created_at DESC").Where("category_id = ? AND NOT id = ?", listing.CategoryId, listing.ID).Limit(3).Find(&relatedListings)

	response := schemas.ListingDetailResponseSchema{
		ResponseSchema: schemas.ResponseSchema{Message: "Listing details fetched"}.Init(),
		Data: schemas.ListingDetailResponseDataSchema{
			Listing:         listing,
			RelatedListings: relatedListings,
		},
	}
	return c.Status(200).JSON(response)
}

// @Summary Retrieve all listings by users watchlist
// @Description This endpoint retrieves all watchlist listings.
// @Tags Listings
// @Success 200 {object} schemas.ListingsResponseSchema
// @Router /listings/watchlist [get]
// @Security BearerAuth
// @Security GuestUserAuth
func GetWatchlistListings(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	client := GetClient(c)
	watchlists := []models.Watchlist{}
	listings := []models.Listing{}

	// Get watchlists
	if client != nil {
		db.Preload("Listing.AuctioneerObj").Preload("Listing.CategoryObj").Preload(clause.Associations).Where("user_id = ?", client.ID).Or("guest_user_id = ?", client.ID).Order("created_at DESC").Find(&watchlists)
	}
	for i := range watchlists {
		listing := watchlists[i].Listing
		listing.Watchlist = true
		listings = append(listings, listing.Init(db))
	}

	response := schemas.ListingsResponseSchema{
		ResponseSchema: schemas.ResponseSchema{Message: "Watchlist Listings fetched"}.Init(),
		Data:           listings,
	}
	return c.Status(200).JSON(response)
}

// @Summary Add or Remove listing from a users watchlist
// @Description This endpoint adds or removes a listing from a user's watchlist, authenticated or not.... As a guest, ensure to store guestuser_id in localstorage and keep passing it to header 'guestuserid' in subsequent requests
// @Tags Listings
// @Param listing_slug body schemas.AddOrRemoveWatchlistSchema true "Add/Remove Watchlist"
// @Success 201 {object} schemas.AddOrRemoveWatchlistResponseSchema
// @Success 200 {object} schemas.AddOrRemoveWatchlistResponseSchema
// @Failure 422 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /listings/watchlist [post]
// @Security BearerAuth
// @Security GuestUserAuth
func AddOrRemoveWatchlistListing(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	validator := utils.Validator()

	addRemoveWatchlistData := schemas.AddOrRemoveWatchlistSchema{}

	// Validate request
	if errCode, errData := DecodeJSONBody(c, &addRemoveWatchlistData); errData != nil {
		return c.Status(errCode).JSON(errData)
	}
	if err := validator.Validate(addRemoveWatchlistData); err != nil {
		return c.Status(422).JSON(err)
	}

	// Get listing
	listing := models.Listing{}
	db.Preload(clause.Associations).Find(&listing, "slug = ?", addRemoveWatchlistData.Slug)
	if listing.ID == uuid.Nil {
		return c.Status(404).JSON(utils.ErrorResponse{Message: "Listing does not exist!"}.Init())
	}

	client := GetClient(c)
	if client == nil {
		guestUser := models.GuestUser{}
		db.Create(&guestUser)
		client = &Client{ID: guestUser.ID, Type: "guest"}
	}

	respMessage := "Listing added to user watchlist"
	statusCode := 201
	// Check if watchlist exists
	watchlist := models.Watchlist{}
	result := db.Where("(user_id = ? OR guest_user_id = ?) AND listing_id = ?", client.ID, client.ID, listing.ID).Take(&watchlist)
	if result.Error == gorm.ErrRecordNotFound {
		// Create Watchlist
		watchlistToCreate := models.Watchlist{ListingId: listing.ID}
		if client.Type == "user" {
			watchlistToCreate.UserId = &client.ID
		} else {
			watchlistToCreate.GuestUserId = &client.ID
		}
		db.Create(&watchlistToCreate)
	} else {
		respMessage = "Listing removed from user watchlist"
		statusCode = 200
		db.Delete(&watchlist)
	}

	var guestUserId *uuid.UUID
	guestUserId = &client.ID
	if client.Type == "user" {
		guestUserId = nil
	}
	response := schemas.AddOrRemoveWatchlistResponseSchema{
		ResponseSchema: schemas.ResponseSchema{Message: respMessage}.Init(),
		Data: schemas.AddOrRemoveWatchlistResponseDataSchema{
			GuestUserId: guestUserId,
		},
	}
	return c.Status(statusCode).JSON(response)
}

// @Summary Retrieve all categories
// @Description This endpoint retrieves all categories
// @Tags Listings
// @Success 200 {object} schemas.CategoriesResponseSchema
// @Router /listings/categories [get]
func GetCategories(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)

	// Get categories
	categories := []models.Category{}
	db.Find(&categories)

	response := schemas.CategoriesResponseSchema{
		ResponseSchema: schemas.ResponseSchema{Message: "Categories fetched"}.Init(),
		Data:           categories,
	}
	return c.Status(200).JSON(response)
}


// @Summary Retrieve all listings by category
// @Description This endpoint retrieves all listings in a particular category. Use slug 'other' for category other
// @Tags Listings
// @Param slug path string true  "Category Slug"
// @Success 200 {object} schemas.ListingsResponseSchema
// @Failure 404 {object} utils.ErrorResponse
// @Router /listings/categories/{slug} [get]
func GetCategoryListings(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	client := GetClient(c)
	categorySlug := c.Params("slug")
	
	// Get Category
	var categoryId *uuid.UUID
	if categorySlug == "other" {
		categoryId = nil
	} else {
		category := models.Category{}
		db.First(&category, "slug = ?", categorySlug)
		if category.ID == uuid.Nil {
			return c.Status(404).JSON(utils.ErrorResponse{Message: "Invalid category!"}.Init())
		}
		categoryId = &category.ID
	}
	
	// Get listings
	listings := []models.Listing{}
	db.Preload(clause.Associations).Order("created_at DESC").Find(&listings, "category_id = ?", categoryId)

	// Initialize each listing object in the slice
	for i := range listings {
		listings[i] = listings[i].Init(db)
		if client != nil {
			watchlist := models.Watchlist{}
			db.Where("(user_id = ? OR guest_user_id = ?) AND listing_id = ?", client.ID, client.ID, listings[i].ID).Take(&watchlist)
			if watchlist.ID != uuid.Nil {
				listings[i].Watchlist = true
			}
		}
	}
	response := schemas.ListingsResponseSchema{
		ResponseSchema: schemas.ResponseSchema{Message: "Category Listings fetched"}.Init(),
		Data:           listings,
	}
	return c.Status(200).JSON(response)
}

// @Summary Retrieve bids in a listing
// @Description This endpoint retrieves at most 3 bids from a particular listing.
// @Tags Listings
// @Param slug path string true  "Listing Slug"
// @Success 200 {object} schemas.BidsResponseSchema
// @Failure 404 {object} utils.ErrorResponse
// @Router /listings/detail/{slug}/bids [get]
func GetListingBids(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	listingSlug := c.Params("slug")

	listing := models.Listing{}
	db.Preload("Bids", func(db *gorm.DB) *gorm.DB {
		return db.Order("updated_at DESC").Limit(3) // Order by updated
	}).Find(&listing, "slug = ?", listingSlug)
	if listing.ID == uuid.Nil {
		return c.Status(404).JSON(utils.ErrorResponse{Message: "Invalid listing!"}.Init())
	}

	// Get Bids
	bids := listing.Bids
	for i := range bids {
		bids[i] = bids[i].Init(db)
	}

	response := schemas.BidsResponseSchema{
		ResponseSchema: schemas.ResponseSchema{Message: "Listing Bids fetched"}.Init(),
		Data:           schemas.BidResponseDataSchema{Listing: listing.Name, Bids: bids},
	}
	return c.Status(200).JSON(response)
}

// @Summary Add a bid to a listing
// @Description This endpoint adds a bid to a particular listing.
// @Tags Listings
// @Param slug path string true  "Listing Slug"
// @Param amount body schemas.CreateBidSchema true "Create Bid"
// @Success 201 {object} schemas.BidResponseSchema
// @Failure 422 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /listings/detail/{slug}/bids [post]
// @Security BearerAuth
func CreateBid(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	user := c.Locals("user").(*models.User)
	listingSlug := c.Params("slug")

	// Get Listing
	listing := models.Listing{}
	db.Preload("Bids", func(db *gorm.DB) *gorm.DB {
		return db.Order("updated_at DESC").Limit(1) // Order by updated (Latest bid is surely the highest bid)
	}).Find(&listing, "slug = ?", listingSlug)
	if listing.ID == uuid.Nil {
		return c.Status(404).JSON(utils.ErrorResponse{Message: "Listing does not exist!"}.Init())
	}
	highestBid := listing.GetHighestBid()

	validator := utils.Validator()
	createBidData := schemas.CreateBidSchema{}

	// Validate request
	if errCode, errData := DecodeJSONBody(c, &createBidData); errData != nil {
		return c.Status(errCode).JSON(errData)
	}
	if err := validator.Validate(createBidData); err != nil {
		return c.Status(422).JSON(err)
	}

	amount := utils.DecimalParser(createBidData.Amount)
	if user.ID == listing.AuctioneerId {
		return c.Status(403).JSON(utils.ErrorResponse{Message: "You cannot bid your own product!"}.Init())
	} else if !listing.Active {
		return c.Status(410).JSON(utils.ErrorResponse{Message: "This auction is closed!"}.Init())
	} else if listing.TimeLeft() < 1 {
		return c.Status(410).JSON(utils.ErrorResponse{Message: "This auction is expired and closed!"}.Init())
	} else if amount.Cmp(listing.Price) < 0 {
		return c.Status(400).JSON(utils.ErrorResponse{Message: "Bid amount cannot be less than the bidding price!"}.Init())
	} else if amount.Cmp(highestBid) <= 0 {
		return c.Status(400).JSON(utils.ErrorResponse{Message: "Bid amount must be more than the highest bid!"}.Init())
	}

	// Check for existing bid
	bid := models.Bid{UserId: user.ID, ListingId: listing.ID}
	db.Where(bid).First(&bid)

	// Create or update
	bid.Amount = amount
	db.Save(&bid)

	response := schemas.BidResponseSchema{
		ResponseSchema: schemas.ResponseSchema{Message: "Bid added to listing"}.Init(),
		Data: bid.Init(db),
	}
	return c.Status(201).JSON(response)
}