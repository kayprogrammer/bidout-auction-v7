package schemas

import (
	"time"
	"github.com/kayprogrammer/bidout-auction-v7/utils"
	"github.com/kayprogrammer/bidout-auction-v7/models"
	"github.com/shopspring/decimal"
)

// REQUEST BODY SCHEMAS
type UpdateProfileSchema struct {
	FirstName				string				`json:"first_name" validate:"required" example:"John"`
	LastName				string				`json:"last_name" validate:"required" example:"Doe"`
	FileType				*string				`json:"file_type" validate:"omitempty,file_type_validator" example:"image/png"`
}

type CreateListingSchema struct {
	Name					string				`json:"name" validate:"required,max=70" example:"Product name"`
	Desc					string				`json:"desc" validate:"required" example:"Product description"`
	Category				string				`json:"category" validate:"required" example:"category_slug"`
	Price					decimal.Decimal		`json:"price" validate:"required" example:"1000.00"`
	ClosingDate				time.Time			`json:"closing_date" validate:"required,closing_date_validator" example:"2006-01-02T15:04:05.000Z"`
	FileType				string				`json:"file_type" validate:"required,file_type_validator" example:"image/jpeg"`
}

type UpdateListingSchema struct {
	Name					*string				`json:"name" validate:"omitempty,max=70" example:"Product name"`
	Desc					*string				`json:"desc" example:"Product description"`
	Category				*string				`json:"category" example:"category_slug"`
	Price					*decimal.Decimal	`json:"price" example:"1000.00"`
	ClosingDate				*time.Time			`json:"closing_date" validate:"omitempty,closing_date_validator" example:"2006-01-02T15:04:05.000Z"`
	FileType				*string				`json:"file_type" validate:"omitempty,file_type_validator" example:"image/jpeg"`
	Active					*bool				`json:"active" example:"true"`
}

// RESPONSE BODY SCHEMAS
type ProfileResponseDataSchema struct {
	FirstName				string				`json:"first_name"`
	LastName				string				`json:"last_name"`
	Avatar					*string				`json:"avatar"`
}

type ProfileResponseSchema struct {
	ResponseSchema
	Data					ProfileResponseDataSchema			`json:"data"`
}

type UpdateProfileResponseDataSchema struct {
	FirstName				string						`json:"first_name" example:"John"`
	LastName				string						`json:"last_name" example:"Doe"`
	FileUploadData			*utils.SignatureFormat		`json:"file_upload_data"`
}

type UpdateProfileResponseSchema struct {
	ResponseSchema
	Data					UpdateProfileResponseDataSchema			`json:"data"`
}

type CreateListingResponseDataSchema struct {
	models.Listing
	FileUploadData			utils.SignatureFormat		`json:"file_upload_data"`
}

type CreateListingResponseSchema struct {
	ResponseSchema
	Data					CreateListingResponseDataSchema			`json:"data"`
}
