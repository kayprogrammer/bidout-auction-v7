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

// @Summary Get Profile
// @Description This endpoint gets the current user's profile.
// @Tags Auctioneer
// @Success 200 {object} schemas.ProfileResponseSchema
// @Router /auctioneer [get]
// @Security BearerAuth
func GetProfile(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	user := c.Locals("user").(*models.User)

	userData := schemas.ProfileResponseDataSchema{
		FirstName: user.FirstName,
		LastName:  user.LastName,
		Avatar:    user.GetAvatarUrl(db),
	}
	response := schemas.ProfileResponseSchema{
		ResponseSchema: schemas.ResponseSchema{Message: "User details fetched!"}.Init(),
		Data:           userData,
	}
	return c.Status(200).JSON(response)
}

// @Summary Update Profile
// @Description This endpoint updates an authenticated user's profile. Note: use the returned upload_url to upload avatar to cloudinary
// @Tags Auctioneer
// @Param user body schemas.UpdateProfileSchema true "Update User"
// @Success 200 {object} schemas.UpdateProfileResponseSchema
// @Failure 422 {object} utils.ErrorResponse
// @Router /auctioneer [put]
// @Security BearerAuth
func UpdateProfile(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	user := c.Locals("user").(*models.User)
	validator := utils.Validator()

	updateProfileData := schemas.UpdateProfileSchema{}

	// Validate request
	if errCode, errData := DecodeJSONBody(c, &updateProfileData); errData != nil {
		return c.Status(errCode).JSON(errData)
	}
	if err := validator.Validate(updateProfileData); err != nil {
		return c.Status(422).JSON(err)
	}

	fileType := updateProfileData.FileType
	if fileType != nil {
		file := models.File{ResourceType: *fileType}
		if user.AvatarId == nil {
			db.Create(&file)
		} else {
			db.Model(models.File{BaseModel: models.BaseModel{ID: *user.AvatarId}}).Updates(&file)
			file.ID = *user.AvatarId
		}
		user.AvatarId = &file.ID
	}
	user.FirstName = updateProfileData.FirstName
	user.LastName = updateProfileData.LastName
	db.Save(&user)

	userData := schemas.UpdateProfileResponseDataSchema{
		FirstName:      user.FirstName,
		LastName:       user.LastName,
		FileUploadData: user.GetAvatarUploadUrl(db),
	}
	response := schemas.UpdateProfileResponseSchema{
		ResponseSchema: schemas.ResponseSchema{Message: "User updated!"}.Init(),
		Data:           userData,
	}
	return c.Status(200).JSON(response)
}

// @Summary Retrieve all listings by the current user
// @Description This endpoint retrieves all listings by the current user.
// @Tags Auctioneer
// @Param quantity query int false  "Listings Quantity"
// @Success 200 {object} schemas.ListingsResponseSchema
// @Router /auctioneer/listings [get]
// @Security BearerAuth
func GetAuctioneerListings(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	user := c.Locals("user").(*models.User)

	listings := []models.Listing{}
	quantity := c.QueryInt("quantity")
	// Get listings
	db.Preload(clause.Associations).Order("created_at DESC").Find(&listings, models.Listing{AuctioneerId: user.ID})

	// Initialize each listing object in the slice
	for i := range listings {
		listings[i] = listings[i].Init(db)
	}
	if quantity > 0 {
		listings = listings[:quantity]
	}
	response := schemas.ListingsResponseSchema{
		ResponseSchema: schemas.ResponseSchema{Message: "Auctioneer Listings fetched"}.Init(),
		Data:           listings,
	}
	return c.Status(200).JSON(response)
}

// @Summary Create a listing
// @Description This endpoint creates a new listing. Note: Use the returned upload_url to upload image to cloudinary
// @Tags Auctioneer
// @Param listing body schemas.CreateListingSchema true "Create Listing"
// @Success 200 {object} schemas.CreateListingResponseSchema
// @Failure 422 {object} utils.ErrorResponse
// @Router /auctioneer/listings [post]
// @Security BearerAuth
func CreateListing(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	user := c.Locals("user").(*models.User)
	validator := utils.Validator()

	createListingData := schemas.CreateListingSchema{}
	// Validate request
	if errCode, errData := DecodeJSONBody(c, &createListingData); errData != nil {
		return c.Status(errCode).JSON(errData)
	}

	if err := validator.Validate(createListingData); err != nil {
		return c.Status(422).JSON(err)
	}
	var categoryId *uuid.UUID
	categorySlug := createListingData.Category
	// Validate Category
	if categorySlug != "other" {
		category := models.Category{Slug: &categorySlug}
		db.Take(&category, category)
		if category.ID == uuid.Nil {
			data := map[string]string{
				"category": "Invalid category!",
			}
			return c.Status(422).JSON(utils.ErrorResponse{Message: "Invalid Entry", Data: &data}.Init())
		}
		categoryId = &category.ID

	} else {
		categoryId = nil
	}
	fileType := createListingData.FileType
	file := models.File{ResourceType: fileType}
	db.Create(&file)

	listing := models.Listing{
		AuctioneerId: user.ID,
		Name:         createListingData.Name,
		Desc:         createListingData.Desc,
		CategoryId:   categoryId,
		Active:       true,
		Price:        utils.DecimalParser(createListingData.Price),
		ClosingDate:  utils.TimeParser(createListingData.ClosingDate),
		ImageId:      file.ID,
	}
	db.Create(&listing)
	db.Preload(clause.Associations).Take(&listing, listing.ID)

	listingData := schemas.CreateListingResponseDataSchema{
		Listing:        listing.Init(db),
		FileUploadData: listing.GetImageUploadData(db),
	}
	response := schemas.CreateListingResponseSchema{
		ResponseSchema: schemas.ResponseSchema{Message: "Listing created successfully"}.Init(),
		Data:           listingData,
	}
	return c.Status(201).JSON(response)
}

// @Summary Update a listing
// @Description This endpoint updates a particular listing. Note: Use the returned upload_url to upload image to cloudinary
// @Tags Auctioneer
// @Param slug path string true  "Listing Slug"
// @Param listing body schemas.UpdateListingSchema true "Update Listing"
// @Success 200 {object} schemas.CreateListingResponseSchema
// @Failure 422 {object} utils.ErrorResponse
// @Router /auctioneer/listings/{slug} [patch]
// @Security BearerAuth
func UpdateListing(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	user := c.Locals("user").(*models.User)
	validator := utils.Validator()

	listingSlug := c.Params("slug")
	listing := models.Listing{Slug: &listingSlug}
	db.Take(&listing, listing)
	if listing.ID == uuid.Nil {
		return c.Status(404).JSON(utils.ErrorResponse{Message: "Invalid listing!"}.Init())
	}

	if listing.AuctioneerId != user.ID {
		return c.Status(400).JSON(utils.ErrorResponse{Message: "This listing doesn't belong to you!"}.Init())
	}

	updateListingData := schemas.UpdateListingSchema{}

	// Validate request
	if errCode, errData := DecodeJSONBody(c, &updateListingData); errData != nil {
		return c.Status(errCode).JSON(errData)
	}
	if err := validator.Validate(updateListingData); err != nil {
		return c.Status(422).JSON(err)
	}
	categorySlug := updateListingData.Category
	if categorySlug != nil {
		// Validate Category
		other := "other"
		if categorySlug != &other {
			category := models.Category{Slug: categorySlug}
			db.Take(&category, category)
			if category.ID == uuid.Nil {
				data := map[string]string{
					"category": "Invalid category!",
				}
				return c.Status(422).JSON(utils.ErrorResponse{Message: "Invalid Entry", Data: &data}.Init())
			}
			listing.CategoryId = &category.ID

		} else {
			listing.CategoryId = nil
		}
	}

	fileType := updateListingData.FileType
	if fileType != nil {
		file := models.File{ResourceType: *fileType}
		db.Model(models.File{BaseModel: models.BaseModel{ID: listing.ImageId}}).Updates(&file)
	}

	// Assign data to listing
	utils.AssignFields(updateListingData, &listing)
	db.Save(&listing)
	db.Preload(clause.Associations).Take(&listing, listing.ID)

	listingData := schemas.CreateListingResponseDataSchema{
		Listing:        listing.Init(db),
		FileUploadData: listing.GetImageUploadData(db),
	}
	response := schemas.CreateListingResponseSchema{
		ResponseSchema: schemas.ResponseSchema{Message: "Listing updated successfully"}.Init(),
		Data:           listingData,
	}
	return c.Status(200).JSON(response)
}

// @Summary Retrieve bids in a listing (current user)
// @Description This endpoint retrieves all bids in a particular listing by the current user.
// @Tags Auctioneer
// @Param slug path string true  "Listing Slug"
// @Success 200 {object} schemas.BidsResponseSchema
// @Failure 404 {object} utils.ErrorResponse
// @Router /auctioneer/listings/{slug}/bids [get]
// @Security BearerAuth
func GetAuctioneerListingBids(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	user := c.Locals("user").(*models.User)

	listingSlug := c.Params("slug")

	listing := models.Listing{Slug: &listingSlug}
	db.Preload("Bids", func(db *gorm.DB) *gorm.DB {
		return db.Order("updated_at DESC").Limit(3) // Order by updated
	}).Find(&listing, listing)
	if listing.ID == uuid.Nil {
		return c.Status(404).JSON(utils.ErrorResponse{Message: "Invalid listing!"}.Init())
	}
	if listing.AuctioneerId != user.ID {
		return c.Status(400).JSON(utils.ErrorResponse{Message: "This listing doesn't belong to you!"}.Init())
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
