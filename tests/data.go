package tests

import (
	"github.com/kayprogrammer/bidout-auction-v7/models"
	"gorm.io/gorm"
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