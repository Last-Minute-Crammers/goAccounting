package userService

import (
	"goAccounting/global"
	"goAccounting/global/constant"
	"goAccounting/global/db"
	userModel "goAccounting/internal/model/user"
	commonService "goAccounting/internal/service/common"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/net/context"
)

type User struct{}

func (userSvc *User) Login(email string, password string, ctx context.Context) (
	user userModel.User, token string, customClaims jwt.RegisteredClaims, err error,
) {
	password, err = commonService.Common.HashPassword(password)
	if err != nil {
		return
	}
	err = global.GlobalDb.Where("email = ? And password = ?", email, password).First(&user).Error
	if err != nil {
		return
	}
	customClaims = commonService.Common.MakeCustomClaims(user.ID)
	token, err = commonService.Common.GenerateJWT(customClaims)
	if err != nil {
		return
	}
	err = userSvc.updateDataAfterLogin(user, ctx)
	if err != nil {
		return
	}
	return
}

func (userSvc *User) updateDataAfterLogin(user userModel.User, ctx context.Context) error {
	err := db.GetDb(ctx).Model(user).Where("user_id = ?", user.ID).Update(
		"login_time", time.Now(),
	).Error
	if err != nil {
		return err
	}
	_, err = userSvc.RecordAction(user, constant.Login, ctx)
	if err != nil {
		return err
	}
	return nil
}

func (userSvc *User) Register(user userModel.User, ctx context.Context) (userModel.User, error) {
	dao := userModel.NewDao(db.Get(ctx))
	err := dao.CheckEmail(user.Email)
	if err != nil {
		return user, err
	}
	user.Password, err = commonService.Common.HashPassword(user.Password)
	if err != nil {
		return user, err
	}
	user, err = dao.AddUser(userModel.AddData{
		Username: user.Username,
		Password: user.Password,
		Email:    user.Email,
	})
	if err != nil {
		return user, err
	}
	return user, nil
}

func (userSvc *User) RecordAction(user userModel.User, action constant.UserAction, ctx context.Context) (
	*userModel.Log, error,
) {
	dao := userModel.NewLogDao(db.Get(ctx))
	log, err := dao.Add(user, &userModel.LogAddData{Action: action})
	if err != nil {
		return nil, err
	}
	return log, err
}
