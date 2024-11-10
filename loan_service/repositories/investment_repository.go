package repositories

import (
	"github.com/alexhendra/amartha_loan/loan_service/models"
	"gorm.io/gorm"
)

type InvestmentRepository interface {
	Create(investment *models.Investment) error
	GetTotalInvestment(loanID uint) (float64, error)
	FindByLoanID(loanID uint) ([]models.Investment, error)
	GetReturnOfInvestment(investorID string) (float64, error)
}

type investmentRepository struct {
	db *gorm.DB
}

func NewInvestmentRepository(db *gorm.DB) InvestmentRepository {
	return &investmentRepository{db}
}

func (r *investmentRepository) Create(investment *models.Investment) error {
	return r.db.Create(investment).Error
}

func (r *investmentRepository) GetTotalInvestment(loanID uint) (float64, error) {
	var total float64
	if err := r.db.Model(&models.Investment{}).Where("loan_id = ?", loanID).Select("COALESCE(sum(amount), 0)").Scan(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}

func (r *investmentRepository) FindByLoanID(loanID uint) ([]models.Investment, error) {
	var investments []models.Investment
	if err := r.db.Where("loan_id = ?", loanID).Find(&investments).Error; err != nil {
		return nil, err
	}
	return investments, nil
}

func (r *investmentRepository) GetReturnOfInvestment(investorID string) (float64, error) {
	var total float64
	if err := r.db.Model(&models.Investment{}).Where("investor_id = ?", investorID).Select("COALESCE(sum(roi), 0)").Scan(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}
