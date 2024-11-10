package repositories

import (
	"github.com/alexhendra/amartha_loan/models"
	"gorm.io/gorm"
)

type LoanRepository interface {
	Create(loan *models.Loan) error
	FindByID(id uint) (*models.Loan, error)
	Update(loan *models.Loan) error
}

type loanRepository struct {
	db *gorm.DB
}

func NewLoanRepository(db *gorm.DB) LoanRepository {
	return &loanRepository{db}
}

func (r *loanRepository) Create(loan *models.Loan) error {
	return r.db.Create(loan).Error
}

func (r *loanRepository) FindByID(id uint) (*models.Loan, error) {
	var loan models.Loan
	if err := r.db.First(&loan, id).Error; err != nil {
		return nil, err
	}
	return &loan, nil
}

func (r *loanRepository) Update(loan *models.Loan) error {
	return r.db.Save(loan).Error
}
