package models

import (
	"time"
	"fmt"

	"github.com/satori/go.uuid"
	"gorm.io/gorm"
)

type BaseModel struct {
	ID 			uuid.UUID `json:"-" gorm:"type:uuid;primary_key"`
	CreatedAt 	time.Time `json:"-"`
	UpdatedAt 	time.Time `json:"-"`
}

func (obj *BaseModel) BeforeCreate(tx *gorm.DB) (err error) {
	fmt.Println("Working")
    obj.ID = uuid.NewV4()
    obj.CreatedAt = time.Now().UTC()
    obj.UpdatedAt = time.Now().UTC()
    return
}

func (obj *BaseModel) BeforeUpdate(tx *gorm.DB) (err error) {
    obj.UpdatedAt = time.Now().UTC()
    return
}

type File struct {
	BaseModel
	ResourceType		string
}