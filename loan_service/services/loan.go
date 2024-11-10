package services

import (
	"errors"
	"log"
	"time"

	"github.com/alexhendra/amartha_loan/loan_service/models"
	"github.com/alexhendra/amartha_loan/loan_service/repositories"
)

type LoanService struct {
	LoanRepo            repositories.LoanRepository
	ApprovalRepo        repositories.ApprovalRepository
	InvestmentRepo      repositories.InvestmentRepository
	DisbursementRepo    repositories.DisbursementRepository
	NotificationService NotificationService
}

func (ls *LoanService) CreateLoan(loan *models.Loan) error {
	return ls.LoanRepo.Create(loan)
}

func (ls *LoanService) GetApprovedLoans() (loans []*models.Loan, err error) {
	return ls.LoanRepo.GetApprovedLoan()
}

func (ls *LoanService) ApproveLoan(loanID uint, approval *models.Approval) error {
	loan, err := ls.LoanRepo.FindByID(loanID)
	if err != nil {
		return err
	}

	if isValidApproval := validateApprovalDate(approval.ApprovalDate); !isValidApproval {
		return errors.New("approval date must be the same as today")
	}

	if loan.State != models.StateProposed {
		return errors.New("loan is not in proposed state")
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
		return errors.New("loan is not in approved state")
	}

	totalInvestment, err := ls.InvestmentRepo.GetTotalInvestment(loanID)
	if err != nil {
		return err
	}

	if totalInvestment+investment.Amount > loan.PrincipalAmount {
		return errors.New("investment exceeds principal amount")
	}

	log.Printf("%+v", loan)

	investment.LoanID = loanID
	investment.Rate = loan.Rate
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
	if isValidApproval := validateApprovalDate(disbursement.DisbursementDate); !isValidApproval {
		return errors.New("disbursement date must be the same as today")
	}

	loan, err := ls.LoanRepo.FindByID(loanID)
	if err != nil || loan.State != models.StateInvested {
		return errors.New("loan is not in invested state")
	}

	loan.State = models.StateDisbursed
	disbursement.LoanID = loanID
	if err := ls.DisbursementRepo.Create(disbursement); err != nil {
		return err
	}
	return ls.LoanRepo.Update(loan)
}

func validateApprovalDate(approvalDate time.Time) bool {
	// Get today's date with time set to midnight
	// To be like this: 00:00:00
	today := time.Now().UTC().Truncate(24 * time.Hour)

	if approvalDate.Before(today) || approvalDate.After(today) {
		return false
	}
	return true
}
