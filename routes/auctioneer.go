package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kayprogrammer/bidout-auction-v7/models"
	"github.com/kayprogrammer/bidout-auction-v7/schemas"
	"github.com/kayprogrammer/bidout-auction-v7/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

var fileTypes = []string{
	"image/bmp",
	"image/gif",
	"image/jpeg",
	"image/png",
	"image/tiff",
	"image/webp",
	"image/svg+xml",
}

// @Summary Get Profile
// @Description This endpoint gets the current user's profile.
// @Tags Auctioneer
// @Success 200 {object} schemas.ProfileResponseSchema
// @Router /api/v7/auctioneer [get]
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
// @Router /api/v7/auctioneer [put]
// @Security BearerAuth
func UpdateProfile(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	user := c.Locals("user").(*models.User)
	validator := utils.Validator()

	updateProfileData := schemas.UpdateProfileSchema{}
	c.BodyParser(&updateProfileData)

	// Validate request
	if err := validator.Validate(updateProfileData); err != nil {
		return c.Status(422).JSON(err)
	}

	fileType := updateProfileData.FileType
	if fileType != nil {
		// Validate file type
		fileTypeFound := false
		for _, value := range fileTypes {
			if value == *fileType {
				fileTypeFound = true
				break
			}
		}
		if !fileTypeFound {
			data := map[string]string{
				"file_type": "Invalid file type!",
			}
			return c.Status(422).JSON(utils.ErrorResponse{Message: "Invalid Entry", Data: &data}.Init())
		}
		file := models.File{ResourceType: *fileType}
		if user.AvatarId == nil {
			db.Create(&file)
		} else {
			db.Model(models.File{}).Where("id = ?", user.AvatarId).Updates(&file)
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
// @Router /api/v7/auctioneer/listings [get]
// @Security BearerAuth
func GetAuctioneerListings(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	user := c.Locals("user").(*models.User)

	listings := []models.Listing{}
	quantity := c.QueryInt("quantity")
	// Get listings
	db.Preload(clause.Associations).Order("created_at DESC").Find(&listings,"auctioneer_id = ?", user.ID)

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
