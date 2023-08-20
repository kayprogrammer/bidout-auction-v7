package models

import (
	"fmt"
	"time"

	"github.com/kayprogrammer/bidout-auction-v7/config"
	"github.com/kayprogrammer/bidout-auction-v7/utils"
	uuid "github.com/satori/go.uuid"
	"gorm.io/gorm"

)

type User struct {
	BaseModel
	FirstName				string			`json:"first_name" gorm:"type: varchar(50);not null" validate:"required,max=50" example:"John"`
	LastName				string			`json:"last_name" gorm:"type: varchar(50);not null" validate:"required,max=50" example:"Doe"`
	Email					string			`json:"email" gorm:"not null;unique;" validate:"required,min=5,email" example:"johndoe@email.com"`
	Password				string			`json:"password" gorm:"not null" validate:"required,min=8,max=50" example:"strongpassword"`
	IsEmailVerified			*bool			`json:"is_email_verified" gorm:"default:false" swaggerignore:"true"`
	IsSuperuser				*bool			`json:"is_superuser" gorm:"default:false" swaggerignore:"true"`
	IsStaff					*bool			`json:"is_staff" gorm:"default:false" swaggerignore:"true"`
	TermsAgreement			bool			`json:"terms_agreement" gorm:"default:false" validate:"eq=true"`
	AvatarId				*uuid.UUID		`json:"avatar_id" gorm:"null" swagger:"ignore" swaggerignore:"true"`
	Avatar					*File			`gorm:"foreignKey:AvatarId;constraint:OnDelete:SET NULL;null;" swaggerignore:"true"`
}

func (user User) FullName() string {
	return fmt.Sprintf("%s %s", user.FirstName, user.LastName)
}

func (user *User) BeforeCreate(tx *gorm.DB) (err error) {
	user.Password = utils.HashPassword(user.Password)
	return
}

func (user User) GetAvatarUrl(db *gorm.DB) *string {
	avatarId := user.AvatarId
	if avatarId != nil {
		avatar := File{}
		db.Take(&avatar, avatarId)
		url := utils.GenerateFileUrl(avatarId.String(), "avatars", avatar.ResourceType)
		return &url
	}
	return nil
}

func (user User) GetAvatarUploadUrl(db *gorm.DB) *utils.SignatureFormat {
	avatarId := user.AvatarId
	if avatarId != nil {
		uploadData := utils.GenerateFileSignature(avatarId.String(), "avatars")
		return &uploadData
	}
	return nil
}

type GuestUser struct {
	BaseModel
}

type Jwt struct {
	BaseModel
	UserId				uuid.UUID		`json:"user_id" gorm:"not null;unique;"`
	User				User			`gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE;unique;not null"`
	Access				string			`json:"access" gorm:"not null"`
	Refresh				string			`json:"refresh" gorm:"not null"`
}

type Otp struct {
	BaseModel
	UserId				uuid.UUID		`json:"user_id" gorm:"not null;unique;"`
	User				User			`gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE;unique;not null"`
	Code				*int			`json:"code" gorm:"not null"`
}

func (otp *Otp) BeforeSave(tx *gorm.DB) (err error) {
	code := utils.GetRandomInt(6)
	otp.Code = &code
	return
}

func (obj Otp) CheckExpiration() bool {
	currentTime := time.Now().UTC()
	diff := int64(currentTime.Sub(obj.UpdatedAt).Seconds())
	emailExpirySecondsTimeout := config.GetConfig().EmailOTPExpireSeconds
	return diff > emailExpirySecondsTimeout
}