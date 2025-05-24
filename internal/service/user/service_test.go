package userService_test

import (
	"context"
	"goAccounting/global"
	userModel "goAccounting/internal/model/user"
	userService "goAccounting/internal/service/user"
	"strconv"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser_RegisterAndLogin(t *testing.T) {
	ctx := context.Background()
	userSvc := userService.User{}

	// 清理数据库
	global.GlobalDb.Exec("DELETE FROM user")

	// 测试注册
	newUser := userModel.User{
		Username: "testuser",
		Email:    "newuser@example.com",
		Password: "password123",
	}

	registeredUser, err := userSvc.Register(newUser, ctx)
	assert.NoError(t, err)
	assert.Equal(t, newUser.Email, registeredUser.Email)
	assert.NotEmpty(t, registeredUser.Password)
	assert.NotEqual(t, newUser.Password, registeredUser.Password) // 确保密码已加密

	// 测试登录
	loggedInUser, token, claims, err := userSvc.Login(newUser.Email, "password123", ctx)
	assert.NoError(t, err)
	assert.Equal(t, registeredUser.ID, loggedInUser.ID)
	assert.NotEmpty(t, token)

	// 将 claims.Subject 转换为 uint
	claimsUserID, err := strconv.ParseUint(claims.Subject, 10, 64)
	assert.NoError(t, err)
	assert.Equal(t, uint(registeredUser.ID), uint(claimsUserID)) // 修复断言
}
