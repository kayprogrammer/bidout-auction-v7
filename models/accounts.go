package models

import (
	"os"
	"time"
	"fmt"
	"strconv"

	"github.com/satori/go.uuid"
)

type User struct {
	BaseModel
	FirstName				string			`json:"first_name" gorm:"type: varchar(50);not null"`
	LastName				string			`json:"last_name" gorm:"type: varchar(50);not null"`
	Email					string			`json:"email" gorm:"not null"`
	Password				string			`json:"password" gorm:"not null"`
	IsEmailVerified			bool			`json:"is_email_verified" gorm:"not null"`
	IsSuperuser				bool			`json:"is_superuser" gorm:"not null"`
	IsStaff					bool			`json:"is_staff" gorm:"not null"`
	TermsAgreement			bool			`json:"terms_agreement" gorm:"not null"`
	AvatarId				uuid.UUID		`json:"avatar_id" gorm:"not null"`
	Avatar					File			`gorm:"foreignKey:AvatarId;constraint:OnDelete:SET NULL;"`
}

func (obj User) FullName() string {
	fullName := "%s %s"
	return fmt.Sprintf(fullName, obj.FirstName, obj.LastName)
}

type Jwt struct {
	BaseModel
	UserId				uuid.UUID		`json:"user_id" gorm:"not null"`
	User				User			`gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE;unique;not null"`
	Access				string			`json:"access" gorm:"not null"`
	Refresh				string			`json:"refresh" gorm:"not null"`
}

type Otp struct {
	BaseModel
	UserId				uuid.UUID		`json:"user_id" gorm:"not null"`
	User				User			`gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE;unique;not null"`
	Code				string			`json:"code" gorm:"type:varchar(6);not null"`
}

func (obj Otp) CheckExpiration() bool {
	currentTime := time.Now().UTC()
	diff := int64(obj.UpdatedAt.Sub(currentTime).Seconds())
	emailExpirySecondsTimeout, err := strconv.ParseInt(os.Getenv("EMAIL_OTP_EXPIRE_SECONDS"), 10, 64)
    if err != nil {
        fmt.Println("Error parsing comparison value:", err)
    }

	if diff > emailExpirySecondsTimeout {
		return true
	}
	return false
}