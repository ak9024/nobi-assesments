package http

import (
	"nobi-assesment/delivery/http/handler"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
)

// SetupRoutes configures all the routes for the API
func SetupRoutes(
	app *fiber.App,
	customerHandler *handler.CustomerHandler,
	investmentHandler *handler.InvestmentHandler,
	transactionHandler *handler.TransactionHandler,
) {
	// Middleware
	app.Use(logger.New())

	// API routes
	api := app.Group("/api")

	// Customer routes
	customers := api.Group("/customers")
	customers.Post("/", customerHandler.Create)
	customers.Get("/", customerHandler.GetAll)
	customers.Get("/:id", customerHandler.GetByID)

	// Investment routes
	investments := api.Group("/investments")
	investments.Post("/", investmentHandler.Create)
	investments.Get("/", investmentHandler.GetAll)
	investments.Get("/:id", investmentHandler.GetByID)

	// Transaction routes
	transactions := api.Group("/transactions")
	transactions.Post("/deposit", transactionHandler.Deposit)
	transactions.Post("/withdraw", transactionHandler.Withdraw)
	transactions.Get("/customer/:id", transactionHandler.GetCustomerTransactions)

	// Portfolio route
	api.Get("/portfolio/:customer_id/:investment_id", transactionHandler.GetCustomerPortfolio)
}
