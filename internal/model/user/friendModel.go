package userModel

import (
	"errors"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

// Friend means user of userId's Frined
type Friend struct {
	ID        uint `gorm:"primarykey"`
	UserId    uint `gorm:"uniqueIndex:idx_mapping,priority:1"`
	FriendId  uint `gorm:"uniqueIndex:idx_mapping,priority:2;"`
	AddMode   AddMode
	CreatedAt time.Time      `gorm:"type:TIMESTAMP"`
	UpdatedAt time.Time      `gorm:"type:TIMESTAMP"`
	DeletedAt gorm.DeletedAt `gorm:"index;type:TIMESTAMP"`
}

func (f *Friend) TableName() string {
	return "user_friend"
}

func (f *Friend) GetFriend(args ...any) (User, error) {
	return NewDao().SelectById(f.FriendId, args)
}

func (f *Friend) GetFriendInfo() (info UserInfo, err error) {
	return NewDao().SelectUserInfoById(f.FriendId)
}

type AddMode string

const FriendAddModeOfFriendInvitation AddMode = "frinedInvitation"

type FriendInvitationStatus int

const (
	InvitationStatsOfWaiting FriendInvitationStatus = iota
	InvitationStatsOfAccept
	InvitationStatsOfRefuse
)

type FriendInvitation struct {
	ID        uint `gorm:"primarykey"`
	Inviter   uint `gorm:"uniqueIndex:idx_mapping,priority:1"`
	Invitee   uint `gorm:"uniqueIndex:idx_mapping,priority:2"`
	Status    FriendInvitationStatus
	CreatedAt time.Time      `gorm:"type:TIMESTAMP"`
	UpdatedAt time.Time      `gorm:"type:TIMESTAMP"`
	DeletedAt gorm.DeletedAt `gorm:"index;type:TIMESTAMP"`
}

func (f *FriendInvitation) TableName() string {
	return "user_friend_invitation"
}

// ForUpadate is an exclusive lock that locks the queried records in
// current tx to prevent other txs fro modifying these records
func (f *FriendInvitation) ForUpdate(tx *gorm.DB) error {
	err := tx.Model(f).Clauses(clause.Locking{Strength: "UPDATE"}).First(f, f.ID).Error
	if err != nil {
		return err

	}
	return nil
}

// ForShare is a shared lock that allows other txs to read these
// records but not to modify or delete them.
func (f *FriendInvitation) ForShare(tx *gorm.DB) error {
	err := tx.Model(f).Clauses(clause.Locking{Strength: "SHARE"}).First(f, f.ID).Error
	if err != nil {
		return err
	}
	return nil
}

func (f *FriendInvitation) GetInviterInfo() (UserInfo, error) {
	return NewDao().SelectUserInfoById(f.Inviter)
}

func (f *FriendInvitation) GetInviteeInfo() (UserInfo, error) {
	return NewDao().SelectUserInfoById(f.Invitee)
}

func (f *FriendInvitation) AddFirend(tx *gorm.DB) (inviterFriend Friend, inviteeFriend Friend, err error) {
	err = f.ForShare(tx)
	if err != nil {
		return
	}
	inviterFriend, err = NewDao(tx).AddFriend(f.Inviter, f.Invitee, FriendAddModeOfFriendInvitation)
	if err != nil {
		return
	}
	inviteeFriend, err = NewDao(tx).AddFriend(f.Invitee, f.Inviter, FriendAddModeOfFriendInvitation)
	if err != nil {
		return
	}
	return

}

func (f *FriendInvitation) UpdateStatus(status FriendInvitationStatus, tx *gorm.DB) error {
	return tx.Model(f).Update("status", status).Error
}

func (f *FriendInvitation) Accept(tx *gorm.DB) (inviterFriend Friend, inviteeFriend Friend, err error) {
	err = f.ForShare(tx)
	if err != nil {
		return
	}
	if f.Status != InvitationStatsOfWaiting {
		err = errors.New("unexpected invitation status")
		return
	}
	err = f.UpdateStatus(InvitationStatsOfAccept, tx)
	if err != nil {
		return
	}
	return f.AddFirend(tx)
}

func (f *FriendInvitation) Refuse(tx *gorm.DB) error {
	err := f.ForShare(tx)
	if err != nil {
		return err
	}
	if f.Status != InvitationStatsOfWaiting {
		err = errors.New("unexpected invitation status")
		return err
	}
	return f.UpdateStatus(InvitationStatsOfRefuse, tx)
}
