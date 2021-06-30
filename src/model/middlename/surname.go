package model

import (
	"github.com/jinzhu/gorm"
)

type SN struct {
	gorm.Model
	SurName string
}
