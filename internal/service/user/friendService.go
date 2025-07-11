package userService

import (
	"errors"
	"goAccounting/global/ctxutil"
	"goAccounting/global/db"
	userModel "goAccounting/internal/model/user"

	"golang.org/x/net/context"
	"gorm.io/gorm"
)

type Friend struct{}

func (f *Friend) CreateInvitation(
	inviter userModel.User, invitee userModel.User, ctx context.Context,
) (invitation userModel.FriendInvitation, err error) {
	return invitation, db.ExecInTransaction(
		ctx, func(ctx *ctxutil.TxContext) (err error) {
			tx := ctx.GetDb()
			dao := userModel.NewDao(tx)
			invitation, err = dao.CreateFriendInvitation(inviter.ID, invitee.ID)
			if !(errors.Is(err, gorm.ErrDuplicatedKey)) {
				return
			}
			// 处理重复键
			invitation, err = dao.SelectFriendInvitation(inviter.ID, invitee.ID, true)
			if err != nil {
				return
			}
			var isRealFriend bool
			isRealFriend, err = dao.IsRealFriend(inviter.ID, invitee.ID)
			if isRealFriend || err != nil {
				return
			}
			if invitation.Status == userModel.InvitationStatsOfWaiting {
				return
			}
			return invitation.UpdateStatus(userModel.InvitationStatsOfWaiting, tx)
		},
	)
}

func (f *Friend) AcceptInvitation(Invitation *userModel.FriendInvitation, ctx context.Context) (
	inviterFriend userModel.Friend, inviteeFriend userModel.Friend, err error,
) {
	err = db.ExecInTransaction(
		ctx, func(ctx *ctxutil.TxContext) error {
			inviterFriend, inviteeFriend, err = Invitation.Accept(ctx.GetDb())
			return err
		},
	)
	return
}

func (f *Friend) RefuseInvitation(Invitation *userModel.FriendInvitation, ctx context.Context) error {
	return db.ExecInTransaction(
		ctx, func(ctx *ctxutil.TxContext) error {
			return Invitation.Refuse(ctx.GetDb())
		},
	)
}
