package sqlite

import (
	"embed"
)

//go:embed sqlite/*.sql
var Migrations embed.FS
