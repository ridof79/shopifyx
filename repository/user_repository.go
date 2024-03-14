package repository

import (
	"database/sql"
	"errors"
	"os"
	"shopifyx/config"
	"shopifyx/domain"
	"strconv"

	"golang.org/x/crypto/bcrypt"
)

func RegisterUser(username, name, password string) (domain.User, error) {
	var user domain.User

	hashedPassword, err := HashPassword(password)
	if err != nil {
		return user, err
	}
	err = config.GetDB().QueryRow("INSERT INTO users (username, name, password) VALUES ($1, $2, $3) RETURNING id, name, username", username, name, hashedPassword).Scan(&user.Id, &user.Name, &user.Username)
	if err != nil {
		return user, err
	}
	return user, nil
}

func LoginUser(username, password string) (domain.User, error) {
	var storedPassword string
	var user domain.User
	err := config.GetDB().QueryRow("SELECT id, username, name, password FROM users WHERE username = $1", username).Scan(&user.Id, &user.Username, &user.Name, &storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, errors.New("user not found")
		}
		return user, err
	}

	err = VerifyPassword(storedPassword, password)
	if err != nil {
		return user, errors.New("wrong password")
	}

	return user, nil
}

func HashPassword(password string) (string, error) {
	cost, _ := strconv.Atoi(os.Getenv("BCRYPT_SALT"))

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), cost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

func VerifyPassword(hashedPassword, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
}
