package job_sms

import (
	"fmt"
	"strings"
	"time"
)

// JobSMS representa un registro de SMS enviado asociado a un trabajo.
type JobSMS struct {
	ID         int64      `json:"id"`
	JobID      int64      `json:"jobId"`
	Recipients []string   `json:"recipients"`
	Type       string     `json:"type"`
	Message    string     `json:"message"`
	CreatedAt  *time.Time `json:"createdAt,omitempty"`
	UpdatedAt  *time.Time `json:"updatedAt,omitempty"`
}

func (j *JobSMS) ValidateCreate() error {
	if j.JobID == 0 {
		return fmt.Errorf("job_id is required")
	}
	if len(j.Recipients) == 0 {
		return fmt.Errorf("recipients is required")
	}
	if strings.TrimSpace(j.Type) == "" {
		return fmt.Errorf("type is required")
	}
	if strings.TrimSpace(j.Message) == "" {
		return fmt.Errorf("message is required")
	}
	return nil
}
