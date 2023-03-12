package entries

import "gorm.io/gorm"

type handler struct {
	db *gorm.DB
}

func NewHandler(db *gorm.DB) handler {
	return handler{
		db: db,
	}
}
