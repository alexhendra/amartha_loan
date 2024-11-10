package models

import (
	"time"

	"gorm.io/gorm"
)

type LoanState string

const (
	StateProposed  LoanState = "proposed"
	StateApproved  LoanState = "approved"
	StateInvested  LoanState = "invested"
	StateDisbursed LoanState = "disbursed"
)

type Loan struct {
	gorm.Model
	BorrowerID      uint
	PrincipalAmount float64
	Rate            float64
	AgreementLink   string
	State           LoanState `gorm:"default:'proposed'"`
}

type Approval struct {
	gorm.Model
	LoanID          uint
	PictureProofURL string
	EmployeeID      string
	ApprovalDate    time.Time
}

type Investment struct {
	gorm.Model
	LoanID     uint
	InvestorID string
	Amount     float64
}

type Disbursement struct {
	gorm.Model
	LoanID             uint
	AgreementSignedURL string
	EmployeeID         string
	DisbursementDate   time.Time
}
