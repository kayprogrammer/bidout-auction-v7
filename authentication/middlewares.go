package authentication

import (
	"github.com/gofiber/fiber/v2"
	"gorm.io/gorm"
	"github.com/kayprogrammer/bidout-auction-v7/utils"
)

func AuthMiddleware(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	db := c.Locals("db").(*gorm.DB)

	if len(token) < 1 {
		return c.Status(401).JSON(utils.ErrorResponse{Message: "Unauthorized User!"}.Init())
	}
	if len(token) < 8 {
		return c.Status(401).JSON(utils.ErrorResponse{Message: "Auth Token is Invalid or Expired!"}.Init())
	}
	user, err := DecodeAccessToken(token[7:], db)
	if err != nil {
		return c.Status(401).JSON(utils.ErrorResponse{Message: *err}.Init())
	}
	c.Locals("user", user)
	return c.Next()
}