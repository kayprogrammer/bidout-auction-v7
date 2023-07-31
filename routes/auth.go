package routes

import (
	"log"
	"github.com/gofiber/fiber/v2"
	"github.com/satori/go.uuid"

	"github.com/kayprogrammer/bidout-auction-v7/models"
	"github.com/kayprogrammer/bidout-auction-v7/schemas"
	"github.com/kayprogrammer/bidout-auction-v7/senders"
	"github.com/kayprogrammer/bidout-auction-v7/utils"
	"gorm.io/gorm"

)

var truth = true

type EmailSender struct{}

func (es *EmailSender) SendEmail(db *gorm.DB, user models.User, emailType string) {
    // Implementation of sending an actual email using your preferred email library.
    // Replace this with your real email sending code.
    senders.SendEmail(db, user, emailType)
}

// @Summary Register a new user
// @Description This endpoint registers new users into our application.
// @Tags Auth
// @Param user body models.User true "User object"
// @Success 201 {object} schemas.RegisterResponseSchema
// @Failure 422 {object} utils.ErrorResponse
// @Router /api/v7/auth/register [post]
func Register(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
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
	emailSender := &EmailSender{}
	go emailSender.SendEmail(db, user, "activate")

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
// @Router /api/v7/auth/verify-email [post]
func VerifyEmail(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	validator := utils.Validator()

	verifyEmail := schemas.VerifyEmailRequestSchema{}
	c.BodyParser(&verifyEmail)

	// Validate request
	if err := validator.Validate(verifyEmail); err != nil {
		return c.Status(422).JSON(err)
	}

	user := models.User{}
	db.Find(&user,"email = ?", verifyEmail.Email)
	if user.ID == uuid.Nil {
		return c.Status(404).JSON(utils.ErrorResponse{Message: "Incorrect Email"}.Init())
	}
	log.Println(user.IsEmailVerified)

	if *user.IsEmailVerified {
		return c.Status(200).JSON(schemas.ResponseSchema{Message: "Email already verified"}.Init())
	}

	otp := models.Otp{}
	db.Find(&otp,"user_id = ?", user.ID)
	if otp.ID == uuid.Nil || *otp.Code !=  verifyEmail.Otp {
		return c.Status(404).JSON(utils.ErrorResponse{Message: "Incorrect Otp"}.Init())
	}

	if otp.CheckExpiration() {
		return c.Status(400).JSON(utils.ErrorResponse{Message: "Expired Otp"}.Init())
	}

	// Update User
	user.IsEmailVerified = &truth
	db.Save(&user)

	// Send Welcome Email
	emailSender := &EmailSender{}
	go emailSender.SendEmail(db, user, "welcome")

	response := schemas.ResponseSchema{Message: "Account verification successful"}.Init()
	return c.Status(200).JSON(response)
}


// @Summary Resend Verification Email
// @Description This endpoint resends new otp to the user's email.
// @Tags Auth
// @Param email body schemas.EmailRequestSchema true "Email object"
// @Success 200 {object} schemas.ResponseSchema
// @Failure 422 {object} utils.ErrorResponse
// @Router /api/v7/auth/resend-verification-email [post]
func ResendVerificationEmail(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	validator := utils.Validator()

	emailSchema := schemas.EmailRequestSchema{}
	c.BodyParser(&emailSchema)

	// Validate request
	if err := validator.Validate(emailSchema); err != nil {
		return c.Status(422).JSON(err)
	}

	user := models.User{}
	db.Find(&user,"email = ?", emailSchema.Email)
	if user.ID == uuid.Nil {
		return c.Status(404).JSON(utils.ErrorResponse{Message: "Incorrect Email"}.Init())
	}

	if *user.IsEmailVerified {
		return c.Status(200).JSON(schemas.ResponseSchema{Message: "Email already verified"}.Init())
	}

	// Send Email
	emailSender := &EmailSender{}
	go emailSender.SendEmail(db, user, "activate")

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
// @Router /api/v7/auth/send-password-reset-otp [post]
func SendPasswordResetOtp(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	validator := utils.Validator()

	emailSchema := schemas.EmailRequestSchema{}
	c.BodyParser(&emailSchema)

	// Validate request
	if err := validator.Validate(emailSchema); err != nil {
		return c.Status(422).JSON(err)
	}

	user := models.User{}
	db.Find(&user,"email = ?", emailSchema.Email)
	if user.ID == uuid.Nil {
		return c.Status(404).JSON(utils.ErrorResponse{Message: "Incorrect Email"}.Init())
	}

	// Send Email
	emailSender := &EmailSender{}
	go emailSender.SendEmail(db, user, "reset")

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
// @Router /api/v7/auth/set-new-password [post]
func SetNewPassword(c *fiber.Ctx) error {
	db := c.Locals("db").(*gorm.DB)
	validator := utils.Validator()

	passwordResetSchema := schemas.SetNewPasswordSchema{}
	c.BodyParser(&passwordResetSchema)

	// Validate request
	if err := validator.Validate(passwordResetSchema); err != nil {
		return c.Status(422).JSON(err)
	}

	user := models.User{}
	db.Find(&user,"email = ?", passwordResetSchema.Email)
	if user.ID == uuid.Nil {
		return c.Status(404).JSON(utils.ErrorResponse{Message: "Incorrect Email"}.Init())
	}

	otp := models.Otp{}
	db.Find(&otp,"user_id = ?", user.ID)
	if otp.ID == uuid.Nil || *otp.Code !=  passwordResetSchema.Otp {
		return c.Status(404).JSON(utils.ErrorResponse{Message: "Incorrect Otp"}.Init())
	}

	if otp.CheckExpiration() {
		return c.Status(400).JSON(utils.ErrorResponse{Message: "Expired Otp"}.Init())
	}

	// Set Password
	user.Password = passwordResetSchema.Password
	db.Save(&user)

	// Send Email
	emailSender := &EmailSender{}
	go emailSender.SendEmail(db, user, "reset-success")

	response := schemas.ResponseSchema{Message: "Password reset successful"}.Init()
	return c.Status(200).JSON(response)
}