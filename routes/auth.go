package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/satori/go.uuid"

	"github.com/kayprogrammer/bidout-auction-v7/database"
	"github.com/kayprogrammer/bidout-auction-v7/models"
	"github.com/kayprogrammer/bidout-auction-v7/schemas"
	"github.com/kayprogrammer/bidout-auction-v7/senders"
	"github.com/kayprogrammer/bidout-auction-v7/utils"
)

// @Summary Register a new user
// @Description This endpoint registers new users into our application.
// @Tags Auth
// @Param user body models.User true "User object"
// @Success 201 {object} schemas.RegisterResponseSchema
// @Failure 422 {object} utils.ErrorResponse
// @Router /api/v7/auth/register [post]
func Register(c *fiber.Ctx) error {
	db := database.Database.Db
	validator := utils.Validator()

	user := models.User{}
	c.BodyParser(&user)

	// Validate request
	if err := validator.Validate(user); err != nil {
		return c.Status(422).JSON(err)
	}

	db.Find(&user,"email = ?", user.Email)
	if user.ID != uuid.Nil {
		data := map[string]string{
			"email": "Email already registered!", 
		}
		return c.Status(422).JSON(utils.ErrorResponse{Message: "Invalid Entry", Data: &data}.Init())
	}

	// Create User
	db.Create(&user)

	// Send Email
	go senders.SendEmail(db, user, "activate")

	response := schemas.RegisterResponseSchema{
		ResponseSchema: schemas.ResponseSchema{Message: "Registration successful"}.Init(),
		Data:           schemas.EmailRequestSchema{Email: user.Email},
	}
	return c.Status(200).JSON(response)
}

