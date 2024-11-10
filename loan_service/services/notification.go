package services

import (
	"fmt"

	"github.com/alexhendra/amartha_loan/repositories"
)

type NotificationService struct {
	InvestmentRepo repositories.InvestmentRepository
	LoanRepo       repositories.LoanRepository
}

func NewNotificationService(investmentRepo repositories.InvestmentRepository, loanRepo repositories.LoanRepository) *NotificationService {
	return &NotificationService{
		InvestmentRepo: investmentRepo,
		LoanRepo:       loanRepo,
	}
}

// SendAgreementLetters sends an email to all investors with the loan's agreement letter link.
func (ns *NotificationService) SendAgreementLetters(loanID uint) error {
	// Retrieve the loan details
	loan, err := ns.LoanRepo.FindByID(loanID)
	if err != nil {
		return fmt.Errorf("error retrieving loan: %v", err)
	}

	// Retrieve all investments for the loan
	investments, err := ns.InvestmentRepo.FindByLoanID(loanID)
	if err != nil {
		return fmt.Errorf("error retrieving investments: %v", err)
	}

	// Send an email to each investor with the link to the agreement letter
	for _, investment := range investments {
		// Simulate sending email
		err := ns.sendEmailToInvestor(investment.InvestorID, loan.AgreementLink)
		if err != nil {
			return fmt.Errorf("error sending email to investor %s: %v", investment.InvestorID, err)
		}
	}

	return nil
}

// sendEmailToInvestor simulates sending an email to an investor.
// In a real-world scenario, this would connect to an email service provider.
func (ns *NotificationService) sendEmailToInvestor(investorID string, agreementLink string) error {
	// Placeholder for sending email logic
	fmt.Printf("Sending email to investor %s with agreement link: %s\n", investorID, agreementLink)
	return nil
}
