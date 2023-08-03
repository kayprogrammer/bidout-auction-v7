package schemas

import "github.com/kayprogrammer/bidout-auction-v7/models"

// REQUEST BODY SCHEMAS


// RESPONSE BODY SCHEMAS
type ListingsResponseSchemas struct {
	ResponseSchema
	Data					[]models.Listing	`json:"data"`
}
