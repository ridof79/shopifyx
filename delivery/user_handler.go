package delivery

import (
	"encoding/json"
	"net/http"

	"shopifyx/auth"
	"shopifyx/domain"
	"shopifyx/repository"
	"shopifyx/util"

	"github.com/labstack/echo/v4"
)

const (
	InvalidUsernameOrPasswordLength = "username or password must be 5 to 15 characters long"
	UsernameAreleadyExists          = "username already exists"
	FailedToGenerateToken           = "failed to generate token"

	UserRegisteredSuccessfully = "User registered successfully"
	UserLoggedSuccessfully     = "User logged successfully"
	UserNotFound               = "user not found"
	UserPasswordFalse          = "wrong password"
)

func RegisterUserHandler(c echo.Context) error {
	var user domain.User

	if err := json.NewDecoder(c.Request().Body).Decode(&user); err != nil {
		return util.ErrorHandler(c, http.StatusBadRequest, InvalidRequestBody)
	}

	if (len(user.Username) < 5 || len(user.Username) > 15) || (len(user.Password) < 5 || len(user.Password) > 15) {
		return util.ErrorHandler(c, http.StatusBadRequest, InvalidUsernameOrPasswordLength)
	}

	user, err := repository.RegisterUser(user.Username, user.Name, user.Password)
	if err != nil {
		if repository.IsDuplicateKeyError(err) {
			return util.ErrorHandler(c, http.StatusConflict, UsernameAreleadyExists)
		}
	}

	token, err := auth.GenerateAccessToken(&user)
	if err != nil {
		return util.ErrorHandler(c, http.StatusInternalServerError, FailedToGenerateToken)
	}

	return util.UserSuccesResponseHandler(c, http.StatusCreated, UserRegisteredSuccessfully, user.Username, user.Name, token)
}

func LoginUserHandler(c echo.Context) error {
	var user domain.User

	if err := json.NewDecoder(c.Request().Body).Decode(&user); err != nil {
		return util.ErrorHandler(c, http.StatusBadRequest, InvalidRequestBody)
	}

	if (len(user.Username) < 5 || len(user.Username) > 15) || (len(user.Password) < 5 || len(user.Password) > 15) {
		return util.ErrorHandler(c, http.StatusBadRequest, InvalidUsernameOrPasswordLength)
	}

	user, err := repository.LoginUser(user.Username, user.Password)

	if err != nil {
		if err == repository.ErrUsernameNotFound {
			return util.ErrorHandler(c, http.StatusNotFound, UserNotFound)
		}
		if err == repository.ErrUsernameNotFound {
			return util.ErrorHandler(c, http.StatusBadRequest, UserPasswordFalse)
		}
		return util.ErrorHandler(c, http.StatusInternalServerError, err.Error())
	}

	token, err := auth.GenerateAccessToken(&user)
	if err != nil {
		return util.ErrorHandler(c, http.StatusInternalServerError, FailedToGenerateToken)
	}

	return util.UserSuccesResponseHandler(c, http.StatusOK, UserLoggedSuccessfully, user.Username, user.Name, token)
}
