package model

import (
	"github.com/jinzhu/gorm"
)

type MiddleName struct {
	gorm.Model
	Id         int
	Mr         string
	LName      string
	SurName    string
	CommonName string
	FName      string
}
