package tests

import (
	"time"
	"github.com/kayprogrammer/bidout-auction-v7/models"
	auth "github.com/kayprogrammer/bidout-auction-v7/authentication"
	"gorm.io/gorm"
	uuid "github.com/satori/go.uuid"
	"github.com/shopspring/decimal"
)

var truth = true

func CreateTestUser(db *gorm.DB) models.User {
	user := models.User{
		FirstName: "Test",
		LastName: "User",
		Email: "testuser@example.com",
		Password: "testpassword",
	}
	db.FirstOrCreate(&user, models.User{Email: user.Email})
	return user
}

func CreateTestVerifiedUser(db *gorm.DB) models.User {
	user := models.User{
		FirstName: "Test",
		LastName: "Verified",
		Email: "testverifieduser@example.com",
		Password: "testpassword",
		IsEmailVerified: &truth,
	}
	db.FirstOrCreate(&user, models.User{Email: user.Email})
	return user
}

func CreateAnotherTestVerifiedUser(db *gorm.DB) models.User {
	user := models.User{
		FirstName: "AnotherTest",
		LastName: "UserVerified",
		Email: "anothertestverifieduser@example.com",
		Password: "testpassword",
		IsEmailVerified: &truth,
	}
	db.FirstOrCreate(&user, models.User{Email: user.Email})
	return user
}

func CreateJwt(db *gorm.DB, userId uuid.UUID) models.Jwt {
	access := auth.GenerateAccessToken(userId)
	refresh := auth.GenerateRefreshToken()
	jwt := models.Jwt{UserId: userId, Access: access, Refresh: refresh}
	db.FirstOrCreate(&jwt, models.Jwt{UserId: userId})
	return jwt
}

func CreateListing(db *gorm.DB) models.Listing {
	auctioneer := CreateTestVerifiedUser(db)

	category := models.Category{Name: "TestCategory"}
	db.Create(&category)

	image := models.File{ResourceType: "image/png"}
	db.Create(&image)

	listing := models.Listing{
		AuctioneerId: auctioneer.ID,
		Name: "New Listing",
		Desc: "New description",
		CategoryId: &category.ID,
		Price: decimal.NewFromInt(int64(1000)),
		ClosingDate: time.Now().Add(24 * time.Hour),
		ImageId: image.ID,
	}

	db.Create(&listing)
	return listing
}