package job_rate_status

import (
	"fmt"
	"time"
)

type JobRateStatus struct {
	ID        int64      `json:"id"`
	Label     string     `json:"label"`
	Class     *string    `json:"class,omitempty"`
	Order     int        `json:"order"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
	DeletedAt *time.Time `json:"deletedAt,omitempty"`
}

func (j *JobRateStatus) ValidateCreate() error {
	if j.Label == "" {
		return fmt.Errorf("label is required")
	}
	return nil
}

func (j *JobRateStatus) ValidateUpdate() error {
	if j.ID == 0 {
		return fmt.Errorf("id is required")
	}
	if j.Label == "" {
		return fmt.Errorf("label is required")
	}
	return nil
}

func (j *JobRateStatus) IsDeleted() bool {
	return j.DeletedAt != nil
}
