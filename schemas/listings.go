package schemas

import (
	"github.com/kayprogrammer/bidout-auction-v7/models"
	"github.com/satori/go.uuid"
)

// REQUEST BODY SCHEMAS
type AddOrRemoveWatchlistSchema struct {
	Slug					string					`json:"slug" validate:"required" example:"listing_slug"`
}

type CreateBidSchema struct {
	Amount					float64			`json:"amount" validate:"required" example:"1000.00"`
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

type CategoriesResponseSchema struct {
	ResponseSchema
	Data					[]models.Category	`json:"data"`
}

type BidResponseDataSchema struct {
	Listing					string				`json:"listing"`
	Bids					[]models.Bid		`json:"bids"`
}

type BidsResponseSchema struct {
	ResponseSchema
	Data					BidResponseDataSchema		`json:"data"`
}

type BidResponseSchema struct {
	ResponseSchema
	Data					models.Bid			`json:"data"`			
}