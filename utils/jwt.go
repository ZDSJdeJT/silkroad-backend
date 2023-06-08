package utils

import (
	"github.com/golang-jwt/jwt"
	"os"
	"time"
)

func GenerateJWT(issuer string, expire time.Duration) (string, error) {
	claims := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Issuer:    issuer,
		ExpiresAt: time.Now().Add(expire).Unix(),
	})

	token, err := claims.SignedString([]byte(os.Getenv(JWTSecretKey)))
	if err != nil {
		return "", err
	}
	return token, nil
}
