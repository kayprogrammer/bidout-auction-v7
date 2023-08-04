package schemas

import (
	"github.com/satori/go.uuid"
	"github.com/kayprogrammer/bidout-auction-v7/models"
)

// REQUEST BODY SCHEMAS
type AddOrRemoveWatchlistSchema struct {
	Slug					string					`json:"slug" validate:"required" example:"listing_slug"`
}

// RESPONSE BODY SCHEMAS
type ListingsResponseSchema struct {
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

type AddOrRemoveWatchlistResponseDataSchema struct {
	GuestUserId				*uuid.UUID						`json:"guestuser_id"`
}

type AddOrRemoveWatchlistResponseSchema struct {
	ResponseSchema
	Data					AddOrRemoveWatchlistResponseDataSchema		`json:"data"`
}