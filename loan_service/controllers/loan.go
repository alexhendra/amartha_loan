package routes

import (
	"net/http"
	"strconv"

	"github.com/alexhendra/amartha_loan/loan_service/models"
	"github.com/alexhendra/amartha_loan/loan_service/services"
	"github.com/labstack/echo/v4"
)

type LoanController struct {
	LoanService *services.LoanService
}

func (lc *LoanController) CreateLoan(c echo.Context) error {
	var loan models.Loan
	if err := c.Bind(&loan); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}
	loan.State = models.StateProposed
	if err := lc.LoanService.CreateLoan(&loan); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, loan)
}

func (lc *LoanController) ApproveLoan(c echo.Context) error {
	loanID, _ := strconv.Atoi(c.Param("id"))
	approval := new(models.Approval)
	if err := c.Bind(approval); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}
	if err := lc.LoanService.ApproveLoan(uint(loanID), approval); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.NoContent(http.StatusOK)
}

func (lc *LoanController) InvestLoan(c echo.Context) error {
	loanID, _ := strconv.Atoi(c.Param("id"))
	investment := new(models.Investment)
	if err := c.Bind(investment); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}
	if err := lc.LoanService.InvestLoan(uint(loanID), investment); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.NoContent(http.StatusOK)
}

func (lc *LoanController) DisburseLoan(c echo.Context) error {
	loanID, _ := strconv.Atoi(c.Param("id"))
	disbursement := new(models.Disbursement)
	if err := c.Bind(disbursement); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}
	if err := lc.LoanService.DisburseLoan(uint(loanID), disbursement); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.NoContent(http.StatusOK)
}
