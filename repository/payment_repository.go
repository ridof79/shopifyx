package repository

import (
	"database/sql"
	"shopifyx/domain"
)

func CreatePayment(tx *sql.Tx, payment *domain.Payment, productId, userId string) error {

	var paymentId string
	err := tx.QueryRow("INSERT INTO payments (bank_account_id, payment_proof_image_url, user_id) VALUES ($1, $2, $3) RETURNING id", payment.BankAccountId, payment.PaymentProofImageURL, userId).Scan(&paymentId)
	if err != nil {
		return err
	}

	_, err = tx.Exec("INSERT INTO payments_counter (product_id, quantity, payment_id) VALUES ($1, $2, $3)", productId, payment.Quantity, paymentId)
	if err != nil {
		return err
	}

	return nil
}
