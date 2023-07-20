package schemas

import "github.com/kayprogrammer/bidout-auction-v7/models"

type SiteDetailDataSchema struct {
	Name        string 		`json:"name"`
	Email 		string 		`json:"email"`
	Phone 		string		`json:"phone"`
	Address 	string		`json:"address"`
	Fb 			string		`json:"fb"`
	Tw 			string		`json:"tw"`
	Wh 			string		`json:"wh"`
	Ig 			string		`json:"ig"`
}

type SiteDetailResponseSchema struct {
	ResponseSchema
	Data			models.SiteDetail		`json:"data"`
}