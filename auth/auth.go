package auth

import (
	"os"
	"shopifyx/domain"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type JwtCustomClaims struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(user *domain.User) (string, error) {

	var jwtSecretKey = os.Getenv("JWT_SECRET")
	var jwtExpiredMinutes = os.Getenv("JWT_EXPIRED_MINUTES")

	var tokenExpirationTime, err = strconv.Atoi(jwtExpiredMinutes)
	if err != nil {
		panic(err)
	}

	claims := &JwtCustomClaims{
		Id:   user.Id,
		Name: user.Name,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(tokenExpirationTime) * time.Minute)),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(jwtSecretKey))
	if err != nil {
		return t, err
	}

	return t, nil
}
