package main

import (
	routes "github.com/alexhendra/amartha_loan/loan_service/controllers"
	"github.com/alexhendra/amartha_loan/loan_service/repositories"
	"github.com/alexhendra/amartha_loan/loan_service/services"
	"github.com/alexhendra/amartha_loan/loan_service/utils"
	"github.com/labstack/echo/v4"
)

func main() {
	e := echo.New()

	db := utils.InitializeDB()
	// Initialize repositories, services, and controllers
	loanRepo := repositories.NewLoanRepository(db)
	approvalRepo := repositories.NewApprovalRepository(db)
	investmentRepo := repositories.NewInvestmentRepository(db)
	disbursementRepo := repositories.NewDisbursementRepository(db)
	notificationService := services.NewNotificationService(investmentRepo, loanRepo)

	loanService := services.LoanService{
		LoanRepo:            loanRepo,
		ApprovalRepo:        approvalRepo,
		InvestmentRepo:      investmentRepo,
		DisbursementRepo:    disbursementRepo,
		NotificationService: notificationService,
	}
	loanController := routes.LoanController{LoanService: &loanService}

	// Define routes
	e.POST("/loans", loanController.CreateLoan)
	e.GET("/approved_loans", loanController.GetApprovedLoans)
	e.PUT("/loans/:id/approve", loanController.ApproveLoan)
	e.PUT("/loans/:id/invest", loanController.InvestLoan)
	e.PUT("/loans/:id/disburse", loanController.DisburseLoan)

	// Start server
	e.Start(":8080")
}
