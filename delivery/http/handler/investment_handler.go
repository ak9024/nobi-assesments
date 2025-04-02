package handler

import (
	"nobi-assesment/internal/domain"
	"nobi-assesment/internal/usecase"
	"nobi-assesment/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type InvestmentHandler struct {
	investmentUsecase usecase.InvestmentUsecase
}

func NewInvestmentHandler(investmentUsecase usecase.InvestmentUsecase) *InvestmentHandler {
	return &InvestmentHandler{
		investmentUsecase: investmentUsecase,
	}
}

func (h *InvestmentHandler) Create(c *fiber.Ctx) error {
	investment := new(domain.Investment)
	if err := c.BodyParser(investment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	if investment.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Investment product name is required"})
	}

	if investment.NAB <= 0 {
		investment.NAB = 1
	}

	if investment.TotalBalance <= 0 {
		investment.TotalBalance = 0
	}

	if investment.TotalUnits <= 0 {
		investment.TotalUnits = 0
	}

	investment.ID = utils.GenerateUUID()

	err := h.investmentUsecase.Create(c.Context(), investment)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save investment product"})
	}

	return c.Status(fiber.StatusCreated).JSON(investment)
}

func (h *InvestmentHandler) GetAll(c *fiber.Ctx) error {
	investments, err := h.investmentUsecase.GetAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve investment products"})
	}

	return c.JSON(investments)
}

func (h *InvestmentHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")

	investment, err := h.investmentUsecase.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Investment product not found"})
	}

	return c.JSON(investment)
}
