package routes

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kayprogrammer/bidout-auction-v7/database"
	"github.com/kayprogrammer/bidout-auction-v7/models"
)
type GeneralSerializer struct {
	Name        string 		`json:"name"`
	Email 		string 		`json:"email"`
	Phone 		string		`json:"phone"`
	Address 	string		`json:"address"`
	Fb 			string		`json:"fb"`
	Tw 			string		`json:"tw"`
	Wh 			string		`json:"wh"`
	Ig 			string		`json:"ig"`
}

func CreateResponseSiteDetail(sitedetailModel models.SiteDetail) GeneralSerializer {
	return GeneralSerializer{Name: sitedetailModel.Name}
}

func GetSiteDetails(c *fiber.Ctx) error {
	var sitedetail models.SiteDetail

	database.Database.Db.FirstOrCreate(&sitedetail, &sitedetail)
	responseSiteDetail := CreateResponseSiteDetail(sitedetail)

	return c.Status(200).JSON(responseSiteDetail)
}