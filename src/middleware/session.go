package middleware

import (
	"log"
	"net/http"

	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

var LoginInfo interface{}

func SessionCheck() gin.HandlerFunc {
	return func(c *gin.Context) {

		session := sessions.Default(c)
		LoginInfo = session.Get("UserJWT")

		if LoginInfo == nil {
			log.Println("ログインしていません")
			c.Redirect(http.StatusMovedPermanently, "/app/middle_name/login")
			c.Abort()
		} else {
			c.Set("UserJWT", LoginInfo)
			c.Next()
		}
		log.Println("ログインチェック終わり")
	}
}
