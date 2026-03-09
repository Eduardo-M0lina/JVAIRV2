package job_activity_log

import (
	"fmt"
	"time"
)

type JobActivityLog struct {
	ID        int64      `json:"id"`
	JobID     int64      `json:"jobId"`
	Type      string     `json:"type"`
	Log       string     `json:"log"`
	UserID    int64      `json:"userId"`
	CreatedAt *time.Time `json:"createdAt,omitempty"`
	UpdatedAt *time.Time `json:"updatedAt,omitempty"`
}

func (j *JobActivityLog) ValidateCreate() error {
	if j.JobID == 0 {
		return fmt.Errorf("job_id is required")
	}
	if j.UserID == 0 {
		return fmt.Errorf("user_id is required")
	}
	if j.Type == "" {
		return fmt.Errorf("type is required")
	}
	if j.Log == "" {
		return fmt.Errorf("log is required")
	}
	return nil
}
