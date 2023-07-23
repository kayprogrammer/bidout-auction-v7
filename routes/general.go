package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kayprogrammer/bidout-auction-v7/database"
	"github.com/kayprogrammer/bidout-auction-v7/models"
	"github.com/kayprogrammer/bidout-auction-v7/schemas"
	"github.com/kayprogrammer/bidout-auction-v7/utils"
)

func GetSiteDetails(c *fiber.Ctx) error {
	var sitedetail models.SiteDetail

	database.Database.Db.FirstOrCreate(&sitedetail, &sitedetail)

	responseSiteDetail := schemas.SiteDetailResponseSchema{
		ResponseSchema: schemas.ResponseSchema{Message: "Site Details Fetched!"}.Init(),
		Data:           sitedetail,
	}
	return c.Status(200).JSON(responseSiteDetail)
}

func Subscribe(c *fiber.Ctx) error {
	validator := utils.Validator()
	subscriber := models.Subscriber{}

	c.BodyParser(&subscriber)

	// Validate request
	if err := validator.Validate(subscriber); err != nil {
		return c.Status(422).JSON(err)
	}

	// Create subscriber
	database.Database.Db.Where(models.Subscriber{Email: subscriber.Email}).FirstOrCreate(&subscriber)

	responseSubscriber := schemas.SubscriberResponseSchema{
		ResponseSchema: schemas.ResponseSchema{Message: "Subscription successful!"}.Init(),
		Data:           subscriber,
	}
	return c.Status(200).JSON(responseSubscriber)
}

func GetReviews(c *fiber.Ctx) error {
	reviews := []models.Review{}
	db := database.Database.Db

	// Get reviews
	db.Where(models.Review{Show: true}).Find(&reviews)

	// Initialize each review object in the slice
	for i := range reviews {
		reviews[i] = reviews[i].Init(db)
	}

	responseReviews := schemas.ReviewsResponseSchema{
		ResponseSchema: schemas.ResponseSchema{Message: "Reviews fetched!"}.Init(),
		Data:           reviews,
	}
	return c.Status(200).JSON(responseReviews)
}
