package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kayprogrammer/bidout-auction-v7/models"
	"github.com/kayprogrammer/bidout-auction-v7/schemas"
	"gorm.io/gorm"
)

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
		LastName: user.LastName,
		Avatar: user.GetAvatarUrl(db),
	}
	response := schemas.ProfileResponseSchema{
		ResponseSchema: schemas.ResponseSchema{Message: "User details fetched!"}.Init(),
		Data:           userData,
	}
	return c.Status(200).JSON(response)
}