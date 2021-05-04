package model

import (
  "github.com/jinzhu/gorm"
)

// Tweetモデル宣言
// モデルはDBのテーブル構造をGOの構造体で表したもの
type Tweet struct {
    gorm.Model
    Content string `form:"content" binding:"required"`
}
