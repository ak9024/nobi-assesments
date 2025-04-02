package handler

import (
	"nobi-assesment/internal/domain"
	"nobi-assesment/internal/usecase"
	"nobi-assesment/pkg/utils"

	"github.com/gofiber/fiber/v2"
)

type CustomerHandler struct {
	customerUsecase usecase.CustomerUsecase
}

func NewCustomerHandler(customerUsecase usecase.CustomerUsecase) *CustomerHandler {
	return &CustomerHandler{
		customerUsecase: customerUsecase,
	}
}

func (h *CustomerHandler) Create(c *fiber.Ctx) error {
	customer := new(domain.Customer)
	if err := c.BodyParser(customer); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	if customer.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Customer name is required"})
	}

	// Check if customer with the same name already exists
	existingCustomers, err := h.customerUsecase.GetAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to check customer uniqueness"})
	}

	for _, existingCustomer := range existingCustomers {
		if existingCustomer.Name == customer.Name {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Customer name must be unique"})
		}
	}

	customer.ID = utils.GenerateUUID()

	err = h.customerUsecase.Create(c.Context(), customer)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save customer"})
	}

	return c.Status(fiber.StatusCreated).JSON(customer)
}

func (h *CustomerHandler) GetAll(c *fiber.Ctx) error {
	customers, err := h.customerUsecase.GetAll(c.Context())
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve customers"})
	}

	return c.JSON(customers)
}

func (h *CustomerHandler) GetByID(c *fiber.Ctx) error {
	id := c.Params("id")

	customer, err := h.customerUsecase.GetByID(c.Context(), id)
	if err != nil {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Customer not found"})
	}

	return c.JSON(customer)
}
