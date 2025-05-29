package categoryService

import (
	"context"
	"errors"
	"goAccounting/global"
	"goAccounting/global/constant"
	"goAccounting/global/db"
	categoryModel "goAccounting/internal/model/category"

	"gorm.io/gorm"
)

type CategoryService struct {
	dao *categoryModel.CategoryDao
}

func NewCategoryService() *CategoryService {
	return &CategoryService{
		dao: categoryModel.NewDao(),
	}
}

func (cs *CategoryService) GetCategories(ie *constant.IncomeExpense, ctx context.Context) ([]categoryModel.Category, error) {
	cs.dao = categoryModel.NewDao(db.GetDb(ctx))
	return cs.dao.List(ie)
}

func (cs *CategoryService) GetCategoryById(id uint, ctx context.Context) (categoryModel.Category, error) {
	cs.dao = categoryModel.NewDao(db.GetDb(ctx))
	return cs.dao.SelectById(id)
}

func (cs *CategoryService) CreateCategory(data categoryModel.CategoryCreateData, ctx context.Context) (categoryModel.Category, error) {
	cs.dao = categoryModel.NewDao(db.GetDb(ctx))
	return cs.dao.Create(data)
}

func (cs *CategoryService) UpdateCategory(id uint, data categoryModel.CategoryUpdateData, ctx context.Context) error {
	cs.dao = categoryModel.NewDao(db.GetDb(ctx))
	return cs.dao.Update(id, data)
}

func (cs *CategoryService) DeleteCategory(id uint, ctx context.Context) error {
	cs.dao = categoryModel.NewDao(db.GetDb(ctx))
	return cs.dao.Delete(id)
}

// Create 创建新的分类
func (cs *CategoryService) Create(
	name string, incomeExpense constant.IncomeExpense) (categoryModel.Category, error) {

	// 1. 参数验证
	if name == "" {
		return categoryModel.Category{}, errors.New("分类名称不能为空")
	}

	// 2. 验证收支类型是否有效
	if incomeExpense != "收入" && incomeExpense != "支出" {
		return categoryModel.Category{}, errors.New("收支类型必须是'收入'或'支出'")
	}

	// 3. 创建数据结构
	data := categoryModel.CategoryCreateData{
		Name:          name,
		IncomeExpense: incomeExpense,
		Icon:          "", // 暂时不设置图标
	}

	// 4. 调用数据访问层创建分类
	return cs.CreateCategory(data, context.Background())
}

// Update 更新分类信息
func (cs *CategoryService) Update(
	id uint, name string, incomeExpense constant.IncomeExpense,
	ctx context.Context) error {

	// 1. 验证ID是否有效
	if id == 0 {
		return errors.New("分类ID不能为空")
	}

	// 2. 验证收支类型是否有效
	if incomeExpense != "收入" && incomeExpense != "支出" {
		return errors.New("收支类型必须是'收入'或'支出'")
	}

	// 3. 先查询要更新的分类是否存在
	category, err := cs.GetCategoryById(id, ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("找不到指定的分类")
	} else if err != nil {
		return errors.New("查询分类时发生错误: " + err.Error())
	}

	// 4. 如果要更新名称，先验证名称
	if name != "" {
		category.Name = name
		// 使用正确的数据库连接进行验证
		if err := category.CheckName(db.GetDb(ctx)); err != nil {
			return err
		}
	}

	// 5. 准备更新数据
	updateData := categoryModel.CategoryUpdateData{
		Name: &category.Name,
	}

	// 6. 执行更新操作
	return cs.UpdateCategory(id, updateData, ctx)
}

// Delete 删除分类（软删除）
func (cs *CategoryService) Delete(id uint, ctx context.Context) error {
	// 1. 验证ID是否有效
	if id == 0 {
		return errors.New("分类ID不能为空")
	}

	// 2. 先检查分类是否存在
	_, err := cs.GetCategoryById(id, ctx)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("找不到指定的分类")
	} else if err != nil {
		return errors.New("查询分类时发生错误: " + err.Error())
	}

	// 3. 执行软删除
	return cs.DeleteCategory(id, ctx)
}

// List 获取分类列表
func (cs *CategoryService) List(incomeExpense *constant.IncomeExpense, ctx context.Context) ([]categoryModel.Category, error) {
	// 1. 调用数据访问层获取列表
	return cs.GetCategories(incomeExpense, ctx)
}
