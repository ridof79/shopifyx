package repository

import (
	"database/sql"
	"errors"

	"shopifyx/config"
	"shopifyx/domain"

	"github.com/lib/pq"
)

func CreateProduct(product *domain.Product, userId string) error {
	_, err := config.GetDB().Exec(
		`INSERT INTO products (name, price, image_url, stock, condition, tags, is_purchaseable, user_id) VALUES($1, $2, $3, $4, $5, $6, $7, $8)`,
		product.Name, product.Price, product.ImageURL, product.Stock, product.Condition, pq.Array(product.Tags), product.IsPurchaseable, userId,
	)
	if err != nil {
		return errors.New("failed to create product")
	}
	return nil
}

func GetProductById(productId string) (domain.ProductResponse, domain.SellerResponse, error) {
	var product domain.ProductResponse
	var seller domain.SellerResponse

	var bankAccountId []string
	var bankNames []string
	var bankAccountNames []string
	var bankAccountNumbers []string

	query := `
		SELECT 
		p.id AS product_id,
		p.name AS product_name,
		p.price AS product_price,
		p.image_url AS product_image_url,
		p.stock AS product_stock,
		p.condition AS product_condition,
		p.tags AS product_tags,
		p.is_purchaseable AS product_purchaseable,
		COALESCE(SUM(pc.quantity), 0) AS purchase_count,
		u.name AS seller_name,
		ARRAY_AGG(ba.id) AS bank_account_id,
		ARRAY_AGG(ba.bank_name) AS bank_names,
		ARRAY_AGG(ba.bank_account_name) AS bank_account_names,
		ARRAY_AGG(ba.bank_account_number) AS bank_account_numbers
	FROM 
		products p
	JOIN 
		users u ON p.user_id = u.id
	LEFT JOIN 
		bank_accounts ba ON u.id = ba.user_id
	LEFT JOIN 
		payments_counter pc ON p.id = pc.product_id
	WHERE 
		p.id = $1
	GROUP BY 
		p.id, p.name, u.name;
	`
	rows, err := config.GetDB().Query(query, productId)
	if err != nil {
		return domain.ProductResponse{}, domain.SellerResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&product.Id, &product.Name, &product.Price, &product.ImageURL, &product.Stock, &product.Condition, pq.Array(&product.Tags), &product.IsPurchaseable,
			&product.PurchaseCount, &seller.Name,
			pq.Array(&bankAccountId), pq.Array(&bankNames), pq.Array(&bankAccountNames), pq.Array(&bankAccountNumbers),
		)
		if err != nil {
			return domain.ProductResponse{}, domain.SellerResponse{}, err
		}
	}

	var bankAccounts []domain.BankAccounts
	for i := range bankAccountId {
		bankAccounts = append(bankAccounts, domain.BankAccounts{
			Id:                bankAccountId[i],
			BankName:          bankNames[i],
			BankAccountName:   bankAccountNames[i],
			BankAccountNumber: bankAccountNumbers[i],
		})
	}
	seller.BankAccounts = bankAccounts

	return product, seller, nil
}

func UpdateProduct(product *domain.Product, userId, productId string) error {
	query := `
		UPDATE products
		SET name = $1, price = $2, image_url = $3, condition = $4, tags = $5, is_purchaseable = $6
		WHERE id = $7 AND user_id = $8
		RETURNING id
	`

	var updatedProductID string
	err := config.GetDB().QueryRow(query,
		product.Name, product.Price, product.ImageURL, product.Condition, pq.Array(product.Tags), product.IsPurchaseable, productId, userId,
	).Scan(&updatedProductID)

	if err != nil {
		return err
	}
	return nil
}

func DeleteProductById(productId, userId string) error {
	result, err := config.GetDB().Exec(
		`DELETE FROM products WHERE id = $1 AND user_id = $2`,
		productId, userId,
	)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil {
		return err
	}
	if rowsAffected == 0 {
		err := &pq.Error{Code: "22P02"}
		return err
	}

	return nil
}

func GetProductStockTx(tx *sql.Tx, productId string) (int, error) {
	var stock int
	err := tx.QueryRow("SELECT stock FROM products WHERE id = $1", productId).Scan(&stock)
	if err != nil {
		return 0, err
	}
	return stock, nil
}

func UpdateProductStockTx(tx *sql.Tx, productId string, newStock int) error {
	_, err := tx.Exec(
		`UPDATE products SET stock = $1 WHERE id = $2`,
		newStock, productId,
	)
	if err != nil {
		return err
	}
	return nil
}

func GetUserIdFromProductId(productId string) (string, error) {
	var userId string
	err := config.GetDB().QueryRow("SELECT user_id FROM products WHERE id = $1", productId).Scan(&userId)
	if err != nil {
		return "", err
	}
	return userId, nil
}
func UpdateProductStock(productId string, newStock int) error {
	_, err := config.GetDB().Exec(
		`UPDATE products SET stock = $1 WHERE id = $2`,
		newStock, productId,
	)
	if err != nil {
		return err
	}
	return nil
}
