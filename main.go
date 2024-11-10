package main

import (
	"github.com/alexhendra/amartha_loan/controllers"
	"github.com/alexhendra/amartha_loan/repositories"
	"github.com/alexhendra/amartha_loan/services"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	// Initialize repositories, services, and controllers
	loanRepo := repositories.NewLoanRepository()
	approvalRepo := repositories.NewApprovalRepository()
	investmentRepo := repositories.NewInvestmentRepository()
	disbursementRepo := repositories.NewDisbursementRepository()
	notificationService := services.NewNotificationService()

	loanService := services.LoanService{
		LoanRepo:            loanRepo,
		ApprovalRepo:        approvalRepo,
		InvestmentRepo:      investmentRepo,
		DisbursementRepo:    disbursementRepo,
		NotificationService: notificationService,
	}
	loanController := controllers.LoanController{LoanService: &loanService}

	// Define routes
	e.POST("/loans", loanController.CreateLoan)
	e.PUT("/loans/:id/approve", loanController.ApproveLoan)
	e.PUT("/loans/:id/invest", loanController.InvestLoan)
	e.PUT("/loans/:id/disburse", loanController.DisburseLoan)

	// Start server
	e.Start(":8080")
}
