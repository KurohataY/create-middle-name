package model

import (
	"github.com/jinzhu/gorm"
)

type CreatedMiddleNames struct {
	gorm.Model
	Mr         string
	LName      string
	SurName    string
	CommonName string
	FName      string
	UserId     string
}
