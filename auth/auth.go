package auth

import (
	"os"
	"shopifyx/domain"
	"strconv"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var jwtSecretKey string

func init() {
	jwtSecretKey = os.Getenv("JWT_SECRET_KEY")
}

func GetJWTSecret() string {
	return jwtSecretKey
}

type JwtCustomClaims struct {
	Id   string `json:"id"`
	Name string `json:"name"`
	jwt.RegisteredClaims
}

func GenerateAccessToken(user *domain.User) (string, error) {

	var tokenExpirationTime, err = strconv.Atoi(os.Getenv("JWT_EXPIRED_MINUTES"))
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

	t, err := token.SignedString([]byte("secret"))
	if err != nil {
		return t, err
	}

	return t, nil
}
