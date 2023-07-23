package schemas

import "github.com/kayprogrammer/bidout-auction-v7/models"

type SiteDetailResponseSchema struct {
	ResponseSchema
	Data			models.SiteDetail		`json:"data"`
}

type SubscriberResponseSchema struct {
	ResponseSchema
	Data			models.Subscriber		`json:"data"`
}

type ReviewsResponseSchema struct {
	ResponseSchema
	Data			[]models.Review		`json:"data"`
}