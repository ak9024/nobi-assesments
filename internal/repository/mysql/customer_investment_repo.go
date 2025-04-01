package mysql

import (
	"context"
	"database/sql"
	"errors"
	"nobi-assesment/internal/domain"
	"nobi-assesment/internal/repository"
)

type mysqlCustomerInvestmentRepository struct {
	db *sql.DB
}

func NewMySQLCustomerInvestmentRepository(db *sql.DB) repository.CustomerInvestmentRepository {
	return &mysqlCustomerInvestmentRepository{db}
}

func (r *mysqlCustomerInvestmentRepository) Create(ctx context.Context, customerInvestment *domain.CustomerInvestment) error {
	query := "INSERT INTO customer_investments (id, customer_id, investment_id, units) VALUES (?, ?, ?, ?)"
	_, err := r.db.ExecContext(ctx, query,
		customerInvestment.ID,
		customerInvestment.CustomerID,
		customerInvestment.InvestmentID,
		customerInvestment.Units)
	return err
}

func (r *mysqlCustomerInvestmentRepository) GetByCustomerAndInvestment(ctx context.Context, customerID, investmentID string) (*domain.CustomerInvestment, error) {
	query := `
		SELECT id, customer_id, investment_id, units 
		FROM customer_investments 
		WHERE customer_id = ? AND investment_id = ?
	`

	var customerInvestment domain.CustomerInvestment
	err := r.db.QueryRowContext(ctx, query, customerID, investmentID).Scan(
		&customerInvestment.ID,
		&customerInvestment.CustomerID,
		&customerInvestment.InvestmentID,
		&customerInvestment.Units)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("customer investment not found")
		}
		return nil, err
	}

	return &customerInvestment, nil
}

func (r *mysqlCustomerInvestmentRepository) UpdateUnits(ctx context.Context, id string, unitsChange float64) error {
	query := "UPDATE customer_investments SET units = units + ? WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, unitsChange, id)
	return err
}

func (r *mysqlCustomerInvestmentRepository) GetCustomerPortfolio(ctx context.Context, customerID, investmentID string) (*domain.CustomerPortfolio, error) {
	// Get customer
	var customer domain.Customer
	customerQuery := "SELECT id, name FROM customers WHERE id = ?"
	err := r.db.QueryRowContext(ctx, customerQuery, customerID).Scan(&customer.ID, &customer.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("customer not found")
		}
		return nil, err
	}

	// Get investment
	var investment domain.Investment
	investmentQuery := "SELECT id, name, total_units, total_balance FROM investments WHERE id = ?"
	err = r.db.QueryRowContext(ctx, investmentQuery, investmentID).Scan(
		&investment.ID, &investment.Name, &investment.TotalUnits, &investment.TotalBalance)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.New("investment not found")
		}
		return nil, err
	}

	investment.NAB = validateNAB(investment.TotalBalance, investment.TotalUnits)

	// Get customer investment
	portfolio := domain.CustomerPortfolio{
		Customer:   customer,
		Investment: investment,
	}

	query := `
		SELECT id, units FROM customer_investments 
		WHERE customer_id = ? AND investment_id = ?
	`
	var units float64
	var portfolioID string
	err = r.db.QueryRowContext(ctx, query, customerID, investmentID).Scan(&portfolioID, &units)
	if err != nil && err != sql.ErrNoRows {
		return nil, err
	}

	portfolio.Portfolio.ID = portfolioID
	portfolio.Portfolio.Units = units
	portfolio.Portfolio.Balance = roundDown(units*investment.NAB, 2)

	return &portfolio, nil
}
