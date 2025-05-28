package util

import (
	"errors"
	"goAccounting/global/db"
	"goAccounting/internal/api/response"
	userModel "goAccounting/internal/model/user"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type checkFunc struct{}

var CheckFunc Checker = new(checkFunc)

type Checker interface {
	FriendInvitationBelongAndGet(id interface{}, ctx *gin.Context) (bool, userModel.FriendInvitation)
}

func (ckf *checkFunc) FriendInvitationBelongAndGet(id interface{}, ctx *gin.Context) (
	pass bool, data userModel.FriendInvitation,
) {
	err := db.Db.First(&data, id).Error
	if err != nil {
		response.FailToError(ctx, err)
		return
	}
	currentUserId := ContextFunc.GetUserId(ctx)
	if data.Inviter != currentUserId && data.Invitee != currentUserId {
		response.Forbidden(ctx)
		return
	}
	return true, data
}
func (ckf *checkFunc) handelForbiddenOrError(err error, ctx *gin.Context) (pass bool) {
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			response.Forbidden(ctx)
			return true
		}
		response.FailToError(ctx, err)
		return true
	}
	return false
}
