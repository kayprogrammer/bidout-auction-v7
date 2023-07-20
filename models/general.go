package models

type SiteDetail struct {
	BaseModel
	Name 		string 		`json:"name" gorm:"default:Kay's Auction House"`
	Email 		string 		`json:"email" gorm:"default:kayprogrammer1@gmail.com"`
	Phone 		string		`json:"phone" gorm:"default:+2348133831036"`
	Address 	string		`json:"address" gorm:"default:234, Lagos, Nigeria"`
	Fb 			string		`json:"fb" gorm:"default:https://facebook.com"`
	Tw 			string		`json:"tw" gorm:"default:https://twitter.com"`
	Wh 			string		`json:"wh" gorm:"default:https://wa.me/2348133831036"`
	Ig 			string		`json:"ig" gorm:"default:https://instagram.com"`
}