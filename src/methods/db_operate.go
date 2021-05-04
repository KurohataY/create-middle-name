package methods

import (
    "os"
    _ "github.com/go-sql-driver/mysql" //直接的な記述が無いが、インポートしたいものに対しては"_"を頭につける決まり
    "github.com/jinzhu/gorm"
    "github.com/joho/godotenv"
    "example.com/m/v2/model"
    "example.com/m/v2/plugins/crypto"
    "fmt"
)

func godotenvConnect() {
  err := godotenv.Load(fmt.Sprintf("%s.env", os.Getenv("GO_ENV")))
  if err != nil {
    panic(err.Error())
  }
}

func gormConnect() *gorm.DB {

  godotenvConnect()

  DBMS := os.Getenv("DBMS")
  USER := os.Getenv("MYSQL_USER")
  PASS := os.Getenv("MYSQL_PASS")
  PROTOCOL := os.Getenv("MYSQL_PROTOCOL")
  DBNAME := os.Getenv("DBNAME")

  CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME + "?charset=utf8&parseTime=true&loc=Asia%2FTokyo"

  db, err := gorm.Open(DBMS, CONNECT)

  if err != nil {
      panic(err.Error())
  }
  return db
}

// DBの初期化
func DbInit() {
    db := gormConnect()

    // コネクション解放解放
    defer db.Close()
    db.AutoMigrate(&model.Tweet{}) //構造体に基づいてテーブルを作成
    db.AutoMigrate(&model.User{}) //構造体に基づいてテーブルを作成
}

// データインサート処理
func DbInsert(content string) {
    db := gormConnect()

    defer db.Close()
    // Insert処理
    db.Create(&model.Tweet{Content: content})
}

//DB更新
func DbUpdate(id int, tweetText string) {
    db := gormConnect()
    var tweet model.Tweet
    db.First(&tweet, id)
    tweet.Content = tweetText
    db.Save(&tweet)
    db.Close()
}

// 全件取得
func DbGetAll() []model.Tweet {
    db := gormConnect()

    defer db.Close()
    var tweets []model.Tweet
    // FindでDB名を指定して取得した後、orderで登録順に並び替え
    db.Order("created_at desc").Find(&tweets)
    return tweets
}

//DB一つ取得
func DbGetOne(id int) model.Tweet {
    db := gormConnect()
    var tweet model.Tweet
    db.First(&tweet, id)
    db.Close()
    return tweet
}

//DB削除
func DbDelete(id int) {
    db := gormConnect()
    var tweet model.Tweet
    db.First(&tweet, id)
    db.Delete(&tweet)
    db.Close()
}


// ユーザー登録処理
func CreateUser(username string, password string) []error {
    passwordEncrypt, _ := plugins.PasswordEncrypt(password)
    db := gormConnect()
    defer db.Close()
    // Insert処理
    if err := db.Create(&model.User{Username: username, Password: passwordEncrypt}).GetErrors(); err != nil {
        return err
    }
    return nil

}

// ユーザーを一件取得
func GetUser(username string) model.User {
    db := gormConnect()
    var user model.User
    db.First(&user, "username = ?", username)
    db.Close()
    return user
}
