package main

import (
	"log"
	"nobi-assesment/delivery/http"
	"nobi-assesment/delivery/http/handler"
	"nobi-assesment/delivery/http/middleware"
	"nobi-assesment/internal/repository/mysql"
	"nobi-assesment/internal/usecase"
	"nobi-assesment/pkg/db"
	"os"

	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
)

func main() {
	// Load environment variables
	err := godotenv.Load()
	if err != nil {
		log.Println("Warning: .env file not found, using environment variables")
	}

	// Database configuration
	dbUser := getEnv("DB_USER", "root")
	dbPass := getEnv("DB_PASSWORD", "29/jSGGz&x0c")
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "3306")
	dbName := getEnv("DB_NAME", "nobi_investment")

	// Initialize database connection
	dbConn, err := db.NewMySQLConnection(dbUser, dbPass, dbHost, dbPort, dbName)
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer dbConn.Close()

	// Repository layer
	customerRepo := mysql.NewMySQLCustomerRepository(dbConn)
	investmentRepo := mysql.NewMySQLInvestmentRepository(dbConn)
	custInvestRepo := mysql.NewMySQLCustomerInvestmentRepository(dbConn)
	transactionRepo := mysql.NewMySQLTransactionRepository(dbConn)

	// Usecase layer
	customerUsecase := usecase.NewCustomerUsecase(customerRepo)
	investmentUsecase := usecase.NewInvestmentUsecase(investmentRepo)
	transactionUsecase := usecase.NewTransactionUsecase(
		transactionRepo,
		customerRepo,
		investmentRepo,
		custInvestRepo,
		dbConn,
	)

	// Handler layer
	customerHandler := handler.NewCustomerHandler(customerUsecase)
	investmentHandler := handler.NewInvestmentHandler(investmentUsecase)
	transactionHandler := handler.NewTransactionHandler(transactionUsecase)

	// Initialize Fiber app
	app := fiber.New(fiber.Config{
		ErrorHandler: customErrorHandler,
	})

	// Setup middleware
	middleware.SetupMiddleware(app)

	// Setup routes
	http.SetupRoutes(app, customerHandler, investmentHandler, transactionHandler)

	// Start server
	port := getEnv("PORT", "3000")
	log.Printf("Server started on http://localhost:%s", port)
	log.Fatal(app.Listen(":" + port))
}

// Custom error handler for Fiber
func customErrorHandler(c *fiber.Ctx, err error) error {
	code := fiber.StatusInternalServerError

	// Check if it's a Fiber error
	if e, ok := err.(*fiber.Error); ok {
		code = e.Code
	}

	return c.Status(code).JSON(fiber.Map{
		"error": err.Error(),
	})
}

// Helper function to get environment variables with fallback
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
