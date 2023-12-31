package routes

import (
	"github.com/gofiber/fiber/v2"
	uuid "github.com/satori/go.uuid"

	auth "github.com/kayprogrammer/bidout-auction-v7/authentication"
	"github.com/kayprogrammer/bidout-auction-v7/models"
	"github.com/kayprogrammer/bidout-auction-v7/schemas"
	"github.com/kayprogrammer/bidout-auction-v7/senders"
	"github.com/kayprogrammer/bidout-auction-v7/utils"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// @Summary Register a new user
// @Description This endpoint registers new users into our application.
// @Tags Auth
// @Param user body models.User true "User object"
// @Success 201 {object} schemas.RegisterResponseSchema
// @Failure 422 {object} utils.ErrorResponse
// @Router /auth/register [post]
func Register(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	validator := utils.Validator()

	user := models.User{}

	// Validate request
	if errCode, errData := DecodeJSONBody(c, &user); errData != nil {
		return c.Status(errCode).JSON(errData)
	}
	if err := validator.Validate(user); err != nil {
		return c.Status(422).JSON(err)
	}

	db.Take(&user, models.User{Email: user.Email})
	if user.ID != uuid.Nil {
		data := map[string]string{
			"email": "Email already registered!",
		}
		return c.Status(422).JSON(utils.ErrorResponse{Message: "Invalid Entry", Data: &data}.Init())
	}

	// Create User
	db.Create(&user)

	// Send Email
	go senders.SendEmail(c.Locals("env"), db, user, "activate")

	response := schemas.RegisterResponseSchema{
		ResponseSchema: schemas.ResponseSchema{Message: "Registration successful"}.Init(),
		Data:           schemas.EmailRequestSchema{Email: user.Email},
	}
	return c.Status(201).JSON(response)
}

// @Summary Verify a user's email
// @Description This endpoint verifies a user's email.
// @Tags Auth
// @Param verify_email body schemas.VerifyEmailRequestSchema true "Verify Email object"
// @Success 200 {object} schemas.ResponseSchema
// @Failure 422 {object} utils.ErrorResponse
// @Router /auth/verify-email [post]
func VerifyEmail(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	validator := utils.Validator()

	verifyEmail := schemas.VerifyEmailRequestSchema{}

	// Validate request
	if errCode, errData := DecodeJSONBody(c, &verifyEmail); errData != nil {
		return c.Status(errCode).JSON(errData)
	}
	if err := validator.Validate(verifyEmail); err != nil {
		return c.Status(422).JSON(err)
	}

	user := models.User{Email: verifyEmail.Email}
	db.Take(&user, user)
	if user.ID == uuid.Nil {
		return c.Status(404).JSON(utils.ErrorResponse{Message: "Incorrect Email"}.Init())
	}

	if *user.IsEmailVerified {
		return c.Status(200).JSON(schemas.ResponseSchema{Message: "Email already verified"}.Init())
	}

	otp := models.Otp{UserId: user.ID}
	db.Take(&otp, otp)
	if otp.ID == uuid.Nil || *otp.Code != verifyEmail.Otp {
		return c.Status(404).JSON(utils.ErrorResponse{Message: "Incorrect Otp"}.Init())
	}

	if otp.CheckExpiration() {
		return c.Status(400).JSON(utils.ErrorResponse{Message: "Expired Otp"}.Init())
	}

	// Update User
	*user.IsEmailVerified = true
	db.Save(&user)

	// Send Welcome Email
	go senders.SendEmail(c.Locals("env"), db, user, "welcome")

	response := schemas.ResponseSchema{Message: "Account verification successful"}.Init()
	return c.Status(200).JSON(response)
}

// @Summary Resend Verification Email
// @Description This endpoint resends new otp to the user's email.
// @Tags Auth
// @Param email body schemas.EmailRequestSchema true "Email object"
// @Success 200 {object} schemas.ResponseSchema
// @Failure 422 {object} utils.ErrorResponse
// @Router /auth/resend-verification-email [post]
func ResendVerificationEmail(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	validator := utils.Validator()

	emailSchema := schemas.EmailRequestSchema{}

	// Validate request
	if errCode, errData := DecodeJSONBody(c, &emailSchema); errData != nil {
		return c.Status(errCode).JSON(errData)
	}
	if err := validator.Validate(emailSchema); err != nil {
		return c.Status(422).JSON(err)
	}

	user := models.User{Email: emailSchema.Email}
	db.Take(&user, user)
	if user.ID == uuid.Nil {
		return c.Status(404).JSON(utils.ErrorResponse{Message: "Incorrect Email"}.Init())
	}

	if *user.IsEmailVerified {
		return c.Status(200).JSON(schemas.ResponseSchema{Message: "Email already verified"}.Init())
	}

	// Send Email
	go senders.SendEmail(c.Locals("env"), db, user, "activate")

	response := schemas.ResponseSchema{Message: "Verification email sent"}.Init()
	return c.Status(200).JSON(response)
}

// @Summary Send Password Reset Otp
// @Description This endpoint sends new password reset otp to the user's email.
// @Tags Auth
// @Param email body schemas.EmailRequestSchema true "Email object"
// @Success 200 {object} schemas.ResponseSchema
// @Failure 422 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /auth/send-password-reset-otp [post]
func SendPasswordResetOtp(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	validator := utils.Validator()

	emailSchema := schemas.EmailRequestSchema{}

	// Validate request
	if errCode, errData := DecodeJSONBody(c, &emailSchema); errData != nil {
		return c.Status(errCode).JSON(errData)
	}
	if err := validator.Validate(emailSchema); err != nil {
		return c.Status(422).JSON(err)
	}

	user := models.User{Email: emailSchema.Email}
	db.Take(&user, user)
	if user.ID == uuid.Nil {
		return c.Status(404).JSON(utils.ErrorResponse{Message: "Incorrect Email"}.Init())
	}

	// Send Email
	go senders.SendEmail(c.Locals("env"), db, user, "reset")

	response := schemas.ResponseSchema{Message: "Password otp sent"}.Init()
	return c.Status(200).JSON(response)
}

// @Summary Set New Password
// @Description This endpoint verifies the password reset otp.
// @Tags Auth
// @Param email body schemas.SetNewPasswordSchema true "Password reset object"
// @Success 200 {object} schemas.ResponseSchema
// @Failure 422 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Router /auth/set-new-password [post]
func SetNewPassword(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	validator := utils.Validator()

	passwordResetSchema := schemas.SetNewPasswordSchema{}

	// Validate request
	if errCode, errData := DecodeJSONBody(c, &passwordResetSchema); errData != nil {
		return c.Status(errCode).JSON(errData)
	}
	if err := validator.Validate(passwordResetSchema); err != nil {
		return c.Status(422).JSON(err)
	}

	user := models.User{Email: passwordResetSchema.Email}
	db.Take(&user, user)
	if user.ID == uuid.Nil {
		return c.Status(404).JSON(utils.ErrorResponse{Message: "Incorrect Email"}.Init())
	}

	otp := models.Otp{UserId: user.ID}
	db.Take(&otp, otp)
	if otp.ID == uuid.Nil || *otp.Code != passwordResetSchema.Otp {
		return c.Status(404).JSON(utils.ErrorResponse{Message: "Incorrect Otp"}.Init())
	}

	if otp.CheckExpiration() {
		return c.Status(400).JSON(utils.ErrorResponse{Message: "Expired Otp"}.Init())
	}

	// Set Password
	user.Password = passwordResetSchema.Password
	db.Save(&user)

	// Send Email
	go senders.SendEmail(c.Locals("env"), db, user, "reset-success")

	response := schemas.ResponseSchema{Message: "Password reset successful"}.Init()
	return c.Status(200).JSON(response)
}

// @Summary Login a user
// @Description This endpoint generates new access and refresh tokens for authentication
// @Tags Auth
// @Param user body schemas.LoginSchema true "User login"
// @Success 201 {object} schemas.ResponseSchema
// @Failure 422 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Security GuestUserAuth
// @Router /auth/login [post]
func Login(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	validator := utils.Validator()

	userLoginSchema := schemas.LoginSchema{}

	// Validate request
	if errCode, errData := DecodeJSONBody(c, &userLoginSchema); errData != nil {
		return c.Status(errCode).JSON(errData)
	}
	if err := validator.Validate(userLoginSchema); err != nil {
		return c.Status(422).JSON(err)
	}

	user := models.User{Email: userLoginSchema.Email}
	db.Take(&user, user)
	if user.ID == uuid.Nil {
		return c.Status(401).JSON(utils.ErrorResponse{Message: "Invalid Credentials"}.Init())
	}
	if !utils.CheckPasswordHash(userLoginSchema.Password, user.Password) {
		return c.Status(401).JSON(utils.ErrorResponse{Message: "Invalid Credentials"}.Init())
	}

	if !*user.IsEmailVerified {
		return c.Status(401).JSON(utils.ErrorResponse{Message: "Verify your email first"}.Init())
	}

	// Create Auth Tokens
	access := auth.GenerateAccessToken(user.ID)
	refresh := auth.GenerateRefreshToken()
	jwt := models.Jwt{UserId: user.ID, Access: access, Refresh: refresh}
	jwtToDelete := models.Jwt{UserId: user.ID}
	db.Where(jwtToDelete).Delete(&jwtToDelete) // Delete existing jwt
	db.Create(&jwt)

	// Move all guest user watchlists to the authenticated user watchlists
	client := GetClient(c)
	if (client != nil) && (client.Type == "guest") {
		// clientId := client.
		watchlists := []models.Watchlist{}
		db.Find(&watchlists, models.Watchlist{GuestUserId: &client.ID})
		if len(watchlists) > 0 {
			watchlistsToCreate := []models.Watchlist{}
			for _, wl := range watchlists {
				watchlist := models.Watchlist{UserId: &user.ID, ListingId: wl.ListingId}
				watchlistsToCreate = append(watchlistsToCreate, watchlist)
			}
			db.Clauses(clause.OnConflict{DoNothing: true}).Create(&watchlistsToCreate)
		}
		db.Delete(&models.GuestUser{}, client.ID)
	}
	response := schemas.LoginResponseSchema{
		ResponseSchema: schemas.ResponseSchema{Message: "Login successful"}.Init(),
		Data:           schemas.TokensResponseSchema{Access: access, Refresh: refresh},
	}
	return c.Status(201).JSON(response)
}

// @Summary Refresh tokens
// @Description This endpoint refresh tokens by generating new access and refresh tokens for a user
// @Tags Auth
// @Param refresh body schemas.RefreshTokenSchema true "Refresh token"
// @Success 201 {object} schemas.ResponseSchema
// @Failure 422 {object} utils.ErrorResponse
// @Failure 404 {object} utils.ErrorResponse
// @Failure 401 {object} utils.ErrorResponse
// @Router /auth/refresh [post]
func Refresh(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	validator := utils.Validator()

	refreshTokenSchema := schemas.RefreshTokenSchema{}

	// Validate request
	if errCode, errData := DecodeJSONBody(c, &refreshTokenSchema); errData != nil {
		return c.Status(errCode).JSON(errData)
	}
	if err := validator.Validate(refreshTokenSchema); err != nil {
		return c.Status(422).JSON(err)
	}

	token := refreshTokenSchema.Refresh
	jwt := models.Jwt{Refresh: token}
	db.Take(&jwt, jwt)
	if jwt.ID == uuid.Nil {
		return c.Status(404).JSON(utils.ErrorResponse{Message: "Refresh token does not exist"}.Init())
	}

	if !auth.DecodeRefreshToken(token) {
		return c.Status(401).JSON(utils.ErrorResponse{Message: "Refresh token is invalid or expired"}.Init())
	}

	// Create and Update Auth Tokens
	access := auth.GenerateAccessToken(jwt.UserId)
	refresh := auth.GenerateRefreshToken()
	jwt.Access = access
	jwt.Refresh = refresh
	db.Save(&jwt)

	response := schemas.LoginResponseSchema{
		ResponseSchema: schemas.ResponseSchema{Message: "Tokens refresh successful"}.Init(),
		Data:           schemas.TokensResponseSchema{Access: access, Refresh: refresh},
	}
	return c.Status(201).JSON(response)
}

// @Summary Logout a user
// @Description This endpoint logs a user out from our application
// @Tags Auth
// @Success 200 {object} schemas.ResponseSchema
// @Failure 401 {object} utils.ErrorResponse
// @Router /auth/logout [get]
// @Security BearerAuth
func Logout(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	user := c.Locals("user").(*models.User)

	jwt := models.Jwt{UserId: user.ID}
	db.Where(jwt).Delete(&jwt) // Delete jwt

	response := schemas.ResponseSchema{Message: "Logout successful"}.Init()
	return c.Status(200).JSON(response)
}
