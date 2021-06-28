package funcDB

import (
	"fmt"
	"os"

	"example.com/m/v2/model"
	plugins "example.com/m/v2/plugins/crypto"
	_ "github.com/go-sql-driver/mysql" //直接的な記述が無いが、インポートしたいものに対しては"_"を頭につける決まり
	"github.com/jinzhu/gorm"
	"github.com/joho/godotenv"
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
	db.AutoMigrate(&model.Tweet{})                //構造体に基づいてテーブルを作成
	db.AutoMigrate(&model.User{})                 //構造体に基づいてテーブルを作成
	db.AutoMigrate(&model.CreateMiddleNameInit{}) //構造体に基づいてテーブルを作成

	// テーブルの初期化
	db.Exec("DROP TABLE mrs")
	db.Exec("DROP TABLE sns")
	db.Exec("DROP TABLE cns")

	db.AutoMigrate(&model.Mr{}) //構造体に基づいてテーブルを作成
	db.AutoMigrate(&model.SN{}) //構造体に基づいてテーブルを作成
	db.AutoMigrate(&model.CN{}) //構造体に基づいてテーブルを作成
	DbInsertSeed()
}

func DbInsertSeed() {
	db := gormConnect()

	defer db.Close()
	// Insert処理
	// TODO なんとかしてまとめたいなぁ〜
	mr := []string{"平", "源"}
	for _, m := range mr {
		db.Create(&model.Mr{Mr: m})
	}

	surname := []string{"朝臣", "臣", "国造", "県主", "和気", "稲置", "連", "直", "首", "史", "村主", "真人", "宿禰", "忌寸", "道師"}
	for _, s := range surname {
		db.Create(&model.SN{SurName: s})
	}

	commonname := []string{"一朗", "二郎", "三郎", "四郎", "五郎", "六郎", "七郎", "八郎", "九郎", "十朗", "十四郎", "二郎三郎"}
	for _, c := range commonname {
		db.Create(&model.CN{CommonName: c})
	}
	// TODO なんとかしてまとめたいなぁ〜

}

// データインサート処理
func DbInsert(content string) {
	db := gormConnect()

	defer db.Close()
	// Insert処理
	db.Create(&model.Tweet{Content: content})
}

//DB更新
func DbSessionUpdate(password string, session string) {
	db := gormConnect()
	var user model.User
	db.Where("password = ?", password).First(&user)
	user.Session, _ = plugins.PasswordEncrypt(session)
	db.Where("password = ?", password).Save(&user)
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
func CreateUser(userid string, username string, password string, email string, session string) []error {
	passwordEncrypt, _ := plugins.PasswordEncrypt(password)
	db := gormConnect()
	defer db.Close()
	// Insert処理
	if err := db.Create(&model.User{UserUUID: userid, Username: username, Password: passwordEncrypt, Email: email, Session: session}).GetErrors(); err != nil {
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

// ユーザーを一件取得(from uuid)
func GetUserFromUuid(uuid string) model.User {
	db := gormConnect()
	var user model.User
	db.First(&user, "user_uuid = ?", uuid)
	db.Close()
	return user
}
