// config/config.go

package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Configuration holds the application configuration loaded from environment variables
type Configuration struct {
	CloudinaryCloudName       string
	CloudinaryAPIKey          string
	CloudinaryAPISecret       string
	ProjectName               string
	Debug                     string
	EmailOTPExpireSeconds     int64
	AccessTokenExpireMinutes  int
	RefreshTokenExpireMinutes int
	SecretKey                 string
	FrontendURL               string
	FirstSuperuserEmail       string
	FirstSuperuserPassword    string
	FirstAuctioneerEmail      string
	FirstAuctioneerPassword   string
	FirstReviewerEmail        string
	FirstReviewerPassword     string
	PostgresUser              string
	PostgresPassword          string
	PostgresServer            string
	PostgresPort              string
	PostgresDB                string
	MailSenderEmail           string
	MailSenderPassword        string
	MailSenderHost            string
	MailSenderPort            string
	CORSAllowedOrigins        string
}

var config *Configuration

func init() {
	// Load environment variables from the .env file (if it exists) into the environment
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}

	// Convert string-based numeric variables to their respective types
	emailOTPExpireSeconds, _ := strconv.ParseInt(os.Getenv("EMAIL_OTP_EXPIRE_SECONDS"), 10, 64)
	accessTokenExpireMinutes, _ := strconv.Atoi(os.Getenv("ACCESS_TOKEN_EXPIRE_MINUTES"))
	refreshTokenExpireMinutes, _ := strconv.Atoi(os.Getenv("REFRESH_TOKEN_EXPIRE_MINUTES"))

	config = &Configuration{
		CloudinaryCloudName:       os.Getenv("CLOUDINARY_CLOUD_NAME"),
		CloudinaryAPIKey:          os.Getenv("CLOUDINARY_API_KEY"),
		CloudinaryAPISecret:       os.Getenv("CLOUDINARY_API_SECRET"),
		ProjectName:               os.Getenv("PROJECT_NAME"),
		Debug:                     os.Getenv("DEBUG"),
		EmailOTPExpireSeconds:     emailOTPExpireSeconds,
		AccessTokenExpireMinutes:  accessTokenExpireMinutes,
		RefreshTokenExpireMinutes: refreshTokenExpireMinutes,
		SecretKey:                 os.Getenv("SECRET_KEY"),
		FrontendURL:               os.Getenv("FRONTEND_URL"),
		FirstSuperuserEmail:       os.Getenv("FIRST_SUPERUSER_EMAIL"),
		FirstSuperuserPassword:    os.Getenv("FIRST_SUPERUSER_PASSWORD"),
		FirstAuctioneerEmail:      os.Getenv("FIRST_AUCTIONEER_EMAIL"),
		FirstAuctioneerPassword:   os.Getenv("FIRST_AUCTIONEER_PASSWORD"),
		FirstReviewerEmail:        os.Getenv("FIRST_REVIEWER_EMAIL"),
		FirstReviewerPassword:     os.Getenv("FIRST_REVIEWER_PASSWORD"),
		PostgresUser:              os.Getenv("POSTGRES_USER"),
		PostgresPassword:          os.Getenv("POSTGRES_PASSWORD"),
		PostgresServer:            os.Getenv("POSTGRES_SERVER"),
		PostgresPort:              os.Getenv("POSTGRES_PORT"),
		PostgresDB:                os.Getenv("POSTGRES_DB"),
		MailSenderEmail:           os.Getenv("MAIL_SENDER_EMAIL"),
		MailSenderPassword:        os.Getenv("MAIL_SENDER_PASSWORD"),
		MailSenderHost:            os.Getenv("MAIL_SENDER_HOST"),
		MailSenderPort:            os.Getenv("MAIL_SENDER_PORT"),
		CORSAllowedOrigins:        os.Getenv("CORS_ALLOWED_ORIGINS"),
	}
}

// GetConfig returns the application configuration
func GetConfig() *Configuration {
	return config
}
