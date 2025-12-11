package v1

import (
	"goAccounting/global"
	"goAccounting/internal/api/request"
	"goAccounting/internal/api/response"
	userModel "goAccounting/internal/model/user"

	"github.com/gin-gonic/gin"
)

type AdminApi struct{}

// ListUsers 查询所有用户（支持 Query 参数和 JSON Body）
func (a *AdminApi) ListUsers(ctx *gin.Context) {
	var filter request.UserQueryFilter

	// 首先尝试从 JSON Body 解析
	if ctx.ContentType() == "application/json" {
		if err := ctx.ShouldBindJSON(&filter); err != nil {
			response.FailToError(ctx, err)
			return
		}
	} else {
		// 否则从 Query 参数解析
		if err := ctx.ShouldBindQuery(&filter); err != nil {
			response.FailToError(ctx, err)
			return
		}
	}

	query := global.GlobalDb.Model(&userModel.User{})

	// 过滤条件
	if filter.Email != "" {
		query = query.Where("email LIKE ?", "%"+filter.Email+"%")
	}
	if filter.Username != "" {
		query = query.Where("username LIKE ?", "%"+filter.Username+"%")
	}
	if filter.IsAdmin != nil {
		query = query.Where("is_admin = ?", *filter.IsAdmin)
	}

	// 分页
	var total int64
	query.Count(&total)

	// 设置默认分页值
	page := filter.Page
	pageSize := filter.PageSize
	if page < 1 {
		page = 1
	}
	if pageSize < 1 {
		pageSize = 10
	}
	if pageSize > 100 {
		pageSize = 100
	}

	var users []userModel.User
	err := query.
		Offset((page - 1) * pageSize).
		Limit(pageSize).
		Find(&users).Error

	if err != nil {
		response.FailToError(ctx, err)
		return
	}

	response.Success(ctx, gin.H{
		"users":    users,
		"total":    total,
		"page":     page,
		"pageSize": pageSize,
	})
}

// DeleteUser 删除用户
func (a *AdminApi) DeleteUser(ctx *gin.Context) {
	var req struct {
		UserID uint `json:"userId" binding:"required"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailToError(ctx, err)
		return
	}

	// 软删除
	err := global.GlobalDb.Delete(&userModel.User{}, req.UserID).Error
	if err != nil {
		response.FailToError(ctx, err)
		return
	}

	response.Success(ctx, nil)
}

// UpdateUserAdmin 设置/取消管理员
func (a *AdminApi) UpdateUserAdmin(ctx *gin.Context) {
	var req struct {
		UserID  uint `json:"userId" binding:"required"`
		IsAdmin bool `json:"isAdmin"`
	}
	if err := ctx.ShouldBindJSON(&req); err != nil {
		response.FailToError(ctx, err)
		return
	}

	err := global.GlobalDb.Model(&userModel.User{}).
		Where("id = ?", req.UserID).
		Update("is_admin", req.IsAdmin).Error

	if err != nil {
		response.FailToError(ctx, err)
		return
	}

	response.Success(ctx, nil)
}
