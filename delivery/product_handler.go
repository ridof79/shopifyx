package delivery

import (
	"encoding/json"
	"net/http"
	"shopifyx/auth"
	"shopifyx/domain"
	"shopifyx/repository"
	"shopifyx/util"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

const (
	FailedToCreateProduct = "failed to create product"
	FailedToUpdateProduct = "failed to update product"
	FailedToDeleteProduct = "failed to delete product"
	FailedToFetchProduct  = "failed to fetch product"
	FailedToUpdateStock   = "failed to update stock"

	ProductAddedSuccessfully   = "product added successfully!"
	ProductUpdatedSuccessfully = "product updated successfully!"
	ProductNotFound            = "product not found"

	ProductDeletedSuccessfully = "product deleted successfully!"
)

func CreateProductHandler(c echo.Context) error {

	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*auth.JwtCustomClaims)
	userId := claims.Id

	var product domain.Product

	if err := json.NewDecoder(c.Request().Body).Decode(&product); err != nil {
		return util.ErrorHandler(c, http.StatusBadRequest, InvalidRequestBody)
	}

	err := repository.CreateProduct(&product, userId)

	if err != nil {
		if repository.IsConstrainViolations(err) {
			return util.ErrorHandler(c, http.StatusBadRequest, RequredFieldsMissing)
		}
		return util.ErrorHandler(c, http.StatusInternalServerError, FailedToCreateProduct)
	}

	return util.ResponseHandler(c, http.StatusCreated, ProductAddedSuccessfully)
}

func UpdateProductHandler(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*auth.JwtCustomClaims)
	userId := claims.Id

	productID := c.Param("productId")

	var updatedProduct domain.Product

	if err := json.NewDecoder(c.Request().Body).Decode(&updatedProduct); err != nil {
		return util.ErrorHandler(c, http.StatusBadRequest, InvalidRequestBody)
	}

	result, err := repository.UpdateProduct(&updatedProduct, productID, userId)

	switch result {
	case 1:
		return util.ResponseHandler(c, http.StatusOK, ProductUpdatedSuccessfully)
	case 2:
		return util.ErrorHandler(c, http.StatusNotFound, ProductNotFound)
	case 3:
		return util.ErrorHandler(c, http.StatusForbidden, DontHavePermission)
	}

	if err != nil {
		if repository.IsConstrainViolations(err) {
			return util.ErrorHandler(c, http.StatusBadRequest, RequredFieldsMissing)
		}

		if repository.IdNotFound(err) {
			return util.ErrorHandler(c, http.StatusNotFound, ProductNotFound)
		}

		return util.ErrorHandler(c, http.StatusInternalServerError, FailedToUpdateProduct)
	}
	return nil
}

func DeleteProductHandler(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*auth.JwtCustomClaims)
	userId := claims.Id

	productID := c.Param("productId")

	result, err := repository.DeleteProductById(productID, userId)

	switch result {
	case 1:
		return util.ResponseHandler(c, http.StatusOK, ProductDeletedSuccessfully)
	case 2:
		return util.ErrorHandler(c, http.StatusNotFound, ProductNotFound)
	case 3:
		return util.ErrorHandler(c, http.StatusForbidden, DontHavePermission)
	}

	if err != nil {
		if repository.IdNotFound(err) {
			return util.ErrorHandler(c, http.StatusNotFound, ProductNotFound)
		}

		return util.ErrorHandler(c, http.StatusInternalServerError, FailedToDeleteProduct)
	}
	return nil
}

func GetProductHandler(c echo.Context) error {
	productID := c.Param("productId")

	product, seller, err := repository.GetProductById(productID)

	if err != nil {
		if repository.IdNotFound(err) {
			return util.ErrorHandler(c, http.StatusNotFound, ProductNotFound)
		}

		return util.ErrorHandler(c, http.StatusInternalServerError, FailedToFetchProduct)
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
		return util.ErrorHandler(c, http.StatusInternalServerError, FailedToFetchProduct)
	}

	if userIdFromProductId != userId {
		return util.ErrorHandler(c, http.StatusForbidden, DontHavePermission)
	}

	var stockUpdate domain.StockUpdate

	if err := json.NewDecoder(c.Request().Body).Decode(&stockUpdate); err != nil {
		return util.ErrorHandler(c, http.StatusBadRequest, InvalidRequestBody)
	}

	err = repository.UpdateProductStock(productId, stockUpdate.Stock)

	if err != nil {
		if repository.IdNotFound(err) {
			return util.ErrorHandler(c, http.StatusNotFound, ProductNotFound)
		}
		return util.ErrorHandler(c, http.StatusInternalServerError, FailedToUpdateStock)
	}

	return util.ResponseHandler(c, http.StatusOK, "stock updated successfully")
}
