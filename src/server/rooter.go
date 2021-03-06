package server

import (
	"fmt"
	"log"
	"net/http"

	"example.com/m/v2/database"
	"example.com/m/v2/function/encryption"
	"example.com/m/v2/middleware"
	"example.com/m/v2/model"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

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

				tokenString, err := middleware.GetTokenHandler(dbUserUuid, formName)
				if err != nil {
					log.Print(err)
					c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{
						"Error": err.Error(),
					})
				}

				session := sessions.Default(c)
				session.Set("UserJWT", tokenString)
				session.Set("Uuid", dbUserUuid)
				// session.Options(sessions.Options{Path: "/", MaxAge: -1})
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
				c.Redirect(302, "/app/middle_name/login")
			}
		})

		login := router.Group("/user")
		login.Use(middleware.LoginCheck())
		{
			//ログイン後
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

			//一覧
			login.GET("/history", func(c *gin.Context) {
				session := sessions.Default(c)
				userId := fmt.Sprintf("%v", session.Get("Uuid"))
				getUser := database.GetUserFromUuid(userId)

				middleNames := database.DbGetCreatedMiddleNames(userId)

				c.HTML(http.StatusOK, "history.html", gin.H{
					"name":   getUser.Username,
					"middle": middleNames,
				})
			})

			login.GET("/create", func(c *gin.Context) {
				c.HTML(200, "create.html", gin.H{})
			})

			// ミッドルネーム作成
			login.POST("/create", func(c *gin.Context) {
				session := sessions.Default(c)
				userId := fmt.Sprintf("%v", session.Get("Uuid"))

				var form model.CreatedMiddleNames

				if err := c.Bind(&form); err != nil {
					middleNames := database.DbGetCreatedMiddleNames(userId)
					c.HTML(http.StatusBadRequest, "create.html", gin.H{"middleNames": middleNames, "err": err})
					c.Abort()
				} else {
					mr := database.DBGetRandomMrData().Mr
					lName := c.PostForm("lname")
					surName := database.DBGetRandomSNData().SurName
					commonName := database.DBGetRandomCNData().CommonName
					fName := c.PostForm("fname")

					database.DbMiddleNameInsert(mr, lName, surName, commonName, fName, userId)
					c.Redirect(302, "/app/middle_name/user/result")
				}
			})

			login.GET("/result", func(c *gin.Context) {
				session := sessions.Default(c)
				userId := fmt.Sprintf("%v", session.Get("Uuid"))
				middleName := database.DbMiddleNameLastFind(userId)
				c.HTML(200, "result.html", gin.H{
					"middleName": middleName,
				})
			})

			login.GET("/logout", func(c *gin.Context) {
				c.HTML(200, "logout.html", gin.H{})
			})

			login.POST("/logout", func(c *gin.Context) {
				err := middleware.Logout(c)
				print(err)
				if err != nil {
					c.HTML(http.StatusBadRequest, "user/logout.html", gin.H{"err": "ログアウトできませんでした。もう一度お試しください！"})
					c.Abort()
				}
				c.Redirect(302, "/app/middle_name/")
			})

		}

	}

	routerGlobal.Run(":8001")

	return routerGlobal
}
