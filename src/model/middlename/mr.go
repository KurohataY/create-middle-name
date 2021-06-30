package model

import (
	"github.com/jinzhu/gorm"
)

type Mr struct {
	gorm.Model
	Mr string
}
