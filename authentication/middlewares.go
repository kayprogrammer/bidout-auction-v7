package authentication

import (
	"github.com/gofiber/fiber/v2"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"

	"github.com/kayprogrammer/bidout-auction-v7/models"
	"github.com/kayprogrammer/bidout-auction-v7/utils"
)

func getUser(c *fiber.Ctx, token string, db *gorm.DB) (*models.User, *string) {
	if len(token) < 8 {
		err := "Auth Token is Invalid or Expired!"
		return nil, &err
	}
	user, err := DecodeAccessToken(token[7:], db)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func AuthMiddleware(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	db := c.Locals("db").(*gorm.DB)

	if len(token) < 1 {
		return c.Status(401).JSON(utils.ErrorResponse{Message: "Unauthorized User!"}.Init())
	}
	user, err := getUser(c, token, db)
	if err != nil {
		return c.Status(401).JSON(utils.ErrorResponse{Message: *err}.Init())
	}
	c.Locals("user", user)
	return c.Next()
}

func isUUID(input string) bool {
    _, err := uuid.FromString(input)
    return err == nil
}

func ClientMiddleware(c *fiber.Ctx) error {
	token := c.Get("Authorization")
	guestId := c.Get("guestuserid")
	db := c.Locals("db").(*gorm.DB)

	if len(token) < 1 {
		// Try for Guest
		c.Locals("client", nil)
		if len(guestId) > 0 {
			guest := models.GuestUser{}
			if !isUUID(guestId) {
				return c.Status(401).JSON(utils.ErrorResponse{Message: "Invalid type for guest id (use a uuid)"}.Init())
			}
			db.Find(&guest, "id = ?", guestId)
			if guest.ID != uuid.Nil {
				c.Locals("client", guest)
			}
		}
	} else {
		// Auth User becomes client
		user, err := getUser(c, token, db)
		if err != nil {
			return c.Status(401).JSON(utils.ErrorResponse{Message: *err}.Init())
		}
		c.Locals("client", user)
	}
	return c.Next()
}
