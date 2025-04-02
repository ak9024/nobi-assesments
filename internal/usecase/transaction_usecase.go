package usecase

import (
	"context"
	"database/sql"
	"errors"
	"nobi-assesment/internal/domain"
	"nobi-assesment/internal/repository"
	"nobi-assesment/pkg/utils"
	"time"
)

type TransactionUsecase interface {
	Deposit(ctx context.Context, req *domain.DepositRequest) (*domain.TransactionResponse, error)
	Withdraw(ctx context.Context, req *domain.WithdrawRequest) (*domain.TransactionResponse, error)
	GetCustomerTransactions(ctx context.Context, customerID string) ([]*domain.Transaction, error)
	GetCustomerPortfolio(ctx context.Context, customerID, investmentID string) (*domain.CustomerPortfolio, error)
}

type transactionUsecase struct {
	transactionRepo repository.TransactionRepository
	customerRepo    repository.CustomerRepository
	investmentRepo  repository.InvestmentRepository
	custInvestRepo  repository.CustomerInvestmentRepository
	db              *sql.DB // For transactions
}

func NewTransactionUsecase(
	transactionRepo repository.TransactionRepository,
	customerRepo repository.CustomerRepository,
	investmentRepo repository.InvestmentRepository,
	custInvestRepo repository.CustomerInvestmentRepository,
	db *sql.DB,
) TransactionUsecase {
	return &transactionUsecase{
		transactionRepo: transactionRepo,
		customerRepo:    customerRepo,
		investmentRepo:  investmentRepo,
		custInvestRepo:  custInvestRepo,
		db:              db,
	}
}

func (u *transactionUsecase) Deposit(ctx context.Context, req *domain.DepositRequest) (*domain.TransactionResponse, error) {
	if req.CustomerID == "" || req.InvestmentID == "" || req.Amount <= 0 {
		return nil, errors.New("invalid parameters")
	}

	// Begin transaction
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Check customer
	customer, err := u.customerRepo.GetByID(ctx, req.CustomerID)
	if err != nil {
		return nil, err
	}
	if !customer.IsActive {
		return nil, errors.New("customer is not active")
	}

	// Get investment
	investment, err := u.investmentRepo.GetByID(ctx, req.InvestmentID)
	if err != nil {
		return nil, err
	}

	// Calculate NAB and new units
	var currentNAB float64
	if investment.NAB <= 0 {
		currentNAB = utils.ValidateNAB(investment.TotalBalance, investment.TotalUnits)
	} else {
		currentNAB = investment.NAB
	}
	newUnits := utils.RoundDown(req.Amount/currentNAB, 4)

	// Update investment
	err = u.investmentRepo.UpdateBalance(ctx, req.InvestmentID, req.Amount, newUnits)
	if err != nil {
		return nil, err
	}

	// Update or create customer investment
	var customerInvestment *domain.CustomerInvestment
	var totalUnitsAfterDeposit float64

	customerInvestment, err = u.custInvestRepo.GetByCustomerAndInvestment(ctx, req.CustomerID, req.InvestmentID)
	if err != nil {
		// Check if it's a "not found" error or empty data rows
		if err.Error() == "sql: no rows in result set" {
			// Create new customer investment
			newCustomerInvestment := &domain.CustomerInvestment{
				ID:           utils.GenerateUUID(),
				CustomerID:   req.CustomerID,
				InvestmentID: req.InvestmentID,
				Units:        newUnits,
			}
			err = u.custInvestRepo.Create(ctx, newCustomerInvestment)
			if err != nil {
				return nil, err
			}
			totalUnitsAfterDeposit = newUnits
		} else {
			return nil, err
		}
	} else {
		// Update existing customer investment
		err = u.custInvestRepo.UpdateUnits(ctx, customerInvestment.ID, newUnits)
		if err != nil {
			return nil, err
		}
		totalUnitsAfterDeposit = customerInvestment.Units + newUnits
	}

	// Create transaction record
	transaction := &domain.Transaction{
		ID:              utils.GenerateUUID(),
		CustomerID:      req.CustomerID,
		InvestmentID:    req.InvestmentID,
		Type:            "DEPOSIT",
		Amount:          req.Amount,
		Units:           newUnits,
		NAB:             currentNAB,
		TransactionDate: time.Now(),
	}

	err = u.transactionRepo.Create(ctx, transaction)
	if err != nil {
		return nil, err
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	currentBalance := utils.RoundDown(totalUnitsAfterDeposit*currentNAB, 2)

	return &domain.TransactionResponse{
		TransactionID:  transaction.ID,
		Message:        "Deposit successful",
		Amount:         req.Amount,
		Units:          newUnits,
		NAB:            currentNAB,
		TotalUnits:     totalUnitsAfterDeposit,
		CurrentBalance: currentBalance,
	}, nil
}

func (u *transactionUsecase) GetCustomerTransactions(ctx context.Context, customerID string) ([]*domain.Transaction, error) {
	// Verify customer exists
	_, err := u.customerRepo.GetByID(ctx, customerID)
	if err != nil {
		return nil, err
	}

	return u.transactionRepo.GetByCustomerID(ctx, customerID)
}

func (u *transactionUsecase) GetCustomerPortfolio(ctx context.Context, customerID, investmentID string) (*domain.CustomerPortfolio, error) {
	return u.custInvestRepo.GetCustomerPortfolio(ctx, customerID, investmentID)
}

func (u *transactionUsecase) Withdraw(ctx context.Context, req *domain.WithdrawRequest) (*domain.TransactionResponse, error) {
	if req.CustomerID == "" || req.InvestmentID == "" || req.Amount <= 0 {
		return nil, errors.New("invalid parameters")
	}

	// Begin transaction
	tx, err := u.db.BeginTx(ctx, nil)
	if err != nil {
		return nil, err
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	// Check customer
	customer, err := u.customerRepo.GetByID(ctx, req.CustomerID)
	if err != nil {
		return nil, err
	}
	if !customer.IsActive {
		return nil, errors.New("customer is not active")
	}

	// Get investment
	investment, err := u.investmentRepo.GetByID(ctx, req.InvestmentID)
	if err != nil {
		return nil, err
	}

	// Calculate NAB
	currentNAB := utils.ValidateNAB(investment.TotalBalance, investment.TotalUnits)
	withdrawUnits := utils.RoundDown(req.Amount/currentNAB, 4)

	// Get customer investment
	customerInvestment, err := u.custInvestRepo.GetByCustomerAndInvestment(ctx, req.CustomerID, req.InvestmentID)
	if err != nil {
		return nil, err
	}

	// Check sufficient balance
	if withdrawUnits > customerInvestment.Units {
		return nil, errors.New("insufficient balance for withdrawal")
	}

	// Update investment
	err = u.investmentRepo.UpdateBalance(ctx, req.InvestmentID, -req.Amount, -withdrawUnits)
	if err != nil {
		return nil, err
	}

	// Update customer investment
	err = u.custInvestRepo.UpdateUnits(ctx, customerInvestment.ID, -withdrawUnits)
	if err != nil {
		return nil, err
	}

	// Create transaction record
	transaction := &domain.Transaction{
		ID:              utils.GenerateUUID(),
		CustomerID:      req.CustomerID,
		InvestmentID:    req.InvestmentID,
		Type:            "WITHDRAW",
		Amount:          req.Amount,
		Units:           withdrawUnits,
		NAB:             currentNAB,
		TransactionDate: time.Now(),
	}

	err = u.transactionRepo.Create(ctx, transaction)
	if err != nil {
		return nil, err
	}

	// Commit transaction
	if err = tx.Commit(); err != nil {
		return nil, err
	}

	remainingUnits := customerInvestment.Units - withdrawUnits
	currentBalance := utils.RoundDown(remainingUnits*currentNAB, 2)

	return &domain.TransactionResponse{
		TransactionID:  transaction.ID,
		Message:        "Withdrawal successful",
		Amount:         req.Amount,
		UnitsReduced:   withdrawUnits,
		NAB:            currentNAB,
		RemainingUnits: remainingUnits,
		CurrentBalance: currentBalance,
	}, nil
}
