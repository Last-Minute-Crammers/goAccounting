package aiModel

import "goAccounting/global/db"

func init() {
	tables := []any{
		ChatRecord{}, // AI聊天记录表
	}
	err := db.InitDb.AutoMigrate(tables...)
	if err != nil {
		panic(err)
	}
}
