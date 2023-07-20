package models

import "gorm.io/gorm"

type SiteDetail struct {
	BaseModel
	Name 		string 		`json:"name"`
	Email 		string 		`json:"email" gorm:"default:kayprogrammer1@gmail.com"`
	Phone 		string		`json:"phone" gorm:"default:+2348133831036"`
	Address 	string		`json:"address" gorm:"default:234, Lagos, Nigeria"`
	Fb 			string		`json:"fb" gorm:"default:https://facebook.com"`
	Tw 			string		`json:"tw" gorm:"default:https://twitter.com"`
	Wh 			string		`json:"wh" gorm:"default:https://wa.me/2348133831036"`
	Ig 			string		`json:"ig" gorm:"default:https://instagram.com"`
}

type Subscriber struct {
	BaseModel
	Email		string		`json:"email"`
	Exported	string		`json:"-"`
}

func (obj *SiteDetail) BeforeCreate(tx *gorm.DB) (err error) {
    obj.Name = "Kay's Auction House"
    return
}
