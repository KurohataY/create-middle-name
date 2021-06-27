package server

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"example.com/m/v2/function/encryption"
	"example.com/m/v2/function/funcDB"
	"example.com/m/v2/model"
	plugins "example.com/m/v2/plugins/crypto"
	"github.com/form3tech-oss/jwt-go"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

var LoginInfo interface{}

func NewRouter() *gin.Engine {
	router := gin.Default()
	router.LoadHTMLGlob("views/*.html")
	// router.LoadHTMLGlob("views/user/*.html")
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
			// LoginInfo = session.Get("Uuid")
			getUser := funcDB.GetUserFromUuid(userId)

			fmt.Fprintf(os.Stdout, "%v", userId)
			fmt.Fprint(os.Stdout, "ペルソナ！")
			c.HTML(http.StatusOK, "index.html", gin.H{
				// htmlに渡す変数を定義
				"uuid":   userId,
				"name":   getUser.Username,
				"email":  getUser.Email,
				"middle": []model.CreateMiddleName{{1, "a", "c_mika", "b", "d", "d"}},
				// "name":
			})
		})
	}

	//トップ画面
	router.GET("/", func(c *gin.Context) {
		tweets := funcDB.DbGetAll()
		c.HTML(200, "home.html", gin.H{"tweets": tweets})
	})

	//登録
	// router.POST("/new", func(c *gin.Context) {
	// 	var form model.Tweet
	// 	// ここがバリデーション部分
	// 	if err := c.Bind(&form); err != nil {
	// 		tweets := funcDB.DbGetAll()
	// 		c.HTML(http.StatusBadRequest, "index.html", gin.H{"tweets": tweets, "err": err})
	// 		c.Abort()
	// 	} else {
	// 		content := c.PostForm("content")
	// 		funcDB.DbInsert(content)
	// 		c.Redirect(302, "/")
	// 	}
	// })

	//投稿詳細
	router.GET("/detail/:id", func(c *gin.Context) {
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic(err)
		}
		tweet := funcDB.DbGetOne(id)
		c.HTML(200, "detail.html", gin.H{"tweet": tweet})
	})

	//更新
	// router.POST("/update/:id", func(c *gin.Context) {
	// 	n := c.Param("id")
	// 	id, err := strconv.Atoi(n)
	// 	if err != nil {
	// 		panic("ERROR")
	// 	}
	// 	tweet := c.PostForm("tweet")
	// 	funcDB.DbUpdate(id, tweet)
	// 	c.Redirect(302, "/")
	// })

	//削除確認
	router.GET("/delete_check/:id", func(c *gin.Context) {
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("ERROR")
		}
		tweet := funcDB.DbGetOne(id)
		c.HTML(200, "delete.html", gin.H{"tweet": tweet})
	})

	//削除
	router.POST("/delete/:id", func(c *gin.Context) {
		n := c.Param("id")
		id, err := strconv.Atoi(n)
		if err != nil {
			panic("ERROR")
		}
		funcDB.DbDelete(id)
		c.Redirect(302, "/")

	})

	// ユーザーログイン画面
	router.GET("/login", func(c *gin.Context) {
		c.HTML(200, "login.html", gin.H{})
	})

	// ユーザーログイン
	router.POST("/login", func(c *gin.Context) {

		// フォームから取得したユーザーパスワード
		formPassword := c.PostForm("password")
		// DBから取得したユーザーパスワード(Hash)
		formName := c.PostForm("username")
		dbPassword := funcDB.GetUser(formName).Password
		dbUserUuid := funcDB.GetUser(formName).UserUUID
		// ユーザーパスワードの比較
		if err := plugins.CompareHashAndPassword(dbPassword, formPassword); err != nil {
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
			fmt.Print(dbUserUuid)
			session.Save()

			text := encryption.Compress(tokenString)

			funcDB.DbSessionUpdate(dbPassword, text)

			c.Redirect(302, "/user/index")
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
			if err := funcDB.CreateUser(userid, username, password, email, session); len(err) != 0 {
				log.Print("同じユーザーが存在します")
				log.Print(len(err))
				log.Print(err)
				c.HTML(http.StatusBadRequest, "signup.html", gin.H{"err": "同じユーザーが存在します"})
				c.Abort()
			}
			c.Redirect(302, "/login")
		}
	})
	router.Run()

	return router
}

func sessionCheck() gin.HandlerFunc {
	return func(c *gin.Context) {

		session := sessions.Default(c)
		LoginInfo = session.Get("UserJWT")

		// セッションがない場合、ログインフォームをだす
		if LoginInfo == nil {
			log.Println("ログインしていません")
			c.Redirect(http.StatusMovedPermanently, "/login")
			c.Abort() // これがないと続けて処理されてしまう
		} else {
			c.Set("UserJWT", LoginInfo) // ユーザidをセット
			c.Next()
		}
		log.Println("ログインチェック終わり")
	}
}
