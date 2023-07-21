package models

import (
	"gorm.io/gorm"
	"github.com/satori/go.uuid"

)

type SiteDetail struct {
	BaseModel
	Name 		string 		`json:"name" gorm:"type:varchar(50);not null"`
	Email 		string 		`json:"email" gorm:"default:kayprogrammer1@gmail.com;not null"`
	Phone 		string		`json:"phone" gorm:"default:+2348133831036;type:varchar(20);not null"`
	Address 	string		`json:"address" gorm:"default:234, Lagos, Nigeria;not null"`
	Fb 			string		`json:"fb" gorm:"default:https://facebook.com;not null"`
	Tw 			string		`json:"tw" gorm:"default:https://twitter.com;not null"`
	Wh 			string		`json:"wh" gorm:"default:https://wa.me/2348133831036;not null"`
	Ig 			string		`json:"ig" gorm:"default:https://instagram.com;not null"`
}

type Subscriber struct {
	BaseModel
	Email		string		`json:"email" gorm:"not null" validate:"required,min=5"`
	Exported	string		`json:"-" gorm:"default:false"`
}

func (obj *SiteDetail) BeforeCreate(tx *gorm.DB) (err error) {
    obj.Name = "Kay's Auction House"
	obj.BaseModel.BeforeCreate(tx)
    return
}

type Review struct {
	BaseModel
    ReviewerId			uuid.UUID		`json:"reviewer_id" gorm:"not null"`
	Reviewer			User			`gorm:"foreignKey:ReviewerId;constraint:OnDelete:CASCADE;unique;not null"`
	Show				bool			`json:"-" gorm:"not null"`
	Text				string			`json:"text" gorm:"type:varchar(200);not null"`
}