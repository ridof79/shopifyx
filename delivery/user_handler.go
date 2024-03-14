package delivery

import (
	"encoding/json"
	"net/http"

	"shopifyx/auth"
	"shopifyx/domain"
	"shopifyx/repository"

	"github.com/labstack/echo/v4"
)

func RegisterUserHandler(c echo.Context) error {
	var user domain.User

	if err := json.NewDecoder(c.Request().Body).Decode(&user); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": err.Error(),
			},
		)
	}

	if (len(user.Username) < 5 || len(user.Username) > 15) || (len(user.Password) < 5 || len(user.Password) > 15) {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": "username or password must be 5 to 15 characters long",
			},
		)
	}

	user, err := repository.RegisterUser(user.Username, user.Name, user.Password)

	if err != nil {
		if repository.IsDuplicateKeyError(err) {
			return c.JSON(
				http.StatusConflict,
				map[string]string{
					"error": "username already exists",
				},
			)
		}
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{
				"error": err.Error(),
			},
		)
	}

	token, err := auth.GenerateAccessToken(&user)

	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{
				"error": "failed to generate token",
			},
		)
	}

	return c.JSON(
		http.StatusCreated,
		map[string]interface{}{
			"message": "User registered successfully",
			"data": map[string]interface{}{
				"username":    user.Username,
				"name":        user.Name,
				"accessToken": token,
			},
		})
}

func LoginUserHandler(c echo.Context) error {
	var user domain.User

	if err := json.NewDecoder(c.Request().Body).Decode(&user); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": "invalid request body",
			},
		)
	}

	if (len(user.Username) < 5 || len(user.Username) > 15) || (len(user.Password) < 5 || len(user.Password) > 15) {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": "username or password must be 5 to 15 characters long",
			},
		)
	}

	user, err := repository.LoginUser(user.Username, user.Password)

	if err != nil {
		if err.Error() == "user not found" {
			return c.JSON(
				http.StatusNotFound,
				map[string]string{
					"error": "user not found",
				},
			)
		}
		if err.Error() == "wrong password" {
			return c.JSON(
				http.StatusBadRequest,
				map[string]string{
					"error": "wrong password",
				},
			)
		}
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{
				"error": err.Error(),
			},
		)
	}

	token, err := auth.GenerateAccessToken(&user)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{
				"error": "failed to generate token",
			},
		)
	}

	return c.JSON(http.StatusOK, map[string]interface{}{
		"message": "User logged successfully",
		"data": map[string]interface{}{
			"username":    user.Username,
			"name":        user.Name,
			"accessToken": token,
		},
	})
}
