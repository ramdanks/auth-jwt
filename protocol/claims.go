package protocol

import (
	"time"

	"github.com/golang-jwt/jwt"
)

type Claims struct {
	jwt.StandardClaims
	Username string `json:"username"`
}

func NewClaims(expirationTime time.Duration, username string) *Claims {
	return &Claims{
		StandardClaims: jwt.StandardClaims{
			Issuer:    "auth-jwt",
			ExpiresAt: time.Now().Add(expirationTime).Unix(),
		},
		Username: username,
	}
}
