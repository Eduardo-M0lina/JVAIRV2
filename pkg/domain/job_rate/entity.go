package job_rate

import (
	"fmt"
	"time"
)

type JobRate struct {
	ID              int64      `json:"id"`
	JobID           int64      `json:"jobId"`
	UserID          int64      `json:"userId"`
	JobRateStatusID int64      `json:"jobRateStatusId"`
	SalePrice       float64    `json:"salePrice"`
	RatePercent     float64    `json:"ratePercent"`
	RateFlat        float64    `json:"rateFlat"`
	TechParts       float64    `json:"techParts"`
	CompanyParts    float64    `json:"companyParts"`
	PartsReplaced   *string    `json:"partsReplaced,omitempty"`
	Deduction       float64    `json:"deduction"`
	Payment         float64    `json:"payment"`
	Paid            bool       `json:"paid"`
	Notes           *string    `json:"notes,omitempty"`
	CreatedAt       *time.Time `json:"createdAt,omitempty"`
	UpdatedAt       *time.Time `json:"updatedAt,omitempty"`
	DeletedAt       *time.Time `json:"deletedAt,omitempty"`
}

func (j *JobRate) ValidateCreate() error {
	if j.JobID == 0 {
		return fmt.Errorf("job_id is required")
	}
	if j.UserID == 0 {
		return fmt.Errorf("user_id is required")
	}
	if j.JobRateStatusID == 0 {
		return fmt.Errorf("job_rate_status_id is required")
	}
	return nil
}

func (j *JobRate) ValidateUpdate() error {
	if j.ID == 0 {
		return fmt.Errorf("id is required")
	}
	if j.UserID == 0 {
		return fmt.Errorf("user_id is required")
	}
	if j.JobRateStatusID == 0 {
		return fmt.Errorf("job_rate_status_id is required")
	}
	return nil
}

func (j *JobRate) IsDeleted() bool {
	return j.DeletedAt != nil
}

// CalculatePayment calculates the payment based on the formula:
// ((salePrice - techParts - companyParts) * (ratePercent/100)) + rateFlat - deduction
func CalculatePayment(salePrice, ratePercent, rateFlat, techParts, companyParts, deduction float64) float64 {
	return ((salePrice - techParts - companyParts) * (ratePercent / 100.0)) + rateFlat - deduction
}
