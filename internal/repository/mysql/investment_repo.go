package mysql

import (
	"context"
	"database/sql"
	"nobi-assesment/internal/domain"
	"nobi-assesment/internal/repository"
	"nobi-assesment/pkg/utils"
)

type mysqlInvestmentRepository struct {
	db *sql.DB
}

func NewMySQLInvestmentRepository(db *sql.DB) repository.InvestmentRepository {
	return &mysqlInvestmentRepository{db}
}

func (r *mysqlInvestmentRepository) Create(ctx context.Context, investment *domain.Investment) error {
	query := "INSERT INTO investments (id, name, total_units, total_balance) VALUES (?, ?, ?, ?)"
	_, err := r.db.ExecContext(ctx, query, investment.ID, investment.Name, investment.TotalUnits, investment.TotalBalance)
	return err
}

func (r *mysqlInvestmentRepository) GetByID(ctx context.Context, id string) (*domain.Investment, error) {
	query := "SELECT id, name, total_units, total_balance FROM investments WHERE id = ?"

	var investment domain.Investment
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&investment.ID, &investment.Name, &investment.TotalUnits, &investment.TotalBalance)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, err
		}
		return nil, err
	}

	investment.NAB = utils.ValidateNAB(investment.TotalBalance, investment.TotalUnits)

	return &investment, nil
}

func (r *mysqlInvestmentRepository) GetAll(ctx context.Context) ([]*domain.Investment, error) {
	query := "SELECT id, name, total_units, total_balance FROM investments"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	investments := []*domain.Investment{}
	for rows.Next() {
		var investment domain.Investment
		if err := rows.Scan(
			&investment.ID, &investment.Name, &investment.TotalUnits, &investment.TotalBalance); err != nil {
			return nil, err
		}

		investment.NAB = utils.ValidateNAB(investment.TotalBalance, investment.TotalUnits)
		investments = append(investments, &investment)
	}

	return investments, nil
}

func (r *mysqlInvestmentRepository) UpdateBalance(ctx context.Context, id string, amountChange float64, unitsChange float64) error {
	query := `
		UPDATE investments 
		SET total_balance = total_balance + ?, total_units = total_units + ? 
		WHERE id = ?
	`
	_, err := r.db.ExecContext(ctx, query, amountChange, unitsChange, id)
	return err
}
