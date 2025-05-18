package commentModel

import "goAccounting/global/db"

func init() {
	tables := []any{
		Comment{},
	}
	err := db.InitDb.AutoMigrate(tables...)
	if err != nil {
		panic(err)
	}
}
