package job_email

import (
	"fmt"
	"strings"
	"time"
)

// JobEmail representa un registro de email enviado asociado a un trabajo.
type JobEmail struct {
	ID         int64      `json:"id"`
	JobID      int64      `json:"jobId"`
	Recipients []string   `json:"recipients"`
	Type       string     `json:"type"`
	CreatedAt  *time.Time `json:"createdAt,omitempty"`
	UpdatedAt  *time.Time `json:"updatedAt,omitempty"`
}

func (j *JobEmail) ValidateCreate() error {
	if j.JobID == 0 {
		return fmt.Errorf("job_id is required")
	}
	if len(j.Recipients) == 0 {
		return fmt.Errorf("recipients is required")
	}
	if strings.TrimSpace(j.Type) == "" {
		return fmt.Errorf("type is required")
	}
	return nil
}
