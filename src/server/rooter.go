package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"example.com/m/v2/database"
	"example.com/m/v2/function/encryption"
	"example.com/m/v2/middleware"
	"example.com/m/v2/model"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var LoginInfo interface{}

func NewRouter() *gin.Engine {
	// リリースモード
	// gin.SetMode(gin.ReleaseMode)

	routerGlobal := gin.Default()

	routerGlobal.LoadHTMLGlob("views/*.html")

	router := routerGlobal.Group("app/middle_name")
	{
		router.Static("/store", "./store")
		router.Static("/assets", "./assets")

		store := cookie.NewStore([]byte("secret"))
		router.Use(sessions.Sessions("mysession", store))

		login := router.Group("/user")
		login.Use(sessionCheck())
		{
			//一覧
			login.GET("/index", func(c *gin.Context) {
				session := sessions.Default(c)
				userId := fmt.Sprintf("%v", session.Get("Uuid"))
				getUser := database.GetUserFromUuid(userId)

				middleNames := database.DbGetCreatedMiddleNames(userId)

				c.HTML(http.StatusOK, "index.html", gin.H{
					"name":   getUser.Username,
					"middle": middleNames,
				})
			})

			// ユーザー登録画面
			login.GET("/create", func(c *gin.Context) {
				c.HTML(200, "create.html", gin.H{})
			})

			// ミッドルネーム作成
			login.POST("/create", func(c *gin.Context) {
				session := sessions.Default(c)
				userId := fmt.Sprintf("%v", session.Get("Uuid"))

				var form model.CreatedMiddleNames
				// ここがバリデーション部分
				if err := c.Bind(&form); err != nil {
					middleNames := database.DbGetCreatedMiddleNames(userId)
					c.HTML(http.StatusBadRequest, "create.html", gin.H{"middleNames": middleNames, "err": err})
					c.Abort()
				} else {
					mr := database.DBGetRandomMrData().Mr
					lName := c.PostForm("lname")
					surName := database.DBGetRandomMrData().Mr
					commonName := database.DBGetRandomMrData().Mr
					fName := c.PostForm("fname")

					database.DbMiddleNameInsert(mr, lName, surName, commonName, fName, userId)
					c.Redirect(302, "/app/middle_name/user/index")
				}
			})

		}

		router.GET("/", func(c *gin.Context) {
			c.HTML(200, "home.html", gin.H{})
		})

		router.GET("/login", func(c *gin.Context) {
			c.HTML(200, "login.html", gin.H{})
		})

		router.POST("/login", func(c *gin.Context) {

			formPassword := c.PostForm("password")
			formName := c.PostForm("username")
			dbPassword := database.GetUser(formName).Password
			dbUserUuid := database.GetUser(formName).UserUUID

			if err := middleware.CompareHashAndPassword(dbPassword, formPassword); err != nil {
				log.Println("login false")
				c.HTML(http.StatusBadRequest, "login.html", gin.H{"err": "ログインできませんでした。"})
				c.Abort()
			} else {
				log.Println("login success")

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
				tokenString, _ := token.SignedString([]byte(os.Getenv("SIGNINGKEY")))

				session := sessions.Default(c)
				session.Set("UserJWT", tokenString)
				session.Set("Uuid", dbUserUuid)
				session.Save()

				text := encryption.Compress(tokenString)

				database.DbSessionUpdate(dbPassword, text)

				c.Redirect(302, "/app/middle_name/user/index")
			}
		})

		// ユーザー登録画面
		router.GET("/signup", func(c *gin.Context) {
			c.HTML(200, "signup.html", gin.H{})
		})

		// ユーザー登録
		router.POST("/signup", func(c *gin.Context) {
			var form model.User
			if err := c.Bind(&form); err != nil {
				log.Print(err)
				c.HTML(http.StatusBadRequest, "signup.html", gin.H{"err": err})
				c.Abort()
			} else {
				username := c.PostForm("username")
				password := c.PostForm("password")
				email := c.PostForm("email")
				userid := uuid.New().String()
				session := "NoLogin"
				formStruct := model.User{
					Username: username,
					Password: password,
					Email:    email,
					Session:  session,
				}
				if ok, err := formStruct.Validate(); !ok {
					log.Print(err)
					c.HTML(http.StatusBadRequest, "signup.html", gin.H{"err": err})
					c.Abort()
				}
				if err := database.CreateUser(userid, username, password, email, session); len(err) != 0 {
					log.Print("同じユーザーが存在します")
					log.Print(len(err))
					log.Print(err)
					c.HTML(http.StatusBadRequest, "signup.html", gin.H{"err": "同じユーザーが存在します"})
					c.Abort()
				}
				c.Redirect(302, "/login")
			}
		})

	}

	routerGlobal.Run(":8001")
	// routerGlobal.Run()

	return routerGlobal
}

func sessionCheck() gin.HandlerFunc {
	return func(c *gin.Context) {

		session := sessions.Default(c)
		LoginInfo = session.Get("UserJWT")

		if LoginInfo == nil {
			log.Println("ログインしていません")
			c.Redirect(http.StatusMovedPermanently, "app/middle_name/login")
			c.Abort()
		} else {
			c.Set("UserJWT", LoginInfo)
			c.Next()
		}
		log.Println("ログインチェック終わり")
	}
}
