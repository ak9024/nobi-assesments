package main

import (
	"database/sql"
	"log"
	"math"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/logger"
	"github.com/google/uuid"
)

type Customer struct {
	ID       string  `json:"id"`
	Name     string  `json:"name"`
	Balance  float64 `json:"balance"`
	Units    float64 `json:"units"`
	IsActive bool    `json:"is_active"`
}

type Investment struct {
	ID           string  `json:"id"`
	Name         string  `json:"name"`
	TotalUnits   float64 `json:"total_units"`
	TotalBalance float64 `json:"total_balance"`
	NAB          float64 `json:"nab"`
}

type Transaction struct {
	ID           string  `json:"id"`
	CustomerID   string  `json:"customer_id"`
	InvestmentID string  `json:"investment_id"`
	Type         string  `json:"type"` // DEPOSIT or WITHDRAW
	Amount       float64 `json:"amount"`
	Units        float64 `json:"units"`
	NAB          float64 `json:"nab"`
	Timestamp    string  `json:"timestamp"`
}

func roundDown(value float64, places int) float64 {
	shift := math.Pow(10, float64(places))
	return math.Floor(value*shift) / shift
}

func validateNAB(totalBalance, totalUnits float64) float64 {
	if totalUnits <= 0 {
		return 1.0
	}

	nab := roundDown(totalBalance/totalUnits, 4)

	log.Printf("Calculated NAB: %.4f, Total Balance: %.2f, Total Units: %.4f", nab, totalBalance, totalUnits)

	return nab
}

func generateUUID() string {
	return uuid.New().String()
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("mysql", "root:29/jSGGz&x0c@tcp(localhost:3306)/nobi_investment?parseTime=true")
	if err != nil {
		log.Fatal("Gagal koneksi ke database:", err)
	}
	defer db.Close()

	app := fiber.New()
	app.Use(logger.New())

	customerRoutes := app.Group("/api/customers")
	customerRoutes.Post("/", createCustomer)
	customerRoutes.Get("/", getAllCustomers)
	customerRoutes.Get("/:id", getCustomerByID)

	investmentRoutes := app.Group("/api/investments")
	investmentRoutes.Post("/", createInvestment)
	investmentRoutes.Get("/", getAllInvestments)
	investmentRoutes.Get("/:id", getInvestmentByID)

	transactionRoutes := app.Group("/api/transactions")
	transactionRoutes.Post("/deposit", depositFunds)
	transactionRoutes.Post("/withdraw", withdrawFunds)
	transactionRoutes.Get("/customer/:id", getCustomerTransactions)

	app.Get("/api/portfolio/:customer_id/:investment_id", getCustomerPortfolio)

	log.Println("Server started on http://localhost:3000")
	log.Fatal(app.Listen(":3000"))
}

func createCustomer(c *fiber.Ctx) error {
	customer := new(Customer)
	if err := c.BodyParser(customer); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	if customer.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Customer name is required"})
	}

	customer.ID = generateUUID()
	customer.IsActive = true

	_, err := db.Exec("INSERT INTO customers (id, name) VALUES (?, ?)", customer.ID, customer.Name)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save customer"})
	}

	return c.Status(fiber.StatusCreated).JSON(customer)
}

func getAllCustomers(c *fiber.Ctx) error {
	rows, err := db.Query(`
		SELECT c.id, c.name, c.is_active 
		FROM customers c
	`)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil data nasabah"})
	}
	defer rows.Close()

	customers := []Customer{}
	for rows.Next() {
		var customer Customer
		if err := rows.Scan(&customer.ID, &customer.Name, &customer.IsActive); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal membaca data nasabah"})
		}

		// Get total balance and units from all investments for this customer
		var totalBalance, totalUnits float64
		investRows, err := db.Query(`
			SELECT ci.units, i.total_balance, i.total_units 
			FROM customer_investments ci
			JOIN investments i ON ci.investment_id = i.id
			WHERE ci.customer_id = ?
		`, customer.ID)

		if err == nil {
			defer investRows.Close()
			for investRows.Next() {
				var units, investTotalBalance, investTotalUnits float64
				if err := investRows.Scan(&units, &investTotalBalance, &investTotalUnits); err != nil {
					continue
				}

				// Calculate NAB and customer's balance for this investment
				var nab float64 = 1.0
				if investTotalUnits > 0 {
					nab = roundDown(investTotalBalance/investTotalUnits, 4)
				}

				customerBalance := roundDown(units*nab, 2)
				totalBalance += customerBalance
				totalUnits += units
			}
		}

		customer.Balance = totalBalance
		customer.Units = totalUnits
		customers = append(customers, customer)
	}

	return c.JSON(customers)
}

func getCustomerByID(c *fiber.Ctx) error {
	id := c.Params("id")

	var customer Customer
	err := db.QueryRow("SELECT id, name, is_active FROM customers WHERE id = ?", id).Scan(&customer.ID, &customer.Name, &customer.IsActive)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Nasabah tidak ditemukan"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mengambil data nasabah"})
	}

	// Get total balance and units from all investments for this customer
	var totalBalance, totalUnits float64
	investRows, err := db.Query(`
		SELECT ci.units, i.total_balance, i.total_units 
		FROM customer_investments ci
		JOIN investments i ON ci.investment_id = i.id
		WHERE ci.customer_id = ?
	`, id)

	if err == nil {
		defer investRows.Close()
		for investRows.Next() {
			var units, investTotalBalance, investTotalUnits float64
			if err := investRows.Scan(&units, &investTotalBalance, &investTotalUnits); err != nil {
				continue
			}

			// Calculate NAB and customer's balance for this investment
			var nab float64 = 1.0
			if investTotalUnits > 0 {
				nab = roundDown(investTotalBalance/investTotalUnits, 4)
			}

			customerBalance := roundDown(units*nab, 2)
			totalBalance += customerBalance
			totalUnits += units
		}
	}

	customer.Balance = totalBalance
	customer.Units = totalUnits

	return c.JSON(customer)
}

func createInvestment(c *fiber.Ctx) error {
	investment := new(Investment)
	if err := c.BodyParser(investment); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	if investment.Name == "" {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Investment product name is required"})
	}

	investment.ID = generateUUID()

	investment.NAB = 1.0
	investment.TotalUnits = 0
	investment.TotalBalance = 0

	_, err := db.Exec("INSERT INTO investments (id, name, total_units, total_balance) VALUES (?, ?, ?, ?)",
		investment.ID, investment.Name, investment.TotalUnits, investment.TotalBalance)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to save investment product"})
	}

	return c.Status(fiber.StatusCreated).JSON(investment)
}

func getAllInvestments(c *fiber.Ctx) error {
	rows, err := db.Query(`
		SELECT id, name, total_units, total_balance FROM investments
	`)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve investment products data"})
	}
	defer rows.Close()

	investments := []Investment{}
	for rows.Next() {
		var investment Investment
		if err := rows.Scan(&investment.ID, &investment.Name, &investment.TotalUnits, &investment.TotalBalance); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to read investment products data"})
		}

		if investment.TotalUnits > 0 {
			investment.NAB = roundDown(investment.TotalBalance/investment.TotalUnits, 4)
		} else {
			investment.NAB = 1.0
		}

		investments = append(investments, investment)
	}

	return c.JSON(investments)
}

func getInvestmentByID(c *fiber.Ctx) error {
	id := c.Params("id")

	var investment Investment
	err := db.QueryRow(`
		SELECT id, name, total_units, total_balance FROM investments WHERE id = ?
	`, id).Scan(&investment.ID, &investment.Name, &investment.TotalUnits, &investment.TotalBalance)

	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Investment product not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve investment product data"})
	}

	investment.NAB = validateNAB(investment.TotalBalance, investment.TotalUnits)

	return c.JSON(investment)
}

func depositFunds(c *fiber.Ctx) error {
	var request struct {
		CustomerID   string  `json:"customer_id"`
		InvestmentID string  `json:"investment_id"`
		Amount       float64 `json:"amount"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	if request.CustomerID == "" || request.InvestmentID == "" || request.Amount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid parameters"})
	}

	tx, err := db.Begin()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to start database transaction"})
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var isActive bool
	var customerName string
	err = tx.QueryRow("SELECT name, is_active FROM customers WHERE id = ?", request.CustomerID).Scan(&customerName, &isActive)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Customer not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to check customer"})
	}
	if !isActive {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Customer is not active"})
	}

	log.Printf("Customer %s is depositing funds of Rp. %.2f", customerName, request.Amount)

	var investment Investment
	err = tx.QueryRow(`
		SELECT id, name, total_units, total_balance FROM investments WHERE id = ?
	`, request.InvestmentID).Scan(&investment.ID, &investment.Name, &investment.TotalUnits, &investment.TotalBalance)

	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Investment product not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve investment data"})
	}

	var currentNAB float64 = validateNAB(investment.TotalBalance, investment.TotalUnits)

	newUnits := roundDown(request.Amount/currentNAB, 4)

	_, err = tx.Exec(`
		UPDATE investments 
		SET total_balance = total_balance + ?, total_units = total_units + ? 
		WHERE id = ?
	`, request.Amount, newUnits, request.InvestmentID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update investment data"})
	}

	var existingUnits float64
	var portofolioExists bool
	var portofolioID string
	err = tx.QueryRow(`
		SELECT id, units FROM customer_investments 
		WHERE customer_id = ? AND investment_id = ?
	`, request.CustomerID, request.InvestmentID).Scan(&portofolioID, &existingUnits)

	if err != nil {
		if err == sql.ErrNoRows {
			portofolioExists = false
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal memeriksa portofolio nasabah"})
		}
	} else {
		portofolioExists = true
	}

	if portofolioExists {
		_, err = tx.Exec(`
			UPDATE customer_investments 
			SET units = units + ? 
			WHERE id = ?
		`, newUnits, portofolioID)
	} else {
		portofolioID = generateUUID()
		_, err = tx.Exec(`
			INSERT INTO customer_investments (id, customer_id, investment_id, units) 
			VALUES (?, ?, ?, ?)
		`, portofolioID, request.CustomerID, request.InvestmentID, newUnits)
	}

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal update portofolio nasabah"})
	}

	transactionID := generateUUID()
	_, err = tx.Exec(`
		INSERT INTO transactions (id, customer_id, investment_id, type, amount, units, nab) 
		VALUES (?, ?, ?, 'DEPOSIT', ?, ?, ?)
	`, transactionID, request.CustomerID, request.InvestmentID, request.Amount, newUnits, currentNAB)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal mencatat transaksi"})
	}

	if err = tx.Commit(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Gagal commit transaksi database"})
	}

	var totalUnitsAfterDeposit float64
	if portofolioExists {
		totalUnitsAfterDeposit = existingUnits + newUnits
	} else {
		totalUnitsAfterDeposit = newUnits
	}

	currentBalance := roundDown(totalUnitsAfterDeposit*currentNAB, 2)

	log.Printf("Deposit: Amount: Rp. %.2f, NAB: %.4f, Units added: %.4f, Total units: %.4f, Balance: Rp. %.2f",
		request.Amount, currentNAB, newUnits, totalUnitsAfterDeposit, currentBalance)

	response := fiber.Map{
		"transaction_id":  transactionID,
		"message":         "Penyetoran dana berhasil",
		"amount":          request.Amount,
		"units_added":     newUnits,
		"nab":             currentNAB,
		"total_units":     totalUnitsAfterDeposit,
		"current_balance": currentBalance,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func withdrawFunds(c *fiber.Ctx) error {
	var request struct {
		CustomerID   string  `json:"customer_id"`
		InvestmentID string  `json:"investment_id"`
		Amount       float64 `json:"amount"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Cannot parse JSON"})
	}

	if request.CustomerID == "" || request.InvestmentID == "" || request.Amount <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid parameters"})
	}

	tx, err := db.Begin()
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to start database transaction"})
	}
	defer func() {
		if err != nil {
			tx.Rollback()
		}
	}()

	var isActive bool
	var customerName string
	err = tx.QueryRow("SELECT name, is_active FROM customers WHERE id = ?", request.CustomerID).Scan(&customerName, &isActive)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Customer not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to check customer"})
	}
	if !isActive {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Customer is not active"})
	}

	log.Printf("Customer %s is withdrawing funds of Rp. %.2f", customerName, request.Amount)

	var investment Investment
	err = tx.QueryRow(`
		SELECT id, name, total_units, total_balance FROM investments WHERE id = ?
	`, request.InvestmentID).Scan(&investment.ID, &investment.Name, &investment.TotalUnits, &investment.TotalBalance)

	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Investment product not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve investment data"})
	}

	var currentNAB float64
	if investment.TotalUnits > 0 {
		currentNAB = roundDown(investment.TotalBalance/investment.TotalUnits, 4)
	} else {
		currentNAB = 1.0
	}

	withdrawUnits := roundDown(request.Amount/currentNAB, 4)

	var existingUnits float64
	var portofolioID string
	err = tx.QueryRow(`
		SELECT id, units FROM customer_investments 
		WHERE customer_id = ? AND investment_id = ?
	`, request.CustomerID, request.InvestmentID).Scan(&portofolioID, &existingUnits)

	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Customer does not have this investment"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to check customer portfolio"})
	}

	if withdrawUnits > existingUnits {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Insufficient balance for withdrawal"})
	}

	_, err = tx.Exec(`
		UPDATE investments 
		SET total_balance = total_balance - ?, total_units = total_units - ? 
		WHERE id = ?
	`, request.Amount, withdrawUnits, request.InvestmentID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update data investments"})
	}

	_, err = tx.Exec(`
		UPDATE customer_investments 
		SET units = units - ? 
		WHERE id = ?
	`, withdrawUnits, portofolioID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update customer portfolio"})
	}

	transactionID := generateUUID()
	_, err = tx.Exec(`
		INSERT INTO transactions (id, customer_id, investment_id, type, amount, units, nab) 
		VALUES (?, ?, ?, 'WITHDRAW', ?, ?, ?)
	`, transactionID, request.CustomerID, request.InvestmentID, request.Amount, withdrawUnits, currentNAB)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to insert data transactions"})
	}

	if err = tx.Commit(); err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to commit transactions"})
	}

	remainingUnits := existingUnits - withdrawUnits
	currentBalance := roundDown(remainingUnits*currentNAB, 2)
	log.Printf("Withdrawal: Amount: Rp. %.2f, NAB: %.4f, Units reduced: %.4f, Remaining units: %.4f, Balance: Rp. %.2f",
		request.Amount, currentNAB, withdrawUnits, remainingUnits, currentBalance)
	response := fiber.Map{
		"transaction_id":  transactionID,
		"message":         "Withdrawal successful",
		"amount":          request.Amount,
		"units_reduced":   withdrawUnits,
		"nab":             currentNAB,
		"remaining_units": remainingUnits,
		"current_balance": currentBalance,
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func getCustomerPortfolio(c *fiber.Ctx) error {
	customerID := c.Params("customer_id")
	investmentID := c.Params("investment_id")

	var customer Customer
	err := db.QueryRow("SELECT id, name FROM customers WHERE id = ?", customerID).Scan(&customer.ID, &customer.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Customer not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve customer data"})
	}

	var investment Investment
	err = db.QueryRow(`
		SELECT id, name, total_units, total_balance FROM investments WHERE id = ?
	`, investmentID).Scan(&investment.ID, &investment.Name, &investment.TotalUnits, &investment.TotalBalance)

	if err != nil {
		if err == sql.ErrNoRows {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Investment product not found"})
		}
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve investment data"})
	}

	if investment.TotalUnits > 0 {
		investment.NAB = roundDown(investment.TotalBalance/investment.TotalUnits, 4)
	} else {
		investment.NAB = 1.0
	}

	var units float64
	var portofolioID string
	err = db.QueryRow(`
		SELECT id, units FROM customer_investments 
		WHERE customer_id = ? AND investment_id = ?
	`, customerID, investmentID).Scan(&portofolioID, &units)

	if err != nil {
		if err == sql.ErrNoRows {
			units = 0
		} else {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve portfolio data"})
		}
	}

	balance := roundDown(units*investment.NAB, 2)

	response := fiber.Map{
		"customer": fiber.Map{
			"id":   customer.ID,
			"name": customer.Name,
		},
		"investment": fiber.Map{
			"id":   investment.ID,
			"name": investment.Name,
			"nab":  investment.NAB,
		},
		"portfolio": fiber.Map{
			"id":      portofolioID,
			"units":   units,
			"balance": balance,
		},
	}

	return c.Status(fiber.StatusOK).JSON(response)
}

func getCustomerTransactions(c *fiber.Ctx) error {
	customerID := c.Params("id")

	var exists bool
	err := db.QueryRow("SELECT EXISTS(SELECT 1 FROM customers WHERE id = ?)", customerID).Scan(&exists)
	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to check customer"})
	}
	if !exists {
		return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "Customer not found"})
	}

	rows, err := db.Query(`
		SELECT t.id, t.investment_id, i.name, t.type, t.amount, t.units, t.nab, t.transaction_date
		FROM transactions t
		JOIN investments i ON t.investment_id = i.id
		WHERE t.customer_id = ?
		ORDER BY t.transaction_date DESC
	`, customerID)

	if err != nil {
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve transaction data"})
	}
	defer rows.Close()

	transactions := []fiber.Map{}
	for rows.Next() {
		var id, investmentID, investmentName, transactionType string
		var amount, units, nab float64
		var transactionDate string

		if err := rows.Scan(&id, &investmentID, &investmentName, &transactionType, &amount, &units, &nab, &transactionDate); err != nil {
			return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to read transaction data"})
		}

		transactions = append(transactions, fiber.Map{
			"id": id,
			"investment": fiber.Map{
				"id":   investmentID,
				"name": investmentName,
			},
			"type":             transactionType,
			"amount":           amount,
			"units":            units,
			"nab":              nab,
			"transaction_date": transactionDate,
		})
	}

	return c.JSON(transactions)
}
