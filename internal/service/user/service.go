package userService

import (
	"goAccounting/global"
	"goAccounting/global/constant"
	"goAccounting/global/ctxutil"
	"goAccounting/global/db"
	categoryModel "goAccounting/internal/model/category"
	userModel "goAccounting/internal/model/user"
	commonService "goAccounting/internal/service/common"
	"log"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"golang.org/x/net/context"
)

type User struct{}

func (userSvc *User) Login(email string, password string, ctx context.Context) (
	user userModel.User, token string, customClaims jwt.RegisteredClaims, err error,
) {
	// 仅通过 email 查询用户记录
	err = global.GlobalDb.Where("email = ?", email).First(&user).Error
	if err != nil {
		return
	}

	// 验证密码
	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return
	}

	// 生成 JWT 和更新登录数据
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
	// err := db.GetDb(ctx).Model(&user).Where("id = ?", user.ID).Update(
	// 	"login_time", time.Now(),
	// ).Error
	// if err != nil {
	// 	return err
	// }
	// tag:consider-user_log
	_, err := userSvc.RecordAction(user, constant.Login, ctx)
	if err != nil {
		return err
	}

	return nil
}

func (userSvc *User) Register(data userModel.AddData, ctx context.Context) (userModel.User, error) {
	log.Println("[service]: starting register")
	log.Printf("[service]: userService.Register param: %+v\n", data)
	dao := userModel.NewDao(db.Get(ctx))

	err := dao.CheckEmail(data.Email)
	if err != nil {
		return userModel.User{}, err
	}

	// 正确 hash 密码
	hashedPW, err := commonService.Common.HashPassword(data.Password)
	if err != nil {
		return userModel.User{}, err
	}
	data.Password = hashedPW

	log.Printf("[service]: addData param: %+v\n", data)
	user, err := dao.AddUser(data)
	if err != nil {
		return user, err
	}

	err = categoryModel.CreateDefaultCategoriesForUser(user.ID)
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

// set avatar for user (not now)

// set username for user
func (userSvc *User) SetUsername(userId uint, username string, ctx context.Context) error {
	dao := userModel.NewDao(db.Get(ctx))
	return dao.UpdateUsername(userId, username)
}

func (userSvc *User) UpdatePassword(user userModel.User, newPassword string, ctx context.Context) error {
	password, _ := commonService.Common.HashPassword(newPassword)
	if password == user.Password {
		logRemark := global.ErrSameAsTheOldPassword.Error()
		log.Printf("%s\n", logRemark)
	}
	return db.ExecInTransaction(
		ctx, func(ctx *ctxutil.TxContext) error {
			tx := ctx.GetDb()
			err := tx.Model(user).Update("password", password).Error
			if err != nil {
				return err
			}
			return nil
		},
	)
}

func (userSvc *User) UpdateInfo(user userModel.User, username string, ctx context.Context) error {
	return db.Get(ctx).Model(&user).Update("username", username).Error
}
