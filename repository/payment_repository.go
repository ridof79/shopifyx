package repository

import (
	"database/sql"
	"shopifyx/domain"
)

func CreatePayment(tx *sql.Tx, payment *domain.Payment, productId, buyerId, sellerId string) error {

	var paymentId string

	query := `SELECT user_id FROM products WHERE id = $1`
	err := tx.QueryRow(query, productId).Scan(&sellerId)
	if err != nil {
		return err
	}

	query = `INSERT INTO payments (bank_account_id, payment_proof_image_url, buyer_id) VALUES ($1, $2, $3) RETURNING id`
	err = tx.QueryRow(
		query,
		payment.BankAccountId,
		payment.PaymentProofImageURL,
		buyerId).Scan(&paymentId)
	if err != nil {
		return err
	}

	query = `INSERT INTO payments_counter (product_id, quantity, payment_id, seller_id) VALUES ($1, $2, $3, $4)`
	_, err = tx.Exec(
		query,
		productId,
		payment.Quantity,
		paymentId,
		sellerId)
	if err != nil {
		return err
	}

	return nil
}

func ProductAndBankAccountValid(tx *sql.Tx, bankAccountId, productId string) (bool, string, error) {
	query := `
	SELECT p.is_purchaseable, p.user_id
	FROM products p
	INNER JOIN bank_accounts b ON p.user_id = b.user_id
	WHERE b.id = $1 AND p.id = $2;`

	var hasMatching bool
	var sellerId string

	err := tx.QueryRow(
		query,
		bankAccountId,
		productId).Scan(&hasMatching, &sellerId)
	if err != nil {
		return false, "", err
	}

	return hasMatching, sellerId, nil
}
