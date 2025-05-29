package categoryModel

import (
	"goAccounting/global/db"
	"log"
)

func init() {
	tables := []any{
		Category{},
	}
	err := db.InitDb.AutoMigrate(tables...)
	if err != nil {
		panic(err)
	}

	// 初始化默认分类数据
	initDefaultCategories()
}

func initDefaultCategories() {
	defaultCategories := GetDefaultCategories()
	
	for _, category := range defaultCategories {
		var existingCategory Category
		result := db.InitDb.Where("id = ?", category.ID).First(&existingCategory)
		
		if result.Error != nil {
			// 分类不存在，创建新分类
			if err := db.InitDb.Create(&category).Error; err != nil {
				log.Printf("Failed to create default category %s: %v", category.Name, err)
			} else {
				log.Printf("Created default category: %s", category.Name)
			}
		}
	}
}
