package models

import (
	"time"

	"github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID 			uuid.UUID `json:"-" gorm:"type:uuid;primary_key"`
	CreatedAt 	time.Time `json:"-"`
	UpdatedAt 	time.Time `json:"-"`
}

func (obj *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
    obj.ID = uuid.NewV4()
    return
}