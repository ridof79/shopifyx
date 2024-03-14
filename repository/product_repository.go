package repository

import (
	"database/sql"

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
		return err
	}
	return nil
}

func GetProductById(productId string) (domain.Product, domain.SellerResponse, error) {
	var product domain.Product
	var seller domain.SellerResponse
	var totalSold int

	sellerId, _ := GetUserIdFromProductId(productId)
	_ = config.GetDB().QueryRow("SELECT COALESCE(SUM(pc.quantity), 0) AS product_purchase_count FROM payments_counter pc WHERE pc.seller_id = $1", sellerId).Scan(&totalSold)

	var arrBankAccountId []sql.NullString
	var arrBankNames []sql.NullString
	var arrBankAccountNames []sql.NullString
	var arrBankAccountNumbers []sql.NullString

	query := `
	SELECT 
		p.name AS product_name,
		p.price AS product_price,
		p.image_url AS product_image_url,
		p.stock AS product_stock,
		p.condition AS product_condition,
		p.tags AS product_tags,
		p.is_purchaseable AS product_purchaseable,
		COALESCE(SUM(pc.quantity), 0) AS product_purchase_count, 
		u.name AS seller_name,
		(
			SELECT ARRAY_AGG(ba.id) 
			FROM bank_accounts ba 
			WHERE ba.user_id = u.id
		) AS bank_account_id,
		(
			SELECT ARRAY_AGG(ba.bank_name) 
			FROM bank_accounts ba 
			WHERE ba.user_id = u.id
		) AS bank_names,
		(
			SELECT ARRAY_AGG(ba.bank_account_name) 
			FROM bank_accounts ba 
			WHERE ba.user_id = u.id
		) AS bank_account_names,
		(
			SELECT ARRAY_AGG(ba.bank_account_number) 
			FROM bank_accounts ba 
			WHERE ba.user_id = u.id
		) AS bank_account_numbers
	FROM 
		products p
	LEFT JOIN 
		users u ON p.user_id = u.id
	LEFT JOIN 
		payments_counter pc ON p.id = pc.product_id
	WHERE 
		p.id = $1
	GROUP BY 
		p.id, p.name, u.name, u.id;`

	rows, err := config.GetDB().Query(query, productId)
	if err != nil {
		return domain.Product{}, domain.SellerResponse{}, err
	}
	defer rows.Close()

	for rows.Next() {
		err := rows.Scan(
			&product.Name,
			&product.Price,
			&product.ImageURL,
			&product.Stock,
			&product.Condition,
			pq.Array(&product.Tags),
			&product.IsPurchaseable,
			&product.PurchaseCount,
			&seller.Name,
			pq.Array(&arrBankAccountId),
			pq.Array(&arrBankNames),
			pq.Array(&arrBankAccountNames),
			pq.Array(&arrBankAccountNumbers),
		)
		if err != nil {
			return domain.Product{}, domain.SellerResponse{}, err
		}
	}

	var bankAccounts []domain.BankAccounts

	for i := range arrBankAccountId {
		bankAccounts = append(bankAccounts, domain.BankAccounts{
			Id:                arrBankAccountId[i].String,
			BankName:          arrBankNames[i].String,
			BankAccountName:   arrBankAccountNames[i].String,
			BankAccountNumber: arrBankAccountNumbers[i].String,
		})
	}

	seller.BankAccounts = bankAccounts
	seller.ProductSoldTotal = totalSold

	return product, seller, nil
}

func UpdateProduct(product *domain.Product, productId, userId string) (int, error) {
	query := `
		WITH updated AS (
			UPDATE products
			SET name = $1, price = $2, image_url = $3, condition = $4, tags = $5, is_purchaseable = $6
			WHERE id = $7 AND user_id = $8
			RETURNING *
		)
		SELECT 
			CASE 
				WHEN EXISTS (SELECT 1 FROM updated) THEN 1 
				WHEN NOT EXISTS (SELECT 1 FROM products WHERE id = $7) THEN 2 
				ELSE 3 
			END AS result_code;
	`

	var resultCode int
	err := config.GetDB().QueryRow(query,
		product.Name, product.Price, product.ImageURL, product.Condition, pq.Array(product.Tags), product.IsPurchaseable, productId, userId,
	).Scan(&resultCode)

	if err != nil {
		return 0, err
	}
	return resultCode, err
}

func DeleteProductById(productId, userId string) (int, error) {
	query :=
		`WITH deleted AS (
		DELETE FROM products 
		WHERE id = $1 AND user_id = $2
		RETURNING *
	)
	SELECT 
		CASE 
			WHEN EXISTS (SELECT 1 FROM deleted) THEN 1 
			WHEN NOT EXISTS (SELECT 1 FROM products WHERE id = $1) THEN 2 
			ELSE 3 
		END AS result_code;`

	var resultCode int
	err := config.GetDB().QueryRow(query, productId, userId).Scan(&resultCode)
	if err != nil {
		return 0, err
	}

	return resultCode, nil
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
