package job_visit

import (
	"fmt"
	"time"
)

// JobVisit representa una visita de trabajo en el sistema
type JobVisit struct {
	ID         int64      `json:"id"`
	JobID      int64      `json:"jobId"`
	UserID     int64      `json:"userId"`
	ViewableBy *string    `json:"viewableBy,omitempty"`
	Date       time.Time  `json:"date"`
	Report     *string    `json:"report,omitempty"`
	CreatedAt  *time.Time `json:"createdAt,omitempty"`
	UpdatedAt  *time.Time `json:"updatedAt,omitempty"`
	DeletedAt  *time.Time `json:"deletedAt,omitempty"`
}

// ValidateCreate valida los campos requeridos para crear una visita
func (jv *JobVisit) ValidateCreate() error {
	if jv.JobID == 0 {
		return fmt.Errorf("job_id is required")
	}
	if jv.UserID == 0 {
		return fmt.Errorf("user_id is required")
	}
	return nil
}

// ValidateUpdate valida los campos para actualizar una visita
func (jv *JobVisit) ValidateUpdate() error {
	if jv.ID == 0 {
		return fmt.Errorf("id is required")
	}
	if jv.UserID == 0 {
		return fmt.Errorf("user_id is required")
	}
	return nil
}

// IsDeleted verifica si la visita está eliminada
func (jv *JobVisit) IsDeleted() bool {
	return jv.DeletedAt != nil
}
