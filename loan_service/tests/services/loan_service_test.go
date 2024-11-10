package services

import (
	"errors"
	"log"
	"testing"
	"time"

	"github.com/alexhendra/amartha_loan/loan_service/models"
	"github.com/alexhendra/amartha_loan/loan_service/services"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
)

// Mock repositories
type MockLoanRepo struct {
	mock.Mock
}

func (m *MockLoanRepo) Create(loan *models.Loan) error {
	args := m.Called(loan)
	return args.Error(0)
}

func (m *MockLoanRepo) GetApprovedLoan() ([]*models.Loan, error) {
	args := m.Called()
	return args.Get(0).([]*models.Loan), args.Error(1)
}

func (m *MockLoanRepo) FindByID(id uint) (*models.Loan, error) {
	args := m.Called(id)
	return args.Get(0).(*models.Loan), args.Error(1)
}

func (m *MockLoanRepo) Update(loan *models.Loan) error {
	args := m.Called(loan)
	return args.Error(0)
}

type MockApprovalRepo struct {
	mock.Mock
}

func (m *MockApprovalRepo) Create(approval *models.Approval) error {
	args := m.Called(approval)
	return args.Error(0)
}

type MockInvestmentRepo struct {
	mock.Mock
}

func (m *MockInvestmentRepo) Create(investment *models.Investment) error {
	args := m.Called(investment)
	return args.Error(0)
}

func (m *MockInvestmentRepo) GetTotalInvestment(loanID uint) (float64, error) {
	args := m.Called(loanID)
	return 100000, args.Error(1)
}

func (m *MockInvestmentRepo) FindByLoanID(loanID uint) ([]models.Investment, error) {
	args := m.Called(loanID)
	return args.Get(0).([]models.Investment), args.Error(1)
}

type MockNotificationService struct {
	mock.Mock
}

func (m *MockNotificationService) SendAgreementLetters(loanID uint) error {
	args := m.Called(loanID)
	return args.Error(0)
}

type MockDisbursementRepo struct {
	mock.Mock
}

func (m *MockDisbursementRepo) Create(disbursement *models.Disbursement) error {
	args := m.Called(disbursement)
	return args.Error(0)
}

// Create Loan
func TestCreateLoan_Positive(t *testing.T) {
	mockLoanRepo := new(MockLoanRepo)
	loanService := &services.LoanService{
		LoanRepo: mockLoanRepo,
	}

	loan := &models.Loan{
		BorrowerID:      1,
		PrincipalAmount: 150000,
		Rate:            10,
		AgreementLink:   "https://",
	}
	loan.ID = 1

	mockLoanRepo.On("Create", loan).Return(nil)

	err := loanService.CreateLoan(loan)
	assert.Nil(t, err)
	mockLoanRepo.AssertExpectations(t)
}

func TestCreateLoan_Negative(t *testing.T) {
	mockLoanRepo := new(MockLoanRepo)
	loanService := &services.LoanService{
		LoanRepo: mockLoanRepo,
	}

	loan := &models.Loan{
		BorrowerID:      1,
		PrincipalAmount: 150000,
		Rate:            10,
		AgreementLink:   "https://",
	}
	loan.ID = 1

	mockLoanRepo.On("Create", loan).Return(errors.New("error creating loan"))
	err := loanService.CreateLoan(loan)
	assert.NotNil(t, err)
	assert.Equal(t, "error creating loan", err.Error())
	mockLoanRepo.AssertExpectations(t)
}

// Approval Loans
func TestApproveLoan_Positive(t *testing.T) {
	mockLoanRepo := new(MockLoanRepo)
	mockApprovalRepo := new(MockApprovalRepo)
	loanService := &services.LoanService{
		LoanRepo:     mockLoanRepo,
		ApprovalRepo: mockApprovalRepo,
	}

	loan := &models.Loan{
		State:           models.StateProposed,
		BorrowerID:      1,
		PrincipalAmount: 150000,
		Rate:            10,
		AgreementLink:   "https://",
	}
	loan.ID = 1

	approval := &models.Approval{
		LoanID:          1,
		PictureProofURL: "https://",
		EmployeeID:      "1234567",
		ApprovalDate:    time.Now().UTC(),
	}

	mockLoanRepo.On("FindByID", uint(1)).Return(loan, nil)
	mockApprovalRepo.On("Create", approval).Return(nil)
	mockLoanRepo.On("Update", loan).Return(nil)

	err := loanService.ApproveLoan(1, approval)
	assert.Nil(t, err)
	assert.Equal(t, models.StateApproved, loan.State)
	mockLoanRepo.AssertExpectations(t)
	mockApprovalRepo.AssertExpectations(t)
}

func TestApproveLoan_Negative(t *testing.T) {
	mockLoanRepo := new(MockLoanRepo)
	mockApprovalRepo := new(MockApprovalRepo)
	loanService := &services.LoanService{
		LoanRepo:     mockLoanRepo,
		ApprovalRepo: mockApprovalRepo,
	}

	// Loan is not in proposed state
	loan := &models.Loan{
		State:           models.StateApproved,
		BorrowerID:      1,
		PrincipalAmount: 150000,
		Rate:            10,
		AgreementLink:   "https://",
	}
	loan.ID = 1

	approval := &models.Approval{
		LoanID:          1,
		PictureProofURL: "https://",
		EmployeeID:      "1234567",
		ApprovalDate:    time.Now().UTC(),
	}

	mockLoanRepo.On("FindByID", uint(1)).Return(loan, nil)

	err := loanService.ApproveLoan(1, approval)
	assert.NotNil(t, err)
	assert.Equal(t, "loan is not in proposed state", err.Error())
	mockLoanRepo.AssertExpectations(t)
	mockApprovalRepo.AssertExpectations(t)
}

// Investment
func TestInvestLoan_Positive(t *testing.T) {
	mockLoanRepo := new(MockLoanRepo)
	mockInvestmentRepo := new(MockInvestmentRepo)
	mockNotificationService := &MockNotificationService{}
	loanService := &services.LoanService{
		LoanRepo:            mockLoanRepo,
		InvestmentRepo:      mockInvestmentRepo,
		NotificationService: mockNotificationService,
	}

	loan := &models.Loan{
		State:           models.StateApproved,
		PrincipalAmount: 500000,
	}
	loan.ID = 1

	investment := &models.Investment{
		Amount: 400000,
	}

	mockLoanRepo.On("FindByID", uint(1)).Return(loan, nil)
	mockInvestmentRepo.On("GetTotalInvestment", uint(1)).Return(0, nil)
	mockInvestmentRepo.On("Create", investment).Return(nil)
	mockNotificationService.On("SendAgreementLetters", uint(1)).Return(nil)
	mockLoanRepo.On("Update", loan).Return(nil)

	log.Print("ID LOAN NYA****")
	log.Print(loan.ID)

	err := loanService.InvestLoan(1, investment)
	assert.Nil(t, err)
	assert.Equal(t, models.StateInvested, loan.State)
	mockLoanRepo.AssertExpectations(t)
	mockInvestmentRepo.AssertExpectations(t)
	mockNotificationService.AssertExpectations(t)
}

func TestInvestLoan_Negative(t *testing.T) {
	mockLoanRepo := new(MockLoanRepo)
	mockInvestmentRepo := new(MockInvestmentRepo)
	loanService := &services.LoanService{
		LoanRepo:       mockLoanRepo,
		InvestmentRepo: mockInvestmentRepo,
	}

	loan := &models.Loan{
		State:           models.StateApproved,
		PrincipalAmount: 1000,
	}
	loan.ID = 1

	investment := &models.Investment{
		Amount: 1500, // Investment exceeds the principal
	}

	mockLoanRepo.On("FindByID", uint(1)).Return(loan, nil)
	mockInvestmentRepo.On("GetTotalInvestment", uint(1)).Return(0, nil)

	err := loanService.InvestLoan(1, investment)
	assert.NotNil(t, err)
	assert.Equal(t, "investment exceeds principal amount", err.Error())
	mockLoanRepo.AssertExpectations(t)
	mockInvestmentRepo.AssertExpectations(t)
}

// Disburse Loan
func TestDisburseLoan_Positive(t *testing.T) {
	mockLoanRepo := new(MockLoanRepo)
	mockDisbursementRepo := new(MockDisbursementRepo)
	loanService := &services.LoanService{
		LoanRepo:         mockLoanRepo,
		DisbursementRepo: mockDisbursementRepo,
	}

	loan := &models.Loan{
		State: models.StateInvested,
	}
	loan.ID = 1

	disbursement := &models.Disbursement{
		LoanID:             1,
		AgreementSignedURL: "https://",
		EmployeeID:         "1234567",
		DisbursementDate:   time.Now().UTC(),
	}

	mockLoanRepo.On("FindByID", uint(1)).Return(loan, nil)
	mockDisbursementRepo.On("Create", disbursement).Return(nil)
	mockLoanRepo.On("Update", loan).Return(nil)

	err := loanService.DisburseLoan(1, disbursement)
	assert.Nil(t, err)
	assert.Equal(t, models.StateDisbursed, loan.State)
	mockLoanRepo.AssertExpectations(t)
	mockDisbursementRepo.AssertExpectations(t)
}

func TestDisburseLoan_Negative(t *testing.T) {
	mockLoanRepo := new(MockLoanRepo)
	mockDisbursementRepo := new(MockDisbursementRepo)
	loanService := &services.LoanService{
		LoanRepo:         mockLoanRepo,
		DisbursementRepo: mockDisbursementRepo,
	}

	loan := &models.Loan{
		State: models.StateApproved, // Loan is not in invested state
	}
	loan.ID = 1

	disbursement := &models.Disbursement{
		LoanID:             1,
		AgreementSignedURL: "https://",
		EmployeeID:         "1234567",
		DisbursementDate:   time.Now().UTC(),
	}

	mockLoanRepo.On("FindByID", uint(1)).Return(loan, nil)

	err := loanService.DisburseLoan(1, disbursement)
	assert.NotNil(t, err)
	assert.Equal(t, "loan is not in invested state", err.Error())
	mockLoanRepo.AssertExpectations(t)
	mockDisbursementRepo.AssertExpectations(t)
}
