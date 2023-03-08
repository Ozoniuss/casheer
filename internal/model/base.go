package model

import (
	"time"

	uuid "github.com/google/uuid"
)

type BaseModel struct {
	Id        uuid.UUID `gorm:"default:gen_random_uuid()"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
