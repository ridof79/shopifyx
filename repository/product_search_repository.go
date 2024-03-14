package repository

import (
	"errors"
	"fmt"
	"shopifyx/config"
	"shopifyx/domain"
	"shopifyx/util"

	"github.com/lib/pq"
)

func SearchProduct(searchPagination *util.SearchPagination, userId string) ([]domain.ProductResponse, int, error) {

	query := `
		SELECT p.id, p.name, p.price, p.image_url, p.stock, p.condition, p.tags, p.is_purchaseable, p.created_at as date,
		COALESCE((SELECT SUM(pc.quantity) 
		FROM payments_counter pc WHERE pc.product_id = p.id), 0) AS total_sold
		FROM products p
		WHERE 1 = 1
	`
	// Buat slice untuk menyimpan nilai parameter prepared statement
	var args []interface{}

	paramIndex := 1
	// Tambahkan filter berdasarkan userOnly
	if searchPagination.UserOnly {
		query += fmt.Sprintf(" AND user_id = $%d", paramIndex)
		paramIndex++
		args = append(args, userId)
	}

	// Tambahkan filter berdasarkan condition
	if searchPagination.Condition != "" {
		if searchPagination.Condition != domain.ConditionEnum("new") {
			query += fmt.Sprintf(" AND condition = $%d", paramIndex)
			args = append(args, searchPagination.Condition)
			paramIndex++
		}
		if searchPagination.Condition != domain.ConditionEnum("second") {
			query += fmt.Sprintf(" AND condition = $%d", paramIndex)
			args = append(args, searchPagination.Condition)
			paramIndex++
		}
	}

	// Tambahkan filter untuk menampilkan produk dengan stok kosong jika showEmptyStock=true
	if !searchPagination.ShowEmptyStock {
		query += " AND stock > 0"
	}

	// Tambahkan filter berdasarkan rentang harga
	if searchPagination.MaxPrice != 0 {
		query += fmt.Sprintf(" AND price <= $%d", paramIndex)
		args = append(args, searchPagination.MaxPrice)
		paramIndex++
	}
	if searchPagination.MinPrice != 0 {
		query += fmt.Sprintf(" AND price >= $%d", paramIndex)
		args = append(args, searchPagination.MinPrice)
		paramIndex++
	}

	// Tambahkan filter berdasarkan pencarian nama produk
	if searchPagination.Search != "" {
		query += fmt.Sprintf(" AND name LIKE $%d", paramIndex)
		args = append(args, "%"+searchPagination.Search+"%")
		paramIndex++
	}

	// Hitung jumlah total produk tanpa paging
	totalQuery := "SELECT COUNT(*) FROM (" + query + ") AS total"
	var total int
	err := config.GetDB().QueryRow(totalQuery, args...).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	// Tambahkan fitur pagination ke query
	query += fmt.Sprintf(" ORDER BY %s %s LIMIT $%d OFFSET $%d", searchPagination.SortBy, searchPagination.OrdedBy, paramIndex, paramIndex+1)
	args = append(args, searchPagination.Limit, searchPagination.Offset)

	// Eksekusi query
	rows, err := config.GetDB().Query(query, args...)
	if err != nil {
		return nil, 0, err
	}
	defer rows.Close()

	var products []domain.ProductResponse
	for rows.Next() {
		var product domain.ProductResponse
		var date string

		err := rows.Scan(&product.Id, &product.Name, &product.Price, &product.ImageURL, &product.Stock, &product.Condition, pq.Array(&product.Tags),
			&product.IsPurchaseable, &date, &product.PurchaseCount)
		if err != nil {
			return nil, 0, err
		}

		products = append(products, product)
	}

	if len(products) == 0 {
		return nil, 0, errors.New(query)
	}

	return products, total, nil
}
