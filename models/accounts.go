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
	FirstName				string			`json:"first_name"`
	LastName				string			`json:"last_name"`
	Email					string			`json:"email"`
	Password				string			`json:"password"`
	IsEmailVerified			bool			`json:"is_email_verified"`
	IsSuperuser				bool			`json:"is_superuser"`
	IsStaff					bool			`json:"is_staff"`
	TermsAgreement			bool			`json:"terms_agreement"`
	AvatarId				uuid.UUID		`json:"avatar_id"`
	Avatar					File			`gorm:"foreignKey:AvatarId;constraint:OnDelete:SET NULL;"`
}

func (obj User) FullName() string {
	fullName := "%s %s"
	return fmt.Sprintf(fullName, obj.FirstName, obj.LastName)
}

type Jwt struct {
	BaseModel
	UserId				uuid.UUID		`json:"user_id"`
	User				User			`gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE;unique;"`
	Access				string			`json:"access"`
	Refresh				string			`json:"refresh"`
}

type Otp struct {
	BaseModel
	UserId				uuid.UUID		`json:"user_id"`
	User				User			`gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE;unique;"`
	Code				string			`json:"code"`
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