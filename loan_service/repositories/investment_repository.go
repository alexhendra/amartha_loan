package repositories

import (
	"github.com/alexhendra/amartha_loan/models"
	"gorm.io/gorm"
)

type InvestmentRepository interface {
	Create(investment *models.Investment) error
	GetTotalInvestment(loanID uint) (float64, error)
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
	if err := r.db.Model(&models.Investment{}).Where("loan_id = ?", loanID).Select("sum(amount)").Scan(&total).Error; err != nil {
		return 0, err
	}
	return total, nil
}
