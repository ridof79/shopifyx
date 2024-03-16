package delivery

import (
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

	if err := c.Bind(&user); err != nil {
		return util.ErrorHandler(c, http.StatusBadRequest, InvalidRequestBody)
	}

	if err := c.Validate(user); err != nil {
		return util.ErrorHandler(c, http.StatusBadRequest, err.Error())
	}

	user, err := repository.RegisterUser(user.Username, user.Name, user.Password)
	if err != nil {
		if repository.IsDuplicateKeyError(err) {
			return util.ErrorHandler(c, http.StatusConflict, UsernameAreleadyExists)
		}
	}

	var userLogin domain.UserLogin
	userLogin.Id = user.Id
	userLogin.Username = user.Username
	userLogin.Name = user.Name

	token, err := auth.GenerateAccessToken(&userLogin)
	if err != nil {
		return util.ErrorHandler(c, http.StatusInternalServerError, FailedToGenerateToken)
	}

	return util.UserSuccesResponseHandler(c, http.StatusCreated, UserRegisteredSuccessfully, user.Username, user.Name, token)
}

func LoginUserHandler(c echo.Context) error {
	var user domain.UserLogin

	if err := c.Bind(&user); err != nil {
		return util.ErrorHandler(c, http.StatusBadRequest, InvalidRequestBody)
	}

	if err := c.Validate(user); err != nil {
		return util.ErrorHandler(c, http.StatusBadRequest, err.Error())
	}

	user, err := repository.LoginUser(user.Username, user.Password)

	if err != nil {
		if err == repository.ErrUsernameNotFound {
			return util.ErrorHandler(c, http.StatusNotFound, UserNotFound)
		}
		if err == repository.ErrPasswordWrong {
			return util.ErrorHandler(c, http.StatusBadRequest, UserPasswordFalse)
		}
		return util.ErrorHandler(c, http.StatusInternalServerError, InvalidRequestBody)
	}

	token, err := auth.GenerateAccessToken(&user)
	if err != nil {
		return util.ErrorHandler(c, http.StatusInternalServerError, FailedToGenerateToken)
	}

	return util.UserSuccesResponseHandler(c, http.StatusOK, UserLoggedSuccessfully, user.Username, user.Name, token)
}
