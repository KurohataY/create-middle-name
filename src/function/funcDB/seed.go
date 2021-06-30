package funcDB

import (
	"example.com/m/v2/model"
	_ "github.com/go-sql-driver/mysql" //直接的な記述が無いが、インポートしたいものに対しては"_"を頭につける決まり
)

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
