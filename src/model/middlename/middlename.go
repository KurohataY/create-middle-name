package model

import (
	"github.com/jinzhu/gorm"
)

type CreateMiddleNameInit struct {
	gorm.Model
	Mr         string
	LName      string
	SurName    string
	CommonName string
	FName      string
}

type CreatedMiddleName struct {
	Id         int
	Mr         string
	LName      string
	SurName    string
	CommonName string
	FName      string
}
