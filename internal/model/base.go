package model

import (
	"time"

	uuid "github.com/google/uuid"
)

type BaseModel struct {
	Id        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt time.Time
}
