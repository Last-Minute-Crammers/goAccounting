package categoryModel

import (
	"errors"
	"goAccounting/global"
	"goAccounting/global/constant"
	"goAccounting/util"

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

func (cd *CategoryDao) SelectById(id uint) (category Category, err error) {
	err = cd.db.First(&category, id).Error
	return
}

func (cd *CategoryDao) SelectByName(accountId uint, name string, incomeExpense constant.IncomeExpense) (category Category, err error) {
	err = cd.db.Where("account_id = ? AND name = ? AND income_expense = ?", accountId, name, incomeExpense).First(&category).Error
	return
}

type CategoryUpdateData struct {
	Name *string
	Icon *string
}

func (cd *CategoryDao) Update(categoryId uint, data CategoryUpdateData) error {
	updateData := &Category{}
	if err := util.Data.CopyNotEmptyStringOptional(data.Name, &updateData.Name); err != nil {
		return err
	}
	if err := util.Data.CopyNotEmptyStringOptional(data.Icon, &updateData.Icon); err != nil {
		return err
	}
	if updateData.Name != "" {
		if err := updateData.CheckName(cd.db); err != nil {
			return err
		}
	}
	err := cd.db.Model(&updateData).Where("id = ?", categoryId).Updates(updateData).Error
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return global.ErrCategorySameName
	}
	return err
}

type CategoryCreateData struct {
	AccountID     uint
	Name          string
	Icon          string
	IncomeExpense constant.IncomeExpense
}

func (cd *CategoryDao) Create(data CategoryCreateData) (Category, error) {
	category := Category{
		AccountID:     data.AccountID,
		Name:          data.Name,
		Icon:          data.Icon,
		IncomeExpense: data.IncomeExpense,
	}

	if err := category.CheckName(cd.db); err != nil {
		return Category{}, err
	}

	err := cd.db.Create(&category).Error
	if errors.Is(err, gorm.ErrDuplicatedKey) {
		return Category{}, global.ErrCategorySameName
	}

	return category, err
}

// Delete 根据分类ID删除分类（软删除）
func (cd *CategoryDao) Delete(categoryId uint) error {
	var category Category
	err := cd.db.First(&category, categoryId).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return errors.New("category not found")
		}
		return err
	}

	return cd.db.Delete(&Category{}, categoryId).Error
}

// HardDelete 完全删除分类（物理删除，慎用）
func (cd *CategoryDao) HardDelete(categoryId uint) error {
	return cd.db.Unscoped().Delete(&Category{}, categoryId).Error
}

// ListOptions 列表查询选项
type ListOptions struct {
	AccountID     *uint
	IncomeExpense *constant.IncomeExpense
	OrderBy       string // 排序字段，默认为 created_at
	OrderDesc     bool   // 是否降序排列，默认为 true
	Limit         int    // 限制返回记录数，默认不限制
	Offset        int    // 跳过记录数，用于分页
}

func (cd *CategoryDao) ListWithOptions(options ListOptions) (list []Category, err error) {
	query := cd.db.Model(&Category{})

	if options.AccountID != nil {
		query = query.Where("account_id = ?", *options.AccountID)
	}

	if options.IncomeExpense != nil {
		query = query.Where("income_expense = ?", *options.IncomeExpense)
	}

	orderBy := "created_at"
	if options.OrderBy != "" {
		orderBy = options.OrderBy
	}

	if options.OrderDesc {
		query = query.Order(orderBy + " DESC")
	} else {
		query = query.Order(orderBy)
	}

	if options.Limit > 0 {
		query = query.Limit(options.Limit)
	}

	if options.Offset > 0 {
		query = query.Offset(options.Offset)
	}

	err = query.Find(&list).Error
	return list, err
}

func (cd *CategoryDao) List(accountId *uint, ie *constant.IncomeExpense) (list []Category, err error) {
	options := ListOptions{
		AccountID:     accountId,
		IncomeExpense: ie,
		OrderDesc:     true,
		OrderBy:       "created_at",
	}
	return cd.ListWithOptions(options)
}
