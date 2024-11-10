package repositories

import (
	"github.com/alexhendra/amartha_loan/loan_service/models"
	"gorm.io/gorm"
)

type DisbursementRepository interface {
	Create(disbursement *models.Disbursement) error
}

type disbursementRepository struct {
	db *gorm.DB
}

func NewDisbursementRepository(db *gorm.DB) DisbursementRepository {
	return &disbursementRepository{db}
}

func (r *disbursementRepository) Create(disbursement *models.Disbursement) error {
	return r.db.Create(disbursement).Error
}
