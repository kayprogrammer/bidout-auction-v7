package models

import (
	"time"

	"github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID 			uuid.UUID `json:"-" gorm:"type:uuid;primary_key;not null;default:uuid_generate_v4()"`
	CreatedAt 	time.Time `json:"-" gorm:"not null"`
	UpdatedAt 	time.Time `json:"-" gorm:"not null"`
}

func (obj *BaseModel) BeforeUpdate(tx *gorm.DB) (err error) {
    obj.UpdatedAt = time.Now().UTC()
    return
}

type File struct {
	BaseModel
	ResourceType		string 		`json:"resource_type" gorm:"not null"`
}

type ShortUserData struct {
	Name				string				`json:"name" example:"John Doe"`
	Avatar				*string				`json:"avatar" example:"https://my-avatar.com"`
}