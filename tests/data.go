package tests

import (
	"github.com/kayprogrammer/bidout-auction-v7/models"
	auth "github.com/kayprogrammer/bidout-auction-v7/authentication"
	"gorm.io/gorm"
	uuid "github.com/satori/go.uuid"
)

var truth = true

func CreateTestUser(db *gorm.DB) models.User {
	user := models.User{
		FirstName: "Test",
		LastName: "User",
		Email: "testuser@example.com",
		Password: "testpassword",
	}
	db.Create(&user)
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
	db.Create(&user)
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
	db.Create(&user)
	return user
}

func CreateJwt(db *gorm.DB, userId uuid.UUID) models.Jwt {
	access := auth.GenerateAccessToken(userId)
	refresh := auth.GenerateRefreshToken()
	jwt := models.Jwt{UserId: userId, Access: access, Refresh: refresh}
	db.Create(&jwt)
	return jwt
}