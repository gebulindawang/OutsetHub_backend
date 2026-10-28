package utils

import (
	"time"
	"errors"
	"github.com/golang-jwt/jwt/v5"
)

var jwtSecret = []byte("gebulinking")

type MyClaims struct {
	UserID   int
	Username string
	Role     int
	jwt.RegisteredClaims
}

func GenToken(userID int,username string,role int) (string,error){
	claims := MyClaims{
		UserID: userID,
		Username: username,
		Role: role,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(2*time.Hour)),
			IssuedAt: jwt.NewNumericDate(time.Now()),
			Issuer: "gebulinking",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256,claims)
	return token.SignedString(jwtSecret)
}

func ParseToken(tokenStr string) (*MyClaims, error) {
	claims := &MyClaims{}
	token, err := jwt.ParseWithClaims(tokenStr, claims, func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("非法的签名算法")
		}
		return jwtSecret, nil
	})
	if err != nil {
		return nil, err
	}
	if !token.Valid {
		return nil, errors.New("token 无效")
	}
	return claims, nil
}