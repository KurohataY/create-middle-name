package model

import (
	"log"

	"github.com/jinzhu/gorm"
	"gopkg.in/go-playground/validator.v9"
)

// User モデルの宣言
type User struct {
	gorm.Model
	UserUUID string
	Username string `form:"username" binding:"required" gorm:"unique;not null" validate:"min=8"`
	Password string `form:"password" binding:"required" validate:"min=8,max=50,required"`
	Email    string `form:"email" binding:"required" validate:"required,email"`
	Session  string
}

// バリデーションチェックを行って、結果とNGの場合にエラーメッセージ群を返す
func (form *User) Validate() (ok bool, result map[string]string) {
	result = make(map[string]string)
	// 構造体のデータをタグで定義した検証方法でチェック
	// err := validator.New().Struct(*form)
	validate := validator.New()
	err := validate.Struct(*form)
	if err != nil {
		errors := err.(validator.ValidationErrors)
		if len(errors) != 0 {
			for i := range errors {
				// フィールドごとに、検証
				switch errors[i].StructField() {
				case "Password":
					switch errors[i].Tag() {
					case "required", "min", "max":
						log.Println("8文字以上50文字以下で入力してください")
						result["Password"] = "8文字以上50文字以下で入力してください"
					case "Email":
						switch errors[i].Tag() {
						case "email":
							log.Println("メールアドレスの形式で入力してください")
							result["email"] = "メールアドレスの形式で入力してください"
						}
					case "Username":
						switch errors[i].Tag() {
						case "min":
							log.Println("８文字以上で入力してください")
							result["name"] = "８文字以上で入力してください"
						}
					}
				}
			}
		}
		return false, result
	}
	return true, result
}
