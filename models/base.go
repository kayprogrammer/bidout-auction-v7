package models

import (
	"time"

	"github.com/satori/go.uuid"
)

type BaseModel struct {
	ID 			uuid.UUID `json:"-" gorm:"type:uuid;primarykey;not null;default:uuid_generate_v4()"`
	CreatedAt 	time.Time `json:"-" gorm:"not null"`
	UpdatedAt 	time.Time `json:"-" gorm:"not null"`
}

type File struct {
	BaseModel
	ResourceType		string 		`json:"resource_type" gorm:"not null"`
}

type ShortUserData struct {
	ID					uuid.UUID			`json:"id" example:"2c64c881-59ca-4916-b2bc-8cfb75c3f09b"`
	Name				string				`json:"name" example:"John Doe"`
	Avatar				*string				`json:"avatar" example:"https://my-avatar.com"`
}