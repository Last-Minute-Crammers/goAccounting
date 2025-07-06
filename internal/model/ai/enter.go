package aiModel

import "goAccounting/global/db"

func init() {
	tables := []any{
		ChatRecord{},     // AI聊天记录表
		FinancialReport{}, // AI财务报告表
	}
	err := db.InitDb.AutoMigrate(tables...)
	if err != nil {
		panic(err)
	}
}
