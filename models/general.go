package models

import (
	"github.com/kayprogrammer/bidout-auction-v7/utils"
	"github.com/satori/go.uuid"
	"gorm.io/gorm"
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
	Email		string		`json:"email" gorm:"not null" validate:"required,min=5,email"`
	Exported	bool		`json:"-" gorm:"default:false"`
}

func (obj *SiteDetail) BeforeCreate(tx *gorm.DB) (err error) {
    obj.Name = "Kay's Auction House"
    return
}


type Review struct {
	BaseModel
    ReviewerId			uuid.UUID			`json:"-" gorm:"not null"`
	ReviewerObj			User				`json:"-" gorm:"foreignKey:ReviewerId;constraint:OnDelete:CASCADE;not null"`
	Reviewer			ShortUserData		`json:"reviewer" gorm:"-"`	
	Show				bool				`json:"-" gorm:"default:false"`
	Text				string				`json:"text" gorm:"type:varchar(200);not null"`
}

func (obj Review) Init(db *gorm.DB) Review{
	reviewer := User{}
	db.Find(&reviewer,"id = ?", obj.ReviewerId)
	name := reviewer.FullName()
	obj.Reviewer.Name = name

	avatarId := reviewer.AvatarId
	if avatarId != nil {
		avatar := File{}
		db.Find(&avatar,"id = ?", avatarId)
		url := utils.GenerateFileUrl(avatarId.String(), "avatars", avatar.ResourceType)
		obj.Reviewer.Avatar = &url
	}
	return obj
}