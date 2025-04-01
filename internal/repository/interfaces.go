package repository

import (
	"context"
	"nobi-assesment/internal/domain"
)

type CustomerRepository interface {
	Create(ctx context.Context, customer *domain.Customer) error
	GetByID(ctx context.Context, id string) (*domain.Customer, error)
	GetAll(ctx context.Context) ([]*domain.Customer, error)
	UpdateActiveStatus(ctx context.Context, id string, isActive bool) error
}

type InvestmentRepository interface {
	Create(ctx context.Context, investment *domain.Investment) error
	GetByID(ctx context.Context, id string) (*domain.Investment, error)
	GetAll(ctx context.Context) ([]*domain.Investment, error)
	UpdateBalance(ctx context.Context, id string, amountChange float64, unitsChange float64) error
}

type CustomerInvestmentRepository interface {
	Create(ctx context.Context, customerInvestment *domain.CustomerInvestment) error
	GetByCustomerAndInvestment(ctx context.Context, customerID, investmentID string) (*domain.CustomerInvestment, error)
	UpdateUnits(ctx context.Context, id string, unitsChange float64) error
	GetCustomerPortfolio(ctx context.Context, customerID, investmentID string) (*domain.CustomerPortfolio, error)
}

type TransactionRepository interface {
	Create(ctx context.Context, transaction *domain.Transaction) error
	GetByCustomerID(ctx context.Context, customerID string) ([]*domain.Transaction, error)
}
