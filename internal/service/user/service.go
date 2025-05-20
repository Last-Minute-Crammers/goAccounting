package userService

import (
	userModel "goAccounting/internal/model/user"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/net/context"
)

type User struct{}

func (userSvc *User) Login(email string, password string, ctx context.Context) (
	user userModel.User, token string, customClaims jwt.RegisteredClaims, err error,
) {
	password = commonService
}
