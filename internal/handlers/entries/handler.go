package entries

import (
	"fmt"

	"github.com/Ozoniuss/casheer/internal/config"
	"gorm.io/gorm"
)

type handler struct {
	db      *gorm.DB
	apiPath string
}

func NewHandler(db *gorm.DB, config config.Server) handler {
	return handler{
		db:      db,
		apiPath: fmt.Sprintf("http://%s:%d/api/entries", config.Address, config.Port),
	}
}
