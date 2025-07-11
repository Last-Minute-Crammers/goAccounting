package categoryService

import (
	"errors"
	"goAccounting/global/constant"
	"goAccounting/global/db"
	categoryModel "goAccounting/internal/model/category"
	"log"

	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type Category struct{}

// Create 创建新的分类
func (c *Category) Create(
	name string, incomeExpense constant.IncomeExpense, icon string, ctx context.Context) (categoryModel.Category, error) {

	// 1. 参数验证
	if name == "" {
		return categoryModel.Category{}, errors.New("分类名称不能为空")
	}

	// 2. 验证收支类型是否有效 - 使用英文常量
	if incomeExpense != constant.Income && incomeExpense != constant.Expense {
		return categoryModel.Category{}, errors.New("收支类型必须是'income'或'expense'")
	}

	// 3. 从上下文获取当前用户ID作为AccountID
	// 假设你的鉴权中间件会将用户ID存储在context中
	userId, exists := ctx.Value("userId").(uint)
	if !exists || userId == 0 {
		return categoryModel.Category{}, errors.New("用户未登录")
	}
	log.Printf("[tagsService] :user Id = %d\n", userId)

	// 4. 创建数据结构
	data := categoryModel.CategoryCreateData{
		AccountID:     userId,
		Name:          name,
		IncomeExpense: incomeExpense,
		Icon:          icon,
	}

	// 5. 调用数据访问层创建分类
	return categoryModel.NewDao(db.GetDb(ctx)).Create(data)
}

// Update 更新分类信息
func (c *Category) Update(
	id uint, name string, incomeExpense constant.IncomeExpense,
	ctx context.Context) error {

	// 1. 验证ID是否有效
	if id == 0 {
		return errors.New("分类ID不能为空")
	}

	// 2. 验证收支类型是否有效 - 使用英文常量
	if incomeExpense != constant.Income && incomeExpense != constant.Expense {
		return errors.New("收支类型必须是'income'或'expense'")
	}

	// 3. 获取数据库连接
	cDAO := categoryModel.NewDao(db.GetDb(ctx))

	// 4. 先查询要更新的分类是否存在
	category, err := cDAO.SelectById(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("找不到指定的分类")
	} else if err != nil {
		return errors.New("查询分类时发生错误: " + err.Error())
	}

	// 5. 如果要更新名称，先验证名称
	if name != "" {
		category.Name = name
		// 使用正确的数据库连接进行验证
		if err := category.CheckName(db.GetDb(ctx)); err != nil {
			return err
		}
	}

	// 6. 准备更新数据
	updateData := categoryModel.CategoryUpdateData{
		Name: &category.Name,
	}

	// 7. 执行更新操作
	return cDAO.Update(id, updateData)
}

// Delete 删除分类（软删除）
func (c *Category) Delete(id uint, ctx context.Context) error {
	// 1. 验证ID是否有效
	if id == 0 {
		return errors.New("分类ID不能为空")
	}

	// 2. 获取数据库连接
	cDAO := categoryModel.NewDao(db.GetDb(ctx))

	// 3. 先检查分类是否存在
	_, err := cDAO.SelectById(id)
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return errors.New("找不到指定的分类")
	} else if err != nil {
		return errors.New("查询分类时发生错误: " + err.Error())
	}

	// 4. 执行软删除
	return cDAO.Delete(id)
}

// List 获取分类列表
func (c *Category) List(incomeExpense *constant.IncomeExpense, ctx context.Context) ([]categoryModel.Category, error) {
	// 从上下文获取当前用户ID
	userId, exists := ctx.Value("userId").(uint)
	if !exists || userId == 0 {
		return nil, errors.New("用户未登录")
	}

	// 1. 获取数据库连接
	cDAO := categoryModel.NewDao(db.GetDb(ctx))

	// 2. 调用数据访问层获取列表
	return cDAO.List(&userId, incomeExpense)
}
