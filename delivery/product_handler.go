package delivery

import (
	"encoding/json"
	"net/http"
	"shopifyx/auth"
	"shopifyx/domain"
	"shopifyx/repository"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func CreateProductHandler(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*auth.JwtCustomClaims)
	userId := claims.Id

	var product domain.Product

	if err := json.NewDecoder(c.Request().Body).Decode(&product); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": err.Error(),
			},
		)
	}

	err := repository.CreateProduct(&product, userId)

	if err != nil {
		if repository.IsConstrainViolations(err) {
			return c.JSON(
				http.StatusBadRequest,
				map[string]string{
					"error": "required fields are missing or invalid (e.g., price, name, category)",
				},
			)
		}

		return c.JSON(
			http.StatusInternalServerError,
			map[string]interface{}{
				"message": err.Error(),
			})
	}

	return c.JSON(
		http.StatusCreated,
		map[string]interface{}{
			"message": "Product added successfully!",
		})
}

func UpdateProductHandler(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*auth.JwtCustomClaims)
	userId := claims.Id

	productID := c.Param("productId")

	var updatedProduct domain.Product

	if err := json.NewDecoder(c.Request().Body).Decode(&updatedProduct); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": err.Error(),
			},
		)
	}

	err := repository.UpdateProduct(&updatedProduct, userId, productID)

	if err != nil {
		if repository.IsConstrainViolations(err) {
			return c.JSON(
				http.StatusBadRequest,
				map[string]string{
					"error": "required fields are missing or invalid (e.g., price, name, category)",
				},
			)
		}

		if repository.IdNotFound(err) {
			return c.JSON(
				http.StatusNotFound,
				map[string]string{
					"error": "product not found",
				},
			)
		}

		if repository.DontHavePermission(err) {
			return c.JSON(
				http.StatusForbidden,
				map[string]string{
					"error": "you don't have permission to update this product",
				},
			)
		}

		return c.JSON(
			http.StatusInternalServerError,
			map[string]interface{}{
				"message": err.Error(),
			})
	}

	return c.JSON(
		http.StatusOK,
		map[string]interface{}{
			"message": "Product updated successfully!",
		})
}

func DeleteProductHandler(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*auth.JwtCustomClaims)
	userId := claims.Id

	productID := c.Param("productId")

	err := repository.DeleteProductById(productID, userId)

	if err != nil {
		if repository.IdNotFound(err) {
			return c.JSON(
				http.StatusNotFound,
				map[string]string{
					"error": "product not found",
				},
			)
		}

		if repository.DontHavePermission(err) {
			return c.JSON(
				http.StatusForbidden,
				map[string]string{
					"error": "you don't have permission to delete this product",
				},
			)
		}

		return c.JSON(
			http.StatusInternalServerError,
			map[string]interface{}{
				"message": err.Error(),
			})
	}

	return c.JSON(
		http.StatusOK,
		map[string]interface{}{
			"message": "product deleted successfully!",
		})
}

func GetProductHandler(c echo.Context) error {
	productID := c.Param("productId")

	product, seller, err := repository.GetProductById(productID)

	if err != nil {
		if repository.IdNotFound(err) {
			return c.JSON(http.StatusNotFound,
				map[string]string{
					"error": "product not found",
				})
		}

		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{
				"error": err.Error(),
			})
	}

	response := map[string]interface{}{
		"message": "ok",
		"data": map[string]interface{}{
			"product": product,
			"seller":  seller,
		},
	}

	return c.JSON(http.StatusOK, response)
}

func UpdateProductStockHandler(c echo.Context) error {

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*auth.JwtCustomClaims)
	userId := claims.Id

	productId := c.Param("productId")
	userIdFromProductId, err := repository.GetUserIdFromProductId(productId)
	if err != nil {
		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{
				"error": err.Error(),
			})
	}

	if userIdFromProductId != userId {
		return c.JSON(
			http.StatusForbidden,
			map[string]string{
				"message": "you don't have permission to update this product",
			},
		)
	}

	var stockUpdate domain.StockUpdate

	if err := json.NewDecoder(c.Request().Body).Decode(&stockUpdate); err != nil {
		return c.JSON(
			http.StatusBadRequest,
			map[string]string{
				"error": err.Error(),
			},
		)
	}

	err = repository.UpdateProductStock(productId, stockUpdate.Stock)

	if err != nil {
		if repository.IdNotFound(err) {
			return c.JSON(
				http.StatusNotFound,
				map[string]string{
					"error": "product not found",
				},
			)
		}

		return c.JSON(
			http.StatusInternalServerError,
			map[string]string{
				"error": err.Error(),
			})
	}

	return c.JSON(
		http.StatusOK,
		map[string]string{
			"message": "stock updated successfully",
		})
}
