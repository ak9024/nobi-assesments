CREATE TABLE IF NOT EXISTS customers (
    id VARCHAR(36) NOT NULL DEFAULT (uuid()),  -- Primary key using UUID format
    name VARCHAR(255) NOT NULL,         -- Full name of the customer
    email VARCHAR(255) UNIQUE,          -- Customer email address
    phone VARCHAR(20),                  -- Customer phone number
    is_active BOOLEAN DEFAULT TRUE,     -- Status flag (true=active, false=inactive)
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- Record creation timestamp
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, -- Last update timestamp
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS investments (
    id VARCHAR(36) NOT NULL DEFAULT (uuid()),  -- Primary key using UUID format
    name VARCHAR(255) NOT NULL,              -- Name of the investment
    description TEXT,                        -- Detailed description of the investment
    risk_level ENUM('LOW', 'MEDIUM', 'HIGH') DEFAULT 'MEDIUM', -- Risk classification
    total_units DECIMAL(20,4) DEFAULT 0,     -- Total units of investment owned
    total_balance DECIMAL(20,2) DEFAULT 0,   -- Total monetary value of investment
    current_nab DECIMAL(20,4) DEFAULT 0,     -- Current Net Asset Value per unit
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- Record creation timestamp
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, -- Last update timestamp
    PRIMARY KEY (id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS customer_investments (
    id VARCHAR(36) NOT NULL DEFAULT (uuid()),  -- Primary key using UUID format
    customer_id VARCHAR(36) NOT NULL,        -- Reference to customers table
    investment_id VARCHAR(36) NOT NULL,      -- Reference to investments table
    units DECIMAL(20,4) DEFAULT 0,           -- Number of investment units owned by customer
    balance DECIMAL(20,2) DEFAULT 0,         -- Monetary value of customer's investment
    purchase_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- When the customer first invested
    last_transaction_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- Date of last transaction
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- Record creation timestamp
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, -- Last update timestamp
    PRIMARY KEY (id),
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    FOREIGN KEY (investment_id) REFERENCES investments(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    UNIQUE KEY unique_customer_investment (customer_id, investment_id)  -- Ensures a customer can't have duplicate investments
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

CREATE TABLE IF NOT EXISTS transactions (
    id VARCHAR(36) NOT NULL DEFAULT (uuid()),  -- Primary key using UUID format
    customer_id VARCHAR(36) NOT NULL,        -- Reference to customer who made the transaction
    investment_id VARCHAR(36) NOT NULL,      -- Reference to investment involved in transaction
    type ENUM('DEPOSIT', 'WITHDRAW') NOT NULL, -- Transaction type (buying or selling)
    status ENUM('PENDING', 'COMPLETED', 'FAILED', 'CANCELLED') DEFAULT 'PENDING', -- Transaction status
    amount DECIMAL(20,2) NOT NULL,           -- Monetary value of the transaction
    units DECIMAL(20,4) NOT NULL,            -- Number of investment units involved
    nab DECIMAL(20,4) NOT NULL,              -- Net Asset Value per unit at transaction time
    transaction_date TIMESTAMP DEFAULT CURRENT_TIMESTAMP, -- When the transaction occurred
    completed_date TIMESTAMP NULL,           -- When the transaction was completed
    notes TEXT,                              -- Additional transaction notes
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,  -- Record creation timestamp
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP, -- Last update timestamp
    PRIMARY KEY (id),
    FOREIGN KEY (customer_id) REFERENCES customers(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    FOREIGN KEY (investment_id) REFERENCES investments(id) ON DELETE RESTRICT ON UPDATE CASCADE,
    INDEX idx_transaction_date (transaction_date),
    INDEX idx_customer_investment (customer_id, investment_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

