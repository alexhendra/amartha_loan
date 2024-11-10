package routes

import (
	"net/http"
	"strconv"

	"github.com/alexhendra/amartha_loan/loan_service/models"
	"github.com/alexhendra/amartha_loan/loan_service/services"
	"github.com/alexhendra/amartha_loan/loan_service/viewmodels"
	"github.com/labstack/echo/v4"
)

type LoanController struct {
	LoanService *services.LoanService
}

func (lc *LoanController) CreateLoan(c echo.Context) error {
	var loanPayload viewmodels.LoanPayload
	if err := c.Bind(&loanPayload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	loanData := &models.Loan{
		State:           models.StateProposed,
		BorrowerID:      loanPayload.BorrowerID,
		PrincipalAmount: loanPayload.PrincipalAmount,
		Rate:            loanPayload.Rate,
	}
	if err := lc.LoanService.CreateLoan(loanData); err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}
	return c.JSON(http.StatusCreated, &viewmodels.MessageResponse{
		Message: "Loan request has been successfully created.",
		Data:    loanPayload,
	})
}

func (lc *LoanController) GetApprovedLoans(c echo.Context) error {
	loans, err := lc.LoanService.GetApprovedLoans()
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, err.Error())
	}

	return c.JSON(http.StatusOK, &viewmodels.MessageResponse{
		Message: "Get approved loans.",
		Data:    loans,
	})
}

func (lc *LoanController) ApproveLoan(c echo.Context) error {
	loanID, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	approvalPayload := new(viewmodels.ApprovalPayload)
	if err := c.Bind(approvalPayload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	approvalData := &models.Approval{
		LoanID:          uint(loanID),
		PictureProofURL: approvalPayload.PictureProofURL,
		EmployeeID:      approvalPayload.EmployeeID,
		ApprovalDate:    approvalPayload.ApprovalDate,
	}
	if err := lc.LoanService.ApproveLoan(uint(loanID), approvalData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, &viewmodels.MessageResponse{
		Message: "Loan request has been approved.",
		Data:    approvalPayload,
	})
}

func (lc *LoanController) InvestLoan(c echo.Context) error {
	loanID, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	investmentPayload := new(viewmodels.InvestmentPayload)
	if err := c.Bind(investmentPayload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	investmentData := &models.Investment{
		LoanID:     uint(loanID),
		InvestorID: investmentPayload.InvestorID,
		Amount:     investmentPayload.Amount,
	}

	if err := lc.LoanService.InvestLoan(uint(loanID), investmentData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, &viewmodels.MessageResponse{
		Message: "Your investment has been submitted.",
		Data:    investmentPayload,
	})
}

func (lc *LoanController) DisburseLoan(c echo.Context) error {
	loanID, _ := strconv.ParseUint(c.Param("id"), 10, 32)
	disbursementPayload := new(viewmodels.DisbursementPayload)
	if err := c.Bind(disbursementPayload); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, "Invalid input")
	}

	disbursementData := &models.Disbursement{
		LoanID:             uint(loanID),
		AgreementSignedURL: disbursementPayload.AgreementSignedURL,
		EmployeeID:         disbursementPayload.EmployeeID,
		DisbursementDate:   disbursementPayload.DisbursementDate,
	}

	if err := lc.LoanService.DisburseLoan(uint(loanID), disbursementData); err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, err.Error())
	}
	return c.JSON(http.StatusOK, &viewmodels.MessageResponse{
		Message: "Your loan has been successfully disbursed.",
		Data:    disbursementPayload,
	})
}
