package delivery

import (
	"net/http"
	"strconv"

	"shopifyx/auth"
	"shopifyx/domain"
	"shopifyx/repository"
	"shopifyx/util"

	"github.com/golang-jwt/jwt/v5"
	"github.com/labstack/echo/v4"
)

func SearchProductHandler(c echo.Context) error {
	user := c.Get("user").(*jwt.Token)
	claims := user.Claims.(*auth.JwtCustomClaims)
	userId := claims.Id

	userOnly, _ := strconv.ParseBool(c.QueryParam("userOnly"))
	limit, _ := strconv.Atoi(c.QueryParam("limit"))
	offset, _ := strconv.Atoi(c.QueryParam("offset"))
	tags := []string{c.QueryParam("tags")}
	condition := domain.ConditionEnum(c.QueryParam("condition"))
	showEmptyStock, _ := strconv.ParseBool(c.QueryParam("showEmptyStock"))
	maxPrice, _ := strconv.Atoi(c.QueryParam("maxPrice"))
	minPrice, _ := strconv.Atoi(c.QueryParam("minPrice"))
	sortBy := util.SortEnum(c.QueryParam("sortBy"))
	orderBy := util.OrderEnum(c.QueryParam("orderBy"))

	if limit == 0 {
		limit = 10
	}
	if sortBy == "" {
		sortBy = util.Price
	}

	searchPagination := &util.SearchPagination{
		UserOnly:       userOnly,
		Limit:          limit,
		Offset:         offset,
		Tags:           tags,
		Condition:      condition,
		ShowEmptyStock: showEmptyStock,
		MaxPrice:       maxPrice,
		MinPrice:       minPrice,
		SortBy:         sortBy,
		OrdedBy:        orderBy,
		Search:         c.QueryParam("search"),
	}

	products, total, err := repository.SearchProduct(searchPagination, userId)
	if err != nil {
		return util.ErrorHandler(c, http.StatusInternalServerError, FailedToFetchProduct)
	}

	return util.SerachProductPaginationResponseHandler(c, http.StatusOK, products, limit, offset, total)
}
