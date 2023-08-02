package tests

import (
	"fmt"
	"testing"
	"encoding/json"
	"net/http/httptest"

	"github.com/gofiber/fiber/v2"
	"github.com/kayprogrammer/bidout-auction-v7/models"
	"github.com/kayprogrammer/bidout-auction-v7/schemas"
	auth "github.com/kayprogrammer/bidout-auction-v7/authentication"
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
		assert.Equal(t, 201, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Registration successful", body["message"])
		expectedData := make(map[string]interface{})
		expectedData["email"] = validEmail
		assert.Equal(t, expectedData, body["data"].(map[string]interface{}))

		// Verify that a user with the same email cannot be registered again
		res = ProcessTestBody(t, app, url, "POST", userData)
		assert.Equal(t, 422, res.StatusCode)

		// Parse and assert body
		body = ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "Invalid Entry", body["message"])
		expectedData = make(map[string]interface{})
		expectedData["email"] = "Email already registered!"
		assert.Equal(t, expectedData, body["data"].(map[string]interface{}))
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
		assert.Equal(t, 404, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "Incorrect Otp", body["message"])

		// Verify that the email verification succeeds with a valid otp
		realOtp := models.Otp{UserId: user.ID}
		db.Create(&realOtp)
		emailOtpData.Otp = *realOtp.Code
		res = ProcessTestBody(t, app, url, "POST", emailOtpData)
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body = ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Account verification successful", body["message"])
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
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Verification email sent", body["message"])

		// Verify that a verified user cannot get a new email
		*user.IsEmailVerified = true
		db.Save(&user)
		res = ProcessTestBody(t, app, url, "POST", emailData)
		
		assert.Equal(t, 200, res.StatusCode)
		// Parse and assert body
		body = ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Email already verified", body["message"])

		// Verify that an error is raised when attempting to resend the verification email for a user that doesn't exist
		emailData.Email = "invalid@example.com"
		res = ProcessTestBody(t, app, url, "POST", emailData)
		
		assert.Equal(t, 404, res.StatusCode)
		// Parse and assert body
		body = ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "Incorrect Email", body["message"])
	})
}

func sendPasswordResetOtp(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	t.Run("Send Password Reset Otp", func(t *testing.T) {

		user := CreateTestVerifiedUser(db)

		url := fmt.Sprintf("%s/send-password-reset-otp", baseUrl)
		emailData := schemas.EmailRequestSchema{
			Email: user.Email,
		}

		emailSenderMock := new(MockEmailSender)
		emailSenderMock.On("sendEmail", db, user, "reset").Return(nil)

		res := ProcessTestBody(t, app, url, "POST", emailData)

		// Verify that an unverified user can get a new email
		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Password otp sent", body["message"])

		// Verify that an error is raised when attempting to send password reset email for a user that doesn't exist
		emailData.Email = "invalid@example.com"
		res = ProcessTestBody(t, app, url, "POST", emailData)
		
		assert.Equal(t, 404, res.StatusCode)
		// Parse and assert body
		body = ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "Incorrect Email", body["message"])
	})
}

func setNewPassword(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	// Drop and Create User Table since the previous test uses the verified_user it...
	DropSingleTable(db, models.User{})
	CreateSingleTable(db, models.User{})

	t.Run("Set New Password", func(t *testing.T) {
		user := CreateTestVerifiedUser(db)

		url := fmt.Sprintf("%s/set-new-password", baseUrl)
		passwordResetData := schemas.SetNewPasswordSchema{
			VerifyEmailRequestSchema: schemas.VerifyEmailRequestSchema{
				Email: "invalid@example.com", // Invalid otp
				Otp: 11111, // Invalid otp
			},
			Password: "newpassword",
		}

		emailSenderMock := new(MockEmailSender)
		emailSenderMock.On("sendEmail", db, user, "reset-success").Return(nil)

		res := ProcessTestBody(t, app, url, "POST", passwordResetData)

		// Verify that the request fails with incorrect email
		// Assert Status code
		assert.Equal(t, 404, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "Incorrect Email", body["message"])

		// Verify that the request fails with incorrect otp
		passwordResetData.Email = user.Email
		res = ProcessTestBody(t, app, url, "POST", passwordResetData)
		// Assert Status code
		assert.Equal(t, 404, res.StatusCode)

		// Parse and assert body
		body = ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "Incorrect Otp", body["message"])

		// Verify that password reset succeeds
		realOtp := models.Otp{UserId: user.ID}
		db.Create(&realOtp)
		passwordResetData.Otp = *realOtp.Code
		res = ProcessTestBody(t, app, url, "POST", passwordResetData)

		// Assert response
		assert.Equal(t, 200, res.StatusCode)
		// Parse and assert body
		body = ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Password reset successful", body["message"])
	})
}

func login(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	t.Run("Login", func(t *testing.T) {
		user := CreateTestUser(db)

		url := fmt.Sprintf("%s/login", baseUrl)
		loginData := schemas.LoginSchema{
			Email: "invalid@example.com", // Invalid email
			Password: "invalidpassword",
		}

		res := ProcessTestBody(t, app, url, "POST", loginData)

		// # Test for invalid credentials
		// Assert Status code
		assert.Equal(t, 401, res.StatusCode)
		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "Invalid Credentials", body["message"])

		// Test for unverified credentials (email)
		loginData.Email = user.Email
		loginData.Password = "testpassword"
		res = ProcessTestBody(t, app, url, "POST", loginData)
		// Assert Status code
		assert.Equal(t, 401, res.StatusCode)
		// Parse and assert body
		body = ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "Verify your email first", body["message"])

		// Test for valid credentials and verified email address
		*user.IsEmailVerified = true
		db.Save(&user)
		res = ProcessTestBody(t, app, url, "POST", loginData)
		// Assert response
		assert.Equal(t, 201, res.StatusCode)
		// Parse and assert body
		body = ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Login successful", body["message"])
		jwt := models.Jwt{}
		db.Find(&jwt,"user_id = ?", user.ID)
		expectedData := map[string]string{
			"access": jwt.Access,
			"refresh": jwt.Refresh,
		}
		data, _ := json.Marshal(body["data"])
		expectedDataJson, _ := json.Marshal(expectedData)
		assert.Equal(t, expectedDataJson, data)
	})
}

func refresh(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	// Drop and Create User Table since the previous test uses the verified_user it...
	DropSingleTable(db, models.User{})
	CreateSingleTable(db, models.User{})

	t.Run("Refresh", func(t *testing.T) {
		user := CreateTestVerifiedUser(db)

		url := fmt.Sprintf("%s/refresh", baseUrl)
		refreshTokenData := schemas.RefreshTokenSchema{
			Refresh: "invalid@example.com", // non-exisitent token
		}

		res := ProcessTestBody(t, app, url, "POST", refreshTokenData)

		// # Test for invalid refresh token (not found)
		// Assert Status code
		assert.Equal(t, 404, res.StatusCode)
		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "Refresh token does not exist", body["message"])

		// Test for invalid refresh token (invalid or expired)
		jwt := models.Jwt{UserId: user.ID, Access: "invalid_access", Refresh: "invalid_refresh"}
		db.Create(&jwt)
		refreshTokenData.Refresh = jwt.Refresh
		res = ProcessTestBody(t, app, url, "POST", refreshTokenData)
		// Assert Status code
		assert.Equal(t, 401, res.StatusCode)
		// Parse and assert body
		body = ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "Refresh token is invalid or expired", body["message"])

		// Test for valid refresh token
		jwt.Refresh = auth.GenerateRefreshToken()
		db.Save(&jwt)
		refreshTokenData.Refresh = jwt.Refresh
		res = ProcessTestBody(t, app, url, "POST", refreshTokenData)
		// Assert response
		assert.Equal(t, 201, res.StatusCode)
		// Parse and assert body
		body = ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Tokens refresh successful", body["message"])
		jwt = models.Jwt{}
		db.Find(&jwt,"user_id = ?", user.ID)
		expectedData := map[string]string{
			"access": jwt.Access,
			"refresh": jwt.Refresh,
		}
		data, _ := json.Marshal(body["data"])
		expectedDataJson, _ := json.Marshal(expectedData)
		assert.Equal(t, expectedDataJson, data)
	})
}

func logout(t *testing.T, app *fiber.App, db *gorm.DB, baseUrl string) {
	// Drop and Create User Table since the previous test uses the verified_user it...
	DropSingleTable(db, models.User{})
	CreateSingleTable(db, models.User{})
	t.Run("Logout", func(t *testing.T) {
		url := fmt.Sprintf("%s/logout", baseUrl)
		req := httptest.NewRequest("GET", url, nil)
		req.Header.Set("Authorization", "Bearer invalid_token")
		res, _ := app.Test(req)

		// Ensures an unauthorized user cannot log out
		// Assert Status code
		assert.Equal(t, 401, res.StatusCode)

		// Parse and assert body
		body := ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "failure", body["status"])
		assert.Equal(t, "Auth Token is Invalid or Expired!", body["message"])

		// Ensures an authorized user can log out
		req = httptest.NewRequest("GET", url, nil)
		user := CreateTestVerifiedUser(db)
		jwt := CreateJwt(db, user.ID)
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", jwt.Access))
		res, _ = app.Test(req)

		// Assert Status code
		assert.Equal(t, 200, res.StatusCode)

		// Parse and assert body
		body = ParseResponseBody(t, res.Body).(map[string]interface{})
		assert.Equal(t, "success", body["status"])
		assert.Equal(t, "Logout successful", body["message"])
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
	sendPasswordResetOtp(t, app, db, BASEURL)
	setNewPassword(t, app, db, BASEURL)
	login(t, app, db, BASEURL)
	logout(t, app, db, BASEURL)
	refresh(t, app, db, BASEURL)

	// Drop Tables and Close Connectiom
	CloseTestDatabase(db)
}