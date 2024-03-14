package util

import "shopifyx/domain"

type SortEnum string

const (
	Price SortEnum = "price"
	Date  SortEnum = "date"
)

type OrderEnum string

const (
	Ascending  OrderEnum = "ASC"
	Descending OrderEnum = "DESC"
)

type SearchPagination struct {
	UserOnly       bool                 `json:"userOnly"`
	Limit          int                  `json:"limit"`
	Offset         int                  `json:"offset"`
	Tags           []string             `json:"tags"`
	Condition      domain.ConditionEnum `json:"condition"`
	ShowEmptyStock bool                 `json:"showEmptyStock"`
	MaxPrice       int                  `json:"maxPrice"`
	MinPrice       int                  `json:"minPrice"`
	SortBy         SortEnum             `json:"sortBy"`
	OrdedBy        OrderEnum            `json:"orderBy"`
	Search         string               `json:"search"`
}
