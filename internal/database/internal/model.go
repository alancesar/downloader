package internal

import (
	"gorm.io/gorm"
)

type (
	Download struct {
		gorm.Model
		ETag string
		URL  string
	}
)
