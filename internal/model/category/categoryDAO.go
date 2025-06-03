package categoryModel

import (
	"goAccounting/global"
	"goAccounting/global/constant"

	"gorm.io/gorm"
)

type CategoryDao struct {
	db *gorm.DB
}

func NewDao(db ...*gorm.DB) *CategoryDao {
	if len(db) > 0 {
		return &CategoryDao{db: db[0]}
	}
	return &CategoryDao{global.GlobalDb}
}

func (c *CategoryDao) SelectById(id uint) (Category, error) {
	var category Category
	err := c.db.Where("id = ?", id).First(&category).Error
	return category, err
}

func (c *CategoryDao) List(ie *constant.IncomeExpense) ([]Category, error) {
	var categories []Category
	query := c.db.Model(&Category{})
	
	if ie != nil {
		query = query.Where("income_expense = ?", *ie)
	}
	
	err := query.Find(&categories).Error
	return categories, err
}

type CategoryCreateData struct {
	Name          string
	IncomeExpense constant.IncomeExpense
	Icon          string
	UserId        *uint
}

func (c *CategoryDao) Create(data CategoryCreateData) (Category, error) {
	category := Category{
		Name:          data.Name,
		IncomeExpense: data.IncomeExpense,
		Icon:          data.Icon,
		UserId:        data.UserId,
	}
	
	err := c.db.Create(&category).Error
	return category, err
}

type CategoryUpdateData struct {
	Name *string
	Icon *string
}

func (c *CategoryDao) Update(id uint, data CategoryUpdateData) error {
	updates := make(map[string]interface{})
	
	if data.Name != nil {
		updates["name"] = *data.Name
	}
	if data.Icon != nil {
		updates["icon"] = *data.Icon
	}
	
	return c.db.Model(&Category{}).Where("id = ?", id).Updates(updates).Error
}

func (c *CategoryDao) Delete(id uint) error {
	return c.db.Delete(&Category{}, id).Error
}
