package userModel

//* DAO means Data Access Object, packages
//* all the oprations the database about user table

import (
	"errors"
	"goAccounting/global"

	"gorm.io/gorm"
)

type UserDao struct {
	db *gorm.DB
}

// UserFilter Conditions
type Condition struct {
	Id                 *uint
	LikePrefixUsername *string
	Offset             int
	Limit              int
}

type AddData struct {
	Username string
	Password string
	Email    string
}

// we user interface just for clearly to see all th function
// all the following function is for user's operation, so DAO
type _baseInterface interface {
	AddUser(data AddData) (User, error)
	SelectById(id uint, args ...interface{}) (User, error)
	PluckNameById(id uint) (string, error)
	SelectByEmail(email string) (User, error)
	SelectUserInfoById(id uint) (result UserInfo, err error)
	SelectUserInfoByCondition(condition Condition) ([]UserInfo, error)
}

type _friendInterface interface {
	CreateFriendInvitation(uint, uint) (FriendInvitation, error)
	// FOR UPDATE is a exclusive LOCK, before other
	// related trasactions are done, it cannot be modified
	SelectFriendInvitation(inviter uint, invitee uint, forUpdate bool) (result FriendInvitation, err error)
	// due to the case that : someone don't have any friend,
	// so we use uint pointer to avoid there's no data
	SelectFriendInvitationList(inviter *uint, invitee *uint) (result []FriendInvitation, err error)
	SelectFriend(uint, uint) (Friend, error)
	IsRealFriend(userId uint, friendId uint) (bool, error)
	AddFriend(uint, uint, AddMode)
	SelectFriendList(uint) ([]Friend, error)
}

func NewDao(db ...*gorm.DB) *UserDao {
	if len(db) > 0 {
		return &UserDao{
			db: db[0],
		}
	}
	return &UserDao{global.GlobalDb}
}

func (u *UserDao) AddUser(data AddData) (User, error) {
	user := User{
		Username: data.Username,
		Password: data.Password,
		Email:    data.Email,
	}
	err := u.db.Create(&user).Error
	return user, err
}

func (u *UserDao) SelectById(id uint, args ...any) (User, error) {
	user := User{}
	var err error
	if len(args) > 0 {
		err = u.db.Where("Id = ?", id).Select(args).First(&user).Error
	} else {
		err = u.db.Where("Id = ?", id).First(&user).Error
	}
	return user, err
}

func (u *UserDao) CheckEmail(email string) error {
	err := u.db.Where("email = ?", email).Take(&User{}).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil

	} else if err == nil {
		return errors.New("The email have already exists.")

	} else {
		return err
	}
}

func (u *UserDao) SelectUserInfoById(id uint) (uInfo UserInfo, err error) {
	dbQuery := u.db.Select("id", "username", "email").Where("id = ?", id)
	err = dbQuery.Model(&User{}).First(&uInfo).Error
	return uInfo, err
}

// pluck means select colomn_name(one colomn) from user(table);
func (u *UserDao) PluckNameById(id uint) (string, error) {
	var name string
	err := u.db.Model(&User{}).Where("Id = ?", id).Pluck("username", &name).Error
	return name, err
}

func (u *UserDao) SelectByEmail(email string) (User, error) {
	user := User{}
	err := u.db.Where("email = ?", email).First(&user).Error
	return user, err
}

// SELECT id, username, email
// FROM users
// WHERE username LIKE 'belikename'
// LIMIT 10 OFFSET 0;
func (u *UserDao) SelectUserInfoByCondition(condition Condition) ([]UserInfo, error) {
	dbQuery := u.db
	if condition.Id != nil {
		dbQuery = dbQuery.Where("id = ?", *condition.Id)
	}
	if condition.LikePrefixUsername != nil {
		dbQuery = dbQuery.Where("username like ?", *&condition.LikePrefixUsername)
	}
	var list []UserInfo
	err := dbQuery.Model(&User{}).Select(
		"id", "username", "email",
	).Offset(condition.Offset).Limit(condition.Limit).Find(&list).Error
	return list, err
}

func (u *UserDao) CreateFriendInvitation(inviter uint, invitee uint) (FriendInvitation, error) {
	var invitation = FriendInvitation{}
	err := u.db.Model(&invitation).Create(&invitation).Error
	return invitation, err
}

func (u *UserDao) SelectFriendInvitation(inviter uint, invitee uint, forUpdate bool) (
	result FriendInvitation, err error,
) {
	dbQuery := u.db.Where("inviter = ? AND invitee = ?", inviter, invitee)
	if forUpdate {
		dbQuery = dbQuery.Set("gorm:query_option", "FOR UPDATE")
	}
	err = dbQuery.Find(&result).Error
	return
}

func (u *UserDao) SelectFriendInvitationList(inviter *uint, invitee *uint) (
	result []FriendInvitation, err error,
) {
	query := u.db
	if inviter != nil {
		query = query.Where("inviter = ?", inviter)
	}
	// Where is a function that ADD conditions
	// you could ctrl + "pointAtIt" to see its intro
	// ---- "Where add conditions"
	if invitee != nil {
		query = query.Where("invitee = ?", invitee)
	}
	err = query.Model(&FriendInvitation{}).Find(&result).Error
	return
}

func (u *UserDao) SelectFriend(userId uint, friendId uint) (Friend, error) {
	var friend Friend
	query := u.db.Model(&Friend{}).Where("user_id = ? AND friend_id = ?", userId, friendId)
	err := query.First(&friend).Error
	return friend, err
}

// if they are friends they both should have records, so count == 2
func (u *UserDao) IsRealFriend(userId uint, friendId uint) (bool, error) {
	var count int64
	whereSql := "user_id = ? AND friend_id = ? OR friend_id = ? AND user_id = ?"
	err := u.db.Model(&Friend{}).Where(whereSql, userId, friendId, userId, friendId).Count(&count).Error
	return count == 2, err
}

func (u *UserDao) AddFriend(userId uint, friendId uint, add AddMode) (friend Friend, err error) {
	friend = Friend{UserId: userId, FriendId: friendId, AddMode: add}
	err = u.db.Create(&friend).Error
	if err != nil {
		if errors.Is(err, gorm.ErrDuplicatedKey) {
			friend, err = u.SelectFriend(userId, friendId)

		} else {
			return friend, err

		}
	}
	return friend, err
}

func (u *UserDao) SelectFriendList(userId uint) (list []Friend, err error) {
	err = u.db.Model(&Friend{}).Where("user_id = ?", userId).Find(&list).Error
	return list, err
}
