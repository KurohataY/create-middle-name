package main

import (
    "log"
    "net/http"
    "strconv"
    "github.com/gin-gonic/gin"
    _ "github.com/go-sql-driver/mysql" //直接的な記述が無いが、インポートしたいものに対しては"_"を頭につける決まり
    _ "github.com/joho/godotenv"
    "example.com/m/v2/model"
    "example.com/m/v2/methods"
    "example.com/m/v2/plugins/crypto"
)

func main() {
    router := gin.Default()
    router.LoadHTMLGlob("views/*.html")
    router.Static("/assets", "./assets")
    router.Static("/store", "./store")

    methods.DbInit()

    //一覧
    router.GET("/", func(c *gin.Context) {
        tweets := methods.DbGetAll()
        c.HTML(200, "index.html", gin.H{"tweets": tweets})
    })

    //一覧
    router.GET("/home", func(c *gin.Context) {
        tweets := methods.DbGetAll()
        c.HTML(200, "home.html", gin.H{"tweets": tweets})
    })

    //登録
    router.POST("/new", func(c *gin.Context) {
        var form model.Tweet
        // ここがバリデーション部分
        if err := c.Bind(&form); err != nil {
            tweets := methods.DbGetAll()
            c.HTML(http.StatusBadRequest, "index.html", gin.H{"tweets": tweets, "err": err})
            c.Abort()
        } else {
            content := c.PostForm("content")
            methods.DbInsert(content)
            c.Redirect(302, "/")
        }
    })

    //投稿詳細
    router.GET("/detail/:id", func(c *gin.Context) {
        n := c.Param("id")
        id, err := strconv.Atoi(n)
        if err != nil {
            panic(err)
        }
        tweet := methods.DbGetOne(id)
        c.HTML(200, "detail.html", gin.H{"tweet": tweet})
    })

    //更新
    router.POST("/update/:id", func(c *gin.Context) {
        n := c.Param("id")
        id, err := strconv.Atoi(n)
        if err != nil {
            panic("ERROR")
        }
        tweet := c.PostForm("tweet")
        methods.DbUpdate(id, tweet)
        c.Redirect(302, "/")
    })

    //削除確認
    router.GET("/delete_check/:id", func(c *gin.Context) {
        n := c.Param("id")
        id, err := strconv.Atoi(n)
        if err != nil {
            panic("ERROR")
        }
        tweet := methods.DbGetOne(id)
        c.HTML(200, "delete.html", gin.H{"tweet": tweet})
    })

    //削除
    router.POST("/delete/:id", func(c *gin.Context) {
        n := c.Param("id")
        id, err := strconv.Atoi(n)
        if err != nil {
            panic("ERROR")
        }
        methods.DbDelete(id)
        c.Redirect(302, "/")

    })

    // ユーザーログイン画面
    router.GET("/login", func(c *gin.Context) {

        c.HTML(200, "login.html", gin.H{})
    })

    // ユーザーログイン
    router.POST("/login", func(c *gin.Context) {

        // DBから取得したユーザーパスワード(Hash)
        dbPassword := methods.GetUser(c.PostForm("username")).Password
        log.Println(dbPassword)
        // フォームから取得したユーザーパスワード
        formPassword := c.PostForm("password")

        // ユーザーパスワードの比較
        if err := plugins.CompareHashAndPassword(dbPassword, formPassword); err != nil {
            log.Println("ログインできませんでした")
            c.HTML(http.StatusBadRequest, "login.html", gin.H{"err": err})
            c.Abort()
        } else {
            log.Println("ログインできました")
            c.Redirect(302, "/")
        }
    })

    // ユーザー登録画面
    router.GET("/signup", func(c *gin.Context) {
        c.HTML(200, "signup.html", gin.H{})
    })

    // ユーザー登録
    router.POST("/signup", func(c *gin.Context) {
        var form model.User
        // バリデーション処理
        if err := c.Bind(&form); err != nil {
            c.HTML(http.StatusBadRequest, "signup.html", gin.H{"err": err})
            c.Abort()
        } else {
            username := c.PostForm("username")
            password := c.PostForm("password")
            // 登録ユーザーが重複していた場合にはじく処理
            if err := methods.CreateUser(username, password); err != nil {
                c.HTML(http.StatusBadRequest, "signup.html", gin.H{"err": err})
            }
            c.Redirect(302, "/login")
        }
    })

    router.Run()
}
