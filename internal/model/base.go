package model

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	Id        int `gorm:"primaryKey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
