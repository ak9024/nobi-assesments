package mysql

import (
	"context"
	"database/sql"
	"nobi-assesment/internal/domain"
	"nobi-assesment/internal/repository"
)

type mysqlTransactionRepository struct {
	db *sql.DB
}

func NewMySQLTransactionRepository(db *sql.DB) repository.TransactionRepository {
	return &mysqlTransactionRepository{db}
}

func (r *mysqlTransactionRepository) Create(ctx context.Context, transaction *domain.Transaction) error {
	query := `
		INSERT INTO transactions (id, customer_id, investment_id, type, amount, units, nab) 
		VALUES (?, ?, ?, ?, ?, ?, ?)
	`
	_, err := r.db.ExecContext(ctx, query,
		transaction.ID,
		transaction.CustomerID,
		transaction.InvestmentID,
		transaction.Type,
		transaction.Amount,
		transaction.Units,
		transaction.NAB)
	return err
}

func (r *mysqlTransactionRepository) GetByCustomerID(ctx context.Context, customerID string) ([]*domain.Transaction, error) {
	query := `
		SELECT t.id, t.customer_id, t.investment_id, i.name, t.type, t.amount, t.units, t.nab, t.transaction_date
		FROM transactions t
		JOIN investments i ON t.investment_id = i.id
		WHERE t.customer_id = ?
		ORDER BY t.transaction_date DESC
	`
	rows, err := r.db.QueryContext(ctx, query, customerID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	transactions := []*domain.Transaction{}
	for rows.Next() {
		var transaction domain.Transaction
		var investmentName string

		if err := rows.Scan(
			&transaction.ID,
			&transaction.CustomerID,
			&transaction.InvestmentID,
			&investmentName,
			&transaction.Type,
			&transaction.Amount,
			&transaction.Units,
			&transaction.NAB,
			&transaction.TransactionDate); err != nil {
			return nil, err
		}

		transactions = append(transactions, &transaction)
	}

	return transactions, nil
}
