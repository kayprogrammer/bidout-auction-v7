package schemas

import "github.com/kayprogrammer/bidout-auction-v7/models"

// REQUEST BODY SCHEMAS


// RESPONSE BODY SCHEMAS
type ListingsResponseSchemas struct {
	ResponseSchema
	Data					[]models.Listing	`json:"data"`
}

type ListingDetailResponseDataSchema struct {
	Listing					models.Listing		`json:"listing"`
	RelatedListings			[]models.Listing	`json:"related_listings"`
}

type ListingDetailResponseSchema struct {
	ResponseSchema
	Data					ListingDetailResponseDataSchema		`json:"data"`
}