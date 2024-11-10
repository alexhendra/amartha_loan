package services

import (
	"errors"

	"github.com/alexhendra/amartha_loan/loan_service/models"
	"github.com/alexhendra/amartha_loan/loan_service/repositories"
)

type LoanService struct {
	LoanRepo            repositories.LoanRepository
	ApprovalRepo        repositories.ApprovalRepository
	InvestmentRepo      repositories.InvestmentRepository
	DisbursementRepo    repositories.DisbursementRepository
	NotificationService *NotificationService
}

func (ls *LoanService) CreateLoan(loan *models.Loan) error {
	return ls.LoanRepo.Create(loan)
}

func (ls *LoanService) ApproveLoan(loanID uint, approval *models.Approval) error {
	loan, err := ls.LoanRepo.FindByID(loanID)
	if err != nil {
		return err
	}
	if loan.State != models.StateProposed {
		return errors.New("Loan is not in proposed state")
	}
	loan.State = models.StateApproved
	approval.LoanID = loanID
	if err := ls.ApprovalRepo.Create(approval); err != nil {
		return err
	}
	return ls.LoanRepo.Update(loan)
}

func (ls *LoanService) InvestLoan(loanID uint, investment *models.Investment) error {
	loan, err := ls.LoanRepo.FindByID(loanID)
	if err != nil || loan.State != models.StateApproved {
		return errors.New("Loan is not in approved state")
	}

	totalInvestment, err := ls.InvestmentRepo.GetTotalInvestment(loanID)
	if err != nil {
		return err
	}

	if totalInvestment+investment.Amount > loan.PrincipalAmount {
		return errors.New("Investment exceeds principal amount")
	}

	investment.LoanID = loanID
	if err := ls.InvestmentRepo.Create(investment); err != nil {
		return err
	}

	if totalInvestment+investment.Amount == loan.PrincipalAmount {
		loan.State = models.StateInvested
		err = ls.NotificationService.SendAgreementLetters(loanID)
		if err != nil {
			return err
		}
	}

	return ls.LoanRepo.Update(loan)
}

func (ls *LoanService) DisburseLoan(loanID uint, disbursement *models.Disbursement) error {
	loan, err := ls.LoanRepo.FindByID(loanID)
	if err != nil || loan.State != models.StateInvested {
		return errors.New("Loan is not in invested state")
	}
	loan.State = models.StateDisbursed
	disbursement.LoanID = loanID
	if err := ls.DisbursementRepo.Create(disbursement); err != nil {
		return err
	}
	return ls.LoanRepo.Update(loan)
}
