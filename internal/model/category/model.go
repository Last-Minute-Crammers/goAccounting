package categoryModel

import (
	"goAccounting/global"
	"goAccounting/global/constant"
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID            uint                   `gorm:"primaryKey"`
	UserID        uint                   `gorm:"comment:'账户ID';index;uniqueIndex:idx_category_unique,priority:1"`
	IncomeExpense constant.IncomeExpense `gorm:"comment:'收支类型';uniqueIndex:idx_category_unique,priority:3"`
	Name          string                 `gorm:"comment:'标签名';size:128;uniqueIndex:idx_category_unique,priority:2"`
	Icon          string                 `gorm:"comment:'图标';size:64"`
	CreatedAt     time.Time              `gorm:"type:TIMESTAMP"`
	UpdatedAt     time.Time              `gorm:"type:TIMESTAMP"`
	DeletedAt     gorm.DeletedAt         `gorm:"index;type:TIMESTAMP"`
}

// Add composite unique index
func (c *Category) BeforeCreate(tx *gorm.DB) error {
	// Create composite unique index: account_id + name + income_expense
	return nil
}

func (c *Category) TableName() string {
	return "category"
}

func (c *Category) Delete() error {
	return global.GlobalDb.Delete(c).Error
}

func (c *Category) CheckName(_ *gorm.DB) error {
	if c.Name == "" {
		return global.NewErrDataIsEmpty("交易类型名称")
	}
	return nil
}
