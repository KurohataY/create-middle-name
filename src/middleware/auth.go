package middleware

import (
	"log"
	"net/http"
	"os"
	"time"

	jwt "github.com/form3tech-oss/jwt-go"
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
)

// GetTokenHandler get token
func GetTokenHandler(dbUserUuid string, formName string) (string, error) {

	// headerのセット
	// headerのセット
	token := jwt.New(jwt.SigningMethodHS256)

	// claimsのセット
	claims := token.Claims.(jwt.MapClaims)
	claims["admin"] = true
	claims["sub"] = dbUserUuid
	claims["name"] = formName
	claims["iat"] = time.Now()
	claims["exp"] = time.Now().Add(time.Hour * 24).Unix()

	// 電子署名
	return token.SignedString([]byte(os.Getenv("SIGNINGKEY")))

}

var LoginInfo interface{}

func LoginCheck() gin.HandlerFunc {
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

func Logout(c *gin.Context) error {
	session := sessions.Default(c)
	// session.Delete("UserJWT")
	// session.Delete("Uuid")
	session.Clear()
	session.Options(sessions.Options{Path: "/", MaxAge: -1})
	err := session.Save()
	return err
}
