package model

import (
	"github.com/jinzhu/gorm"
)

type CN struct {
	gorm.Model
	CommonName string
}
