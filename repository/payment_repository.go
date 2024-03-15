package repository

import (
	"database/sql"
	"shopifyx/domain"
)

func CreatePayment(tx *sql.Tx, payment *domain.Payment, productId, buyerId, sellerId string) error {

	query := `
	INSERT INTO payments 
	(bank_account_id, payment_proof_image_url, buyer_id, product_id, quantity) 
	VALUES ($1, $2, $3, $4, $5)`

	_, err := tx.Exec(
		query,
		payment.BankAccountId,
		payment.PaymentProofImageURL,
		buyerId,
		productId,
		payment.Quantity,
	)
	if err != nil {
		return err
	}

	return nil
}

func CheckStockProductAndBankAccountValid(tx *sql.Tx, bankAccountId, productId string) (bool, int, string, error) {
	query := `
	SELECT is_purchaseable, stock, seller_id 
	FROM seller_bank_account 
	WHERE bank_account_id = $1
	AND product_id = $2;`

	var isPurchaseable bool
	var stock int
	var sellerId string

	err := tx.QueryRow(
		query,
		bankAccountId,
		productId).Scan(&isPurchaseable, &stock, &sellerId)
	if err != nil {
		return false, 0, "", err
	}

	return isPurchaseable, stock, sellerId, err
}
