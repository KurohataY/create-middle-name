package model

import (
  "github.com/jinzhu/gorm"
)

// User モデルの宣言
type User struct {
    gorm.Model
    Username string `form:"username" binding:"required" gorm:"unique;not null"`
    Password string `form:"password" binding:"required"`
}
