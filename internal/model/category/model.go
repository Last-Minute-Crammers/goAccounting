package categoryModel

import (
	"errors"
	"goAccounting/global/constant"
	commonModel "goAccounting/internal/model/common"
	"time"

	"gorm.io/gorm"
)

type Category struct {
	commonModel.BaseModel
	ID            uint                   `gorm:"primarykey"`
	Name          string                 `gorm:"type:varchar(128);not null"`
	Icon          string                 `gorm:"type:varchar(64)"`
	IncomeExpense constant.IncomeExpense `gorm:"type:varchar(16);not null"`
	UserId        *uint                  `gorm:"index"` // null表示系统默认分类
	CreatedAt     time.Time              `gorm:"type:TIMESTAMP"`
	UpdatedAt     time.Time              `gorm:"type:TIMESTAMP"`
	DeletedAt     gorm.DeletedAt         `gorm:"index;type:TIMESTAMP"`
}

func (c *Category) TableName() string {
	return "category"
}

func (c *Category) Delete() error {
	return global.GlobalDb.Delete(c).Error
}

func (c *Category) CheckName(db *gorm.DB) error {
	var count int64
	err := db.Model(&Category{}).Where("name = ? AND income_expense = ?", c.Name, c.IncomeExpense).Count(&count).Error
	if err != nil {
		return err
	}
	if count > 0 {
		return errors.New("分类名称已存在")
	}
	return nil
}

// GetDefaultCategories 获取系统默认分类
func GetDefaultCategories() []Category {
	return []Category{
		{ID: 1, Name: "餐饮", Icon: "🍔", IncomeExpense: constant.Expense, UserId: nil},
		{ID: 2, Name: "购物", Icon: "🛒", IncomeExpense: constant.Expense, UserId: nil},
		{ID: 3, Name: "交通", Icon: "🚗", IncomeExpense: constant.Expense, UserId: nil},
		{ID: 4, Name: "住房", Icon: "🏠", IncomeExpense: constant.Expense, UserId: nil},
		{ID: 5, Name: "娱乐", Icon: "🎬", IncomeExpense: constant.Expense, UserId: nil},
		{ID: 6, Name: "数码", Icon: "📱", IncomeExpense: constant.Expense, UserId: nil},
		{ID: 7, Name: "服饰", Icon: "👗", IncomeExpense: constant.Expense, UserId: nil},
		{ID: 8, Name: "其他", Icon: "➕", IncomeExpense: constant.Expense, UserId: nil},
		{ID: 9, Name: "工资", Icon: "💰", IncomeExpense: constant.Income, UserId: nil},
		{ID: 10, Name: "投资", Icon: "📈", IncomeExpense: constant.Income, UserId: nil},
		{ID: 11, Name: "兼职", Icon: "💼", IncomeExpense: constant.Income, UserId: nil},
		{ID: 12, Name: "其他收入", Icon: "➕", IncomeExpense: constant.Income, UserId: nil},
	}
}
