package tests

import (
	"fmt"
	"testing"

	"github.com/gofiber/fiber/v2"
	"github.com/kayprogrammer/bidout-auction-v7/models"
	"github.com/kayprogrammer/bidout-auction-v7/schemas"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"gorm.io/gorm"
)

type MockEmailSender struct {
    mock.Mock
}

func (m *MockEmailSender) SendEmail(db *gorm.DB, user models.User, emailType string) {
    m.Called(db, user, emailType)
}

func register(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	t.Run("Register User", func(t *testing.T) {
		url := fmt.Sprintf("%s/register", baseUrl)
		validEmail := "testregisteruser@email.com"
		userData := models.User{
			FirstName: "TestRegister", 
			LastName: "User", 
			Email: validEmail,
			Password: "testregisteruserpassword",
			TermsAgreement: true,
		}

		emailSenderMock := new(MockEmailSender)
		emailSenderMock.On("sendEmail", db, userData, "activate").Return(nil)

		res := ProcessTestBody(t, app, url, "POST", userData)

		// Assert Status code
		assert.Equal(t, res.StatusCode, 201)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, body["status"], "success")
		assert.Equal(t, body["message"], "Registration successful")
		expectedData := make(map[string]interface{})
		expectedData["email"] = validEmail
		assert.Equal(t, body["data"].(map[string]interface{}), expectedData)

		// Verify that a user with the same email cannot be registered again
		res = ProcessTestBody(t, app, url, "POST", userData)
		assert.Equal(t, res.StatusCode, 422)

		// Parse and assert body
		body = ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, body["status"], "failure")
		assert.Equal(t, body["message"], "Invalid Entry")
		expectedData = make(map[string]interface{})
		expectedData["email"] = "Email already registered!"
		assert.Equal(t, body["data"].(map[string]interface{}), expectedData)
	})
}

func verifyEmail(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	t.Run("Verify Email", func(t *testing.T) {
		user := CreateTestUser(db)
		otp := 1111

		url := fmt.Sprintf("%s/verify-email", baseUrl)
		emailOtpData := schemas.VerifyEmailRequestSchema{
			Email: user.Email,
			Otp: otp,
		}

		emailSenderMock := new(MockEmailSender)
		emailSenderMock.On("sendEmail", db, user, "welcome").Return(nil)

		res := ProcessTestBody(t, app, url, "POST", emailOtpData)

		// Verify that the email verification fails with an invalid otp
		// Assert Status code
		assert.Equal(t, res.StatusCode, 404)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, body["status"], "failure")
		assert.Equal(t, body["message"], "Incorrect Otp")

		// Verify that the email verification succeeds with a valid otp
		realOtp := models.Otp{UserId: user.ID}
		db.Create(&realOtp)
		emailOtpData.Otp = *realOtp.Code
		res = ProcessTestBody(t, app, url, "POST", emailOtpData)
		assert.Equal(t, res.StatusCode, 200)

		// Parse and assert body
		body = ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, body["status"], "success")
		assert.Equal(t, body["message"], "Account verification successful")
	})
}

func resendVerificationEmail(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	t.Run("Resend Verification Email", func(t *testing.T) {
		// Drop and Create User Table since the previous test uses it...
		DropSingleTable(db, models.User{})
		CreateSingleTable(db, models.User{})

		user := CreateTestUser(db)

		url := fmt.Sprintf("%s/resend-verification-email", baseUrl)
		emailData := schemas.EmailRequestSchema{
			Email: user.Email,
		}

		emailSenderMock := new(MockEmailSender)
		emailSenderMock.On("sendEmail", db, user, "activate").Return(nil)

		res := ProcessTestBody(t, app, url, "POST", emailData)

		// Verify that an unverified user can get a new email
		// Assert Status code
		assert.Equal(t, res.StatusCode, 200)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, body["status"], "success")
		assert.Equal(t, body["message"], "Verification email sent")

		// Verify that a verified user cannot get a new email
		user.IsEmailVerified = &truth
		db.Save(&user)
		res = ProcessTestBody(t, app, url, "POST", emailData)
		
		assert.Equal(t, res.StatusCode, 200)
		// Parse and assert body
		body = ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, body["status"], "success")
		assert.Equal(t, body["message"], "Email already verified")

		// Verify that an error is raised when attempting to resend the verification email for a user that doesn't exist
		emailData.Email = "invalid@example.com"
		res = ProcessTestBody(t, app, url, "POST", emailData)
		
		assert.Equal(t, res.StatusCode, 404)
		// Parse and assert body
		body = ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, body["status"], "failure")
		assert.Equal(t, body["message"], "Incorrect Email")
	})
}


func TestAuth(t *testing.T) {
	app := fiber.New()
	db := Setup(t, app)
	BASEURL := "/api/v7/auth"

	// Run Auth Endpoint Tests
	register(t, app, db, BASEURL)
	verifyEmail(t, app, db, BASEURL)
	resendVerificationEmail(t, app, db, BASEURL)

	// Drop Tables and Close Connectiom
	CloseTestDatabase(db)
}