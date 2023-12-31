package schemas

import (
	"github.com/kayprogrammer/bidout-auction-v7/models"
	"github.com/kayprogrammer/bidout-auction-v7/utils"
)

// REQUEST BODY SCHEMAS
type UpdateProfileSchema struct {
	FirstName 		string  			`json:"first_name" validate:"required" example:"John"`
	LastName  		string  			`json:"last_name" validate:"required" example:"Doe"`
	FileType  		*string 			`json:"file_type" validate:"omitempty,file_type_validator" example:"image/png"`
}

type CreateListingSchema struct {
	Name        	string	          `json:"name" validate:"required,max=70" example:"Product name"`
	Desc        	string	          `json:"desc" validate:"required" example:"Product description"`
	Category    	string	          `json:"category" validate:"required" example:"category_slug"`
	Price       	float64	 		  `json:"price" validate:"required,gt=0" example:"1000.00"`
	ClosingDate 	string	 		  `json:"closing_date" validate:"required,date,closing_date_validator" example:"2006-01-02T15:04:05.000Z"`
	FileType    	string	          `json:"file_type" validate:"required,file_type_validator" example:"image/jpeg"`
}

type UpdateListingSchema struct {
	Name        *string          `json:"name" validate:"omitempty,max=70" example:"Product name"`
	Desc        *string          `json:"desc" example:"Product description"`
	Category    *string          `json:"category" example:"category_slug"`
	Price       *float64 		 `json:"price" validate:"omitempty,gt=0" example:"1000.00"`
	ClosingDate *string       	 `json:"closing_date" validate:"omitempty,date,closing_date_validator" example:"2006-01-02T15:04:05.000Z"`
	FileType    *string          `json:"file_type" validate:"omitempty,file_type_validator" example:"image/jpeg"`
	Active      *bool            `json:"active" example:"true"`
}

// RESPONSE BODY SCHEMAS
type ProfileResponseDataSchema struct {
	FirstName string  `json:"first_name"`
	LastName  string  `json:"last_name"`
	Avatar    *string `json:"avatar"`
}

type ProfileResponseSchema struct {
	ResponseSchema
	Data ProfileResponseDataSchema `json:"data"`
}

type UpdateProfileResponseDataSchema struct {
	FirstName      string                 `json:"first_name" example:"John"`
	LastName       string                 `json:"last_name" example:"Doe"`
	FileUploadData *utils.SignatureFormat `json:"file_upload_data"`
}

type UpdateProfileResponseSchema struct {
	ResponseSchema
	Data UpdateProfileResponseDataSchema `json:"data"`
}

type CreateListingResponseDataSchema struct {
	models.Listing
	FileUploadData utils.SignatureFormat `json:"file_upload_data"`
}

type CreateListingResponseSchema struct {
	ResponseSchema
	Data CreateListingResponseDataSchema `json:"data"`
}
