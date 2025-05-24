package userModel

import "goAccounting/global/db"

func init() {
	tables := []any{
		User{},
		Friend{}, FriendInvitation{},
		TransactionShareConfig{},
		Log{}, // 添加 Log 模型
	}
	err := db.InitDb.AutoMigrate(tables...)
	if err != nil {
		panic(err)
	}
}
