package response

import (
	userModel "goAccounting/internal/model/user"
	"time"
)

type UserOne struct {
	Id         uint
	Username   string
	Email      string
	CreateTime time.Time
}

type Login struct {
	Token               string
	TokenExpirationTime time.Time
	User                UserOne
}
type Register struct {
	User                UserOne
	Token               string
	TokenExpirationTime time.Time
}

func (u *UserOne) SetData(data userModel.User) error {
	u.Id = data.ID
	u.Email = data.Email
	u.Username = data.Username
	u.CreateTime = data.CreatedAt
	return nil
}

type UserHomeHeaderCard struct {
	*TransactionStatistic
}

type UserHome struct {
	HeaderCard *UserHomeHeaderCard
}

type UserFriendInvitation struct {
	Id         uint
	Inviter    UserInfo
	Invitee    UserInfo
	CreateTime time.Time
}

type UserInfo struct {
	Id       uint
	Username string
	Email    string
}

func (u *UserInfo) SetData(data userModel.UserInfo) {
	u.Id = data.ID
	u.Username = data.Username
	u.Email = data.Email
}
