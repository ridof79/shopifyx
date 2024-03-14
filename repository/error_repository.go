package repository

import (
	"database/sql"

	"github.com/lib/pq"
)

func IsConstrainViolations(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == "23514"
	}
	return false
}

func IdNotFound(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == "22P02"
	}
	return false
}

func DontHavePermission(err error) bool {
	return err == sql.ErrNoRows
}

func IsDuplicateKeyError(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == "23505"
	}
	return false
}

func InvalidUsernameAndPasswod(err error) bool {
	if pqErr, ok := err.(*pq.Error); ok {
		return pqErr.Code == "23514"
	}
	return false
}
