package userModel

import (
	"time"

	"gorm.io/gorm"
)

type Friend struct {
	ID        uint `gorm:"primarykey"`
	UserId    uint `gorm:"uniqueIndex:idx_mapping,priority:1"`
	FriendId  uint `gorm:"uniqueIndex:idx_mapping,priority:2;"`
	AddMode   AddMode
	CreatedAt time.Time      `gorm:"type:TIMESTAMP"`
	UpdatedAt time.Time      `gorm:"type:TIMESTAMP"`
	DeletedAt gorm.DeletedAt `gorm:"index;type:TIMESTAMP"`
}

type AddMode string

const FriendAddModeOfFriendInvitation AddMode = "frinedInvitation"

func (f *Friend) TableName() string {
	return "user_friend"
}

func (f *Friend) GetFriend(args ...any) (User, error) {
	return NewDao().SelectById(f.FriendId, args)
}

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
