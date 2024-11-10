package repositories

import (
	"github.com/alexhendra/amartha_loan/loan_service/models"
	"gorm.io/gorm"
)

type ApprovalRepository interface {
	Create(approval *models.Approval) error
}

type approvalRepository struct {
	db *gorm.DB
}

func NewApprovalRepository(db *gorm.DB) ApprovalRepository {
	return &approvalRepository{db}
}

func (r *approvalRepository) Create(approval *models.Approval) error {
	return r.db.Create(approval).Error
}
