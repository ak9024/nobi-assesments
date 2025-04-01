package mysql

import (
	"context"
	"database/sql"
	"nobi-assesment/internal/domain"
	"nobi-assesment/internal/repository"
)

type mysqlCustomerRepository struct {
	db *sql.DB
}

func NewMySQLCustomerRepository(db *sql.DB) repository.CustomerRepository {
	return &mysqlCustomerRepository{db}
}

func (r *mysqlCustomerRepository) Create(ctx context.Context, customer *domain.Customer) error {
	query := "INSERT INTO customers (id, name, is_active) VALUES (?, ?, ?)"
	_, err := r.db.ExecContext(ctx, query, customer.ID, customer.Name, customer.IsActive)
	return err
}

func (r *mysqlCustomerRepository) GetByID(ctx context.Context, id string) (*domain.Customer, error) {
	query := "SELECT id, name, is_active FROM customers WHERE id = ?"

	var customer domain.Customer
	err := r.db.QueryRowContext(ctx, query, id).Scan(&customer.ID, &customer.Name, &customer.IsActive)
	if err != nil {
		return nil, err
	}

	// Calculate balance and units
	query = `
		SELECT ci.units, i.total_balance, i.total_units 
		FROM customer_investments ci
		JOIN investments i ON ci.investment_id = i.id
		WHERE ci.customer_id = ?
	`
	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var totalBalance, totalUnits float64
	for rows.Next() {
		var units, investTotalBalance, investTotalUnits float64
		if err := rows.Scan(&units, &investTotalBalance, &investTotalUnits); err != nil {
			continue
		}

		var nab float64 = 1.0
		if investTotalUnits > 0 {
			nab = roundDown(investTotalBalance/investTotalUnits, 4)
		}

		customerBalance := roundDown(units*nab, 2)
		totalBalance += customerBalance
		totalUnits += units
	}

	customer.Balance = totalBalance
	customer.Units = totalUnits

	return &customer, nil
}

func (r *mysqlCustomerRepository) GetAll(ctx context.Context) ([]*domain.Customer, error) {
	query := "SELECT id, name, is_active FROM customers"
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	customers := []*domain.Customer{}
	for rows.Next() {
		var customer domain.Customer
		if err := rows.Scan(&customer.ID, &customer.Name, &customer.IsActive); err != nil {
			return nil, err
		}

		customers = append(customers, &customer)
	}

	// Calculate balance and units for each customer
	for _, customer := range customers {
		query = `
			SELECT ci.units, i.total_balance, i.total_units 
			FROM customer_investments ci
			JOIN investments i ON ci.investment_id = i.id
			WHERE ci.customer_id = ?
		`
		investRows, err := r.db.QueryContext(ctx, query, customer.ID)
		if err != nil {
			continue
		}

		var totalBalance, totalUnits float64
		for investRows.Next() {
			var units, investTotalBalance, investTotalUnits float64
			if err := investRows.Scan(&units, &investTotalBalance, &investTotalUnits); err != nil {
				continue
			}

			var nab float64 = 1.0
			if investTotalUnits > 0 {
				nab = roundDown(investTotalBalance/investTotalUnits, 4)
			}

			customerBalance := roundDown(units*nab, 2)
			totalBalance += customerBalance
			totalUnits += units
		}
		investRows.Close()

		customer.Balance = totalBalance
		customer.Units = totalUnits
	}

	return customers, nil
}

func (r *mysqlCustomerRepository) UpdateActiveStatus(ctx context.Context, id string, isActive bool) error {
	query := "UPDATE customers SET is_active = ? WHERE id = ?"
	_, err := r.db.ExecContext(ctx, query, isActive, id)
	return err
}
