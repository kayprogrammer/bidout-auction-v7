package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kayprogrammer/bidout-auction-v7/models"
	"github.com/satori/go.uuid"
)

type Client struct {
	ID					uuid.UUID
	Type				string			// guest or user
}

func GetClient(c *fiber.Ctx) *Client {
	clientContext := c.Locals("client")
	client := Client{}
	if clientContext == nil {
		return nil
	} else if user, ok := clientContext.(*models.User); ok {
		client.ID = user.ID
		client.Type = "user"
	} else if guestUser, ok := clientContext.(models.GuestUser); ok {
		client.ID = guestUser.ID
		client.Type = "guest"
	}
	return &client
}

// func ParseRequestBody()