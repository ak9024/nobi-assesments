package usecase

import (
	"context"
	"nobi-assesment/internal/domain"
	"nobi-assesment/internal/repository"
)

type InvestmentUsecase interface {
	Create(ctx context.Context, investment *domain.Investment) error
	GetByID(ctx context.Context, id string) (*domain.Investment, error)
	GetAll(ctx context.Context) ([]*domain.Investment, error)
}

type investmentUsecase struct {
	investmentRepo repository.InvestmentRepository
}

func NewInvestmentUsecase(investmentRepo repository.InvestmentRepository) InvestmentUsecase {
	return &investmentUsecase{
		investmentRepo: investmentRepo,
	}
}

func (u *investmentUsecase) Create(ctx context.Context, investment *domain.Investment) error {
	investment.NAB = 1.0
	investment.TotalUnits = 0
	investment.TotalBalance = 0
	return u.investmentRepo.Create(ctx, investment)
}

func (u *investmentUsecase) GetByID(ctx context.Context, id string) (*domain.Investment, error) {
	return u.investmentRepo.GetByID(ctx, id)
}

func (u *investmentUsecase) GetAll(ctx context.Context) ([]*domain.Investment, error) {
	return u.investmentRepo.GetAll(ctx)
}
