package models

import (
	"time"
	"github.com/satori/go.uuid"

)

type BaseModel struct {
	ID uuid.UUID `json:"id" gorm:"type:uuid;primary_key;"`
	CreatedAt time.Time
	UpdatedAt time.Time
}