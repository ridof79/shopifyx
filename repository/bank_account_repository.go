package repository

import (
	"shopifyx/config"
	"shopifyx/domain"

	"github.com/lib/pq"
)

func AddBankAccount(bankAccount *domain.BankAccount, userId string) error {
	query := `INSERT INTO bank_accounts (bank_name, bank_account_name, bank_account_number, user_id) VALUES($1, $2, $3, $4)`
	_, err := config.GetDB().Exec(
		query,
		bankAccount.BankName,
		bankAccount.BankAccountName,
		bankAccount.BankAccountNumber,
		userId,
	)
	if err != nil {
		return err
	}
	return err
}

func GetBankAccounts(userId string) ([]domain.BankAccount, error) {
	query := `SELECT id, bank_name, bank_account_name, bank_account_number FROM bank_accounts WHERE user_id = $1`
	rows, err := config.GetDB().Query(
		query,
		userId,
	)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var bankAccounts []domain.BankAccount
	for rows.Next() {
		var bankAccount domain.BankAccount
		err := rows.Scan(
			&bankAccount.Id,
			&bankAccount.BankName,
			&bankAccount.BankAccountName,
			&bankAccount.BankAccountNumber)
		if err != nil {
			return nil, err
		}
		bankAccounts = append(bankAccounts, bankAccount)
	}
	return bankAccounts, nil
}

func UpdateBankAccount(bankAccount *domain.BankAccount, bankAccountId, userId string) (int, error) {
	query := `
	WITH updated AS (
		UPDATE bank_accounts
		SET bank_name = $1, bank_account_name = $2, bank_account_number = $3
		WHERE id = $4 AND user_id = $5
		RETURNING *
	)
	SELECT 
		CASE 
			WHEN EXISTS (SELECT 1 FROM updated) THEN 1 
			WHEN NOT EXISTS (SELECT 1 FROM bank_accounts WHERE id = $4) THEN 2 
			ELSE 3 
		END AS result_code;`

	var resultCode int
	err := config.GetDB().QueryRow(
		query,
		bankAccount.BankName,
		bankAccount.BankAccountName,
		bankAccount.BankAccountNumber,
		bankAccountId,
		userId,
	).Scan(&resultCode)

	if err != nil {
		return 0, err
	}
	return resultCode, err
}

func DeleteBankAccount(bankAccountId, userId string) error {
	query := `DELETE FROM bank_accounts WHERE id = $1 AND user_id = $2`
	result, err := config.GetDB().Exec(
		query,
		bankAccountId, userId,
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
