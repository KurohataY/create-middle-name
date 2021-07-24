package database

import (
	"fmt"
	"os"

	"example.com/m/v2/model"
	plugins "example.com/m/v2/plugins/crypto"
	_ "github.com/go-sql-driver/mysql"
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
	CHARSET := os.Getenv("CHARSET")
	LOC := os.Getenv("LOC")

	CONNECT := USER + ":" + PASS + "@" + PROTOCOL + "/" + DBNAME + "?charset=" + CHARSET + "&parseTime=true&loc=" + LOC

	db, err := gorm.Open(DBMS, CONNECT)

	if err != nil {
		panic(err.Error())
	}
	return db
}

// DBの初期化
func Init() {
	db := gormConnect()

	// コネクション解放解放
	defer db.Close()
	db.AutoMigrate(&model.User{})
	db.AutoMigrate(&model.CreatedMiddleNames{})

	// テーブルの初期化
	db.Exec("DROP TABLE mrs")
	db.Exec("DROP TABLE sns")
	db.Exec("DROP TABLE cns")

	db.AutoMigrate(&model.Mr{})
	db.AutoMigrate(&model.SN{})
	db.AutoMigrate(&model.CN{})

	// DB初期値の挿入
	DbInsertSeed()
}

// データインサート処理
func DbMiddleNameInsert(mr string, lName string, sName string, cName string, fName string, userId string) []error {
	db := gormConnect()
	defer db.Close()
	// Insert処理
	if err := db.Create(&model.CreatedMiddleNames{Mr: mr, LName: lName, SurName: sName, CommonName: cName, FName: fName, UserId: userId}).GetErrors(); err != nil {
		return err
	}
	return nil
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
func DbGetCreatedMiddleNames(userId string) []model.CreatedMiddleNames {
	db := gormConnect()

	defer db.Close()
	var createdMiddleNames []model.CreatedMiddleNames

	db.Where("user_id =  ?", userId).Find(&createdMiddleNames)
	print(fmt.Sprintf("%v", createdMiddleNames))
	return createdMiddleNames
}

func DBGetRandomMrData() model.Mr {
	db := gormConnect()

	defer db.Close()
	var mr model.Mr
	// db.Raw("SELECT * FROM mrs ORDER BY RAND() LIMIT 1").Find(&mr)
	// db.Limit(1).Offset(5).Find(&users)
	db.Order("RAND()").Find(&mr)
	// SELECT * FROM users ORDER BY age desc, name;

	return mr
}

// //DB一つ取得
// func DbGetOne(id int) model.Tweet {
// 	db := gormConnect()
// 	var tweet model.Tweet
// 	db.First(&tweet, id)
// 	db.Close()
// 	return tweet
// }

// //DB削除
// func DbDelete(id int) {
// 	db := gormConnect()
// 	var tweet model.Tweet
// 	db.First(&tweet, id)
// 	db.Delete(&tweet)
// 	db.Close()
// }

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
