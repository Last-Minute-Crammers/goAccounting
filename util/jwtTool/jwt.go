package jwtTool

import (
	"errors"
	"fmt"
	"strconv"

	"github.com/golang-jwt/jwt/v5"
)

var (
	TokenExpired     error  = errors.New("Token is expired")
	TokenNotValidYet error  = errors.New("Token not active yet")
	TokenMalformed   error  = errors.New("That's not even a token")
	TokenInvalid     error  = errors.New("Couldn't handle this token:")
	SignKey          string = "test"
)

func CreateToken(claims jwt.RegisteredClaims, key []byte) (string, error) {
	return jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(key)
}

func ParseToken(tokenStr string, key []byte) (claims jwt.RegisteredClaims, err error) {
	token, err := jwt.ParseWithClaims(
		tokenStr, &claims, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return key, nil
		},
	)
	if err != nil {
		return
	}
	if !token.Valid {
		err = errors.New("parse token fail")
	}
	return
}

func ParseUserIdFromToken(tokenStr string, key []byte) (uint, error) {
	claims, err := ParseToken(tokenStr, key)
	if err != nil {
		return 0, err
	}
	// userId 存在于 claims.Subject
	// 你登录时生成 token 时，claims.Subject 应该写成用户id的字符串
	if claims.Subject == "" {
		return 0, errors.New("userId(subject) not found in token")
	}
	id, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return 0, errors.New("userId(subject) not uint")
	}
	return uint(id), nil
}
