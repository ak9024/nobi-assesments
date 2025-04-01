package usecase

import (
	"context"
	"nobi-assesment/internal/domain"
	"nobi-assesment/internal/repository"
)

type CustomerUsecase interface {
	Create(ctx context.Context, customer *domain.Customer) error
	GetByID(ctx context.Context, id string) (*domain.Customer, error)
	GetAll(ctx context.Context) ([]*domain.Customer, error)
}

type customerUsecase struct {
	customerRepo repository.CustomerRepository
}

func NewCustomerUsecase(customerRepo repository.CustomerRepository) CustomerUsecase {
	return &customerUsecase{
		customerRepo: customerRepo,
	}
}

func (u *customerUsecase) Create(ctx context.Context, customer *domain.Customer) error {
	customer.IsActive = true
	return u.customerRepo.Create(ctx, customer)
}

func (u *customerUsecase) GetByID(ctx context.Context, id string) (*domain.Customer, error) {
	return u.customerRepo.GetByID(ctx, id)
}

func (u *customerUsecase) GetAll(ctx context.Context) ([]*domain.Customer, error) {
	return u.customerRepo.GetAll(ctx)
}
