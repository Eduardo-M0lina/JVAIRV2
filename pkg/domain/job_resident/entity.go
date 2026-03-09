package job_resident

import (
	"fmt"
	"time"
)

type JobResident struct {
	ID          int64      `json:"id"`
	JobID       int64      `json:"jobId"`
	Name        string     `json:"name"`
	MobilePhone *string    `json:"mobilePhone,omitempty"`
	HomePhone   *string    `json:"homePhone,omitempty"`
	Email       *string    `json:"email,omitempty"`
	CreatedAt   *time.Time `json:"createdAt,omitempty"`
	UpdatedAt   *time.Time `json:"updatedAt,omitempty"`
	DeletedAt   *time.Time `json:"deletedAt,omitempty"`
}

func (j *JobResident) ValidateCreate() error {
	if j.JobID == 0 {
		return fmt.Errorf("job_id is required")
	}
	if j.Name == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}

func (j *JobResident) ValidateUpdate() error {
	if j.ID == 0 {
		return fmt.Errorf("id is required")
	}
	if j.Name == "" {
		return fmt.Errorf("name is required")
	}
	return nil
}

func (j *JobResident) IsDeleted() bool {
	return j.DeletedAt != nil
}
