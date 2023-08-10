package schemas

import (
	"github.com/kayprogrammer/bidout-auction-v7/utils"
)

// REQUEST BODY SCHEMAS
type UpdateProfileSchema struct {
	FirstName				string				`json:"first_name" validate:"required" example:"John"`
	LastName				string				`json:"last_name" validate:"required" example:"Doe"`
	FileType				*string				`json:"file_type" example:"image/png"`
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