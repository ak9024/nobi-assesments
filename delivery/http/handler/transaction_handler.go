package handler

import (
	"nobi-assesment/internal/domain"
	"nobi-assesment/internal/usecase"

	"github.com/gofiber/fiber/v2"
)

type TransactionHandler struct {
	transactionUsecase usecase.TransactionUsecase
}

func NewTransactionHandler(transactionUsecase usecase.TransactionUsecase) *TransactionHandler {
	return &TransactionHandler{
		transactionUsecase: transactionUsecase,
	}
}

func (h *TransactionHandler) Deposit(c *fiber.Ctx) error {
	var req domain.DepositRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	resp, err := h.transactionUsecase.Deposit(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *TransactionHandler) Withdraw(c *fiber.Ctx) error {
	var req domain.WithdrawRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	resp, err := h.transactionUsecase.Withdraw(c.Context(), &req)
	if err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	return c.Status(fiber.StatusOK).JSON(resp)
}

func (h *TransactionHandler) GetCustomerTransactions(c *fiber.Ctx) error {
	customerID := c.Params("id")

	transactions, err := h.transactionUsecase.GetCustomerTransactions(c.Context(), customerID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(transactions)
}

func (h *TransactionHandler) GetCustomerPortfolio(c *fiber.Ctx) error {
	customerID := c.Params("customer_id")
	investmentID := c.Params("investment_id")

	portfolio, err := h.transactionUsecase.GetCustomerPortfolio(c.Context(), customerID, investmentID)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": err.Error()})
	}

	return c.JSON(portfolio)
}
