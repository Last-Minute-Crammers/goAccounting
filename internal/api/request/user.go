package request

// here, we have something about accpect and check the
// parameters from client, so we can call this part by "DTO"
// which means Data Transfer Object

type UserLogin struct {
	Email    string `binding:"required"`
	Password string `binding:"required"`
	// PicCaptcha
}

type UserRegister struct {
	Username string `binding:"required"`
	Password string `binding:"required"`
	Email    string `binding:"required,email"`
	Captcha  string `binding:"required"`
}

type UserUpdateInfo struct {
	Username string `binding:"required"`
}

type UserSearch struct {
	Id       *uint  `binding:"omitempty"`
	Username string `binding:"required"`
	PageData
}

// the following is not that important
type UserForgetPassword struct {
	Email    string `binding:"required,email"`
	Password string `binding:"required"`
	Captcha  string `binding:"required"`
}

type UserUpdatePassword struct {
	Password string `binding:"required"`
	Captcha  string `binding:"required"`
}

// we only have violent way to shutdown share for all frieneds
type UserTransactionShareConfigUpdate struct {
	Status bool
}

type UserCreateFriendInvitation struct {
	Invitee uint
}

type UserGetFriendInvitation struct {
	IsInvite bool
}

type UserGetAccountInvitationList struct {
	PageData
}
