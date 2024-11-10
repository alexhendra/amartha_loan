package utils

import (
	"log"

	"github.com/alexhendra/amartha_loan/loan_service/models"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

// initializeDB sets up and returns a GORM database connection.
func InitializeDB() *gorm.DB {
	// Here we're using SQLite for simplicity. Adjust the driver and DSN as needed for other databases.
	dsn := "loanengine.db" // Database file for SQLite
	db, err := gorm.Open(sqlite.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("failed to connect to database: %v", err)
	}

	// Run migrations to create/update the database schema
	if err := db.AutoMigrate(&models.Loan{}, &models.Approval{}, &models.Investment{}, &models.Disbursement{}); err != nil {
		log.Fatalf("failed to run migrations: %v", err)
	}

	return db
}
