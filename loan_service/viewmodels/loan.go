package viewmodels

import "time"

type LoanPayload struct {
	BorrowerID      uint    `json:"borrower_id"`
	PrincipalAmount float64 `json:"principal_amount"`
	Rate            float64 `json:"rate"`
}

type ApprovalPayload struct {
	PictureProofURL string    `json:"picture_proof_url"`
	EmployeeID      string    `json:"employee_id"`
	ApprovalDate    time.Time `json:"approval_date"`
}

type InvestmentPayload struct {
	Amount     float64 `json:"amount"`
	InvestorID string  `json:"investor_id"`
}

type DisbursementPayload struct {
	AgreementSignedURL string    `json:"agreement_signed_url"`
	EmployeeID         string    `json:"employee_id"`
	DisbursementDate   time.Time `json:"disbursement_date"`
}

type MessageResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}
