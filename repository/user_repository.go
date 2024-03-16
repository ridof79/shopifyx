package repository

import (
	"database/sql"
	"shopifyx/auth"
	"shopifyx/db"
	"shopifyx/domain"
)

func RegisterUser(username, name, password string) (domain.User, error) {
	var user domain.User

	hashedPassword, err := auth.HashPassword(password)
	if err != nil {
		return user, err
	}

	query := `INSERT INTO users (username, name, password) VALUES ($1, $2, $3) 
			  RETURNING id, name, username`
	err = db.GetDB().QueryRow(
		query,
		username,
		name,
		hashedPassword).Scan(&user.Id, &user.Name, &user.Username)
	if err != nil {
		return user, err
	}
	return user, nil
}

func LoginUser(username, password string) (domain.UserLogin, error) {
	var storedPassword string
	var user domain.UserLogin

	query := `SELECT id, name, username, password FROM users WHERE username = $1`
	err := db.GetDB().QueryRow(query,
		username).Scan(
		&user.Id,
		&user.Name,
		&user.Username,
		&storedPassword)
	if err != nil {
		if err == sql.ErrNoRows {
			return user, ErrUsernameNotFound
		}
		return user, err
	}

	err = auth.VerifyPassword(storedPassword, password)
	if err != nil {
		return user, ErrPasswordWrong
	}

	return user, nil
}
