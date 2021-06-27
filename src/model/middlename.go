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

type CreateMiddleName struct {
	Id         int
	Mr         string
	LName      string
	SurName    string
	CommonName string
	FName      string
}

type Mr struct {
	gorm.Model
	Mr string
}

type SN struct {
	gorm.Model
	SurName string
}

type CN struct {
	gorm.Model
	CommonName string
}
