package commonService

import (
	"goAccounting/global"
	utilJwt "goAccounting/util/jwtTool"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type common struct{}

var Common = new(common)

func (cm *common) HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}

	return string(hashedPassword), nil
}

func CheckInt(v any) (i int) {
	switch v := v.(type) {
	case int:
		i = v
	default:
		i = 0
	}
	return
}

const ExpireAt time.Duration = 90 * 24 * time.Hour

func (cm *common) MakeCustomClaims(userId uint) jwt.RegisteredClaims {
	expirationTime := time.Now().Add(ExpireAt)
	return jwt.RegisteredClaims{
		ID:        strconv.Itoa(int(userId)),
		ExpiresAt: jwt.NewNumericDate(expirationTime),
		Issuer:    "server",
		Subject:   "user",
	}
}

func (cm *common) parseToken(tokenStr string) (jwt.RegisteredClaims, error) {
	return utilJwt.ParseToken(tokenStr, []byte(global.Config.System.JwtKey))
}

func (cm *common) GenerateJWT(claims jwt.RegisteredClaims) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(global.Config.System.JwtKey))
	if err != nil {
		return "", err
	}
	return tokenString, nil
}
