package models

import (
	"gorm.io/gorm"
)

type DenylistedIP struct {
	gorm.Model
	HTTPPath string
	IP       string
}
