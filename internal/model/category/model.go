package categoryModel

import (
	"goAccounting/global"
	"goAccounting/global/constant"
	"time"

	"gorm.io/gorm"
)

type Category struct {
	ID            uint                   `gorm:"primaryKey"`
	IncomeExpense constant.IncomeExpense `gorm:"comment:'收支类型'"`
	Name          string                 `gorm:"comment:'标签名';size:128;uniqueIndex"`
	Icon          string                 `gorm:"comment:'图标';size:64"`
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

func (c *Category) CheckName(_ *gorm.DB) error {
	if c.Name == "" {
		return global.NewErrDataIsEmpty("交易类型名称")
	}
	return nil
}
