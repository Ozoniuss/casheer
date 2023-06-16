package model

import (
	"time"

	"gorm.io/gorm"
)

type BaseModel struct {
	Id        int
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt gorm.DeletedAt
}
