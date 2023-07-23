package initials

import (
	"errors"

	"github.com/kayprogrammer/bidout-auction-v7/models"
	"github.com/kayprogrammer/bidout-auction-v7/utils"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"
)

func createSuperUser(db *gorm.DB) models.User {
	user := models.User{
		FirstName:       "Test",
		LastName:        "Admin",
		Email:           "testadmin@email.com",
		Password:        utils.HashPassword("testadmin"),
		IsSuperuser:     true,
		IsStaff:         true,
		IsEmailVerified: true,
	}
	db.Where(models.User{Email: user.Email}).FirstOrCreate(&user)
	return user
}

func createAutioneer(db *gorm.DB) models.User {
	user := models.User{
		FirstName:       "Test",
		LastName:        "Auctioneer",
		Email:           "testauctioneer@email.com",
		Password:        utils.HashPassword("testauctioneer"),
		IsEmailVerified: true,
	}
	db.Where(models.User{Email: user.Email}).FirstOrCreate(&user)
	return user
}

func createReviewer(db *gorm.DB) models.User {
	user := models.User{
		FirstName:       "Test",
		LastName:        "Reviewer",
		Email:           "testreviewer@email.com",
		Password:        utils.HashPassword("testreviewer"),
		IsEmailVerified: true,
	}
	db.Where(models.User{Email: user.Email}).FirstOrCreate(&user)
	return user
}

func createReviews(db *gorm.DB, reviewerId uuid.UUID) {
	review := &models.Review{
		ReviewerId: reviewerId,
		Show:       true,
		Text:       "Maecenas vitae porttitor neque, ac porttitor nunc. Duis venenatis lacinia libero. Nam nec augue ut nunc vulputate tincidunt at suscipit nunc.",
	}
	reviews := []*models.Review{review, review, review}

	// Create Reviews if none in the database
	result := db.Take(&models.Review{})
	if errors.Is(result.Error, gorm.ErrRecordNotFound) {
		db.Create(&reviews)
	}
}

func CreateInitialData(db *gorm.DB) {
	createSuperUser(db)
	createAutioneer(db)
	reviewer := createReviewer(db)
	createReviews(db, reviewer.ID)
}
