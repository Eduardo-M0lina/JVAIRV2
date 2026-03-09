package job_task

import (
	"fmt"
	"time"
)

type JobTask struct {
	ID           int64      `json:"id"`
	JobID        int64      `json:"jobId"`
	UserID       int64      `json:"userId"`
	DueDate      *time.Time `json:"dueDate,omitempty"`
	Task         string     `json:"task"`
	TaskStatusID int64      `json:"taskStatusId"`
	CreatedAt    *time.Time `json:"createdAt,omitempty"`
	UpdatedAt    *time.Time `json:"updatedAt,omitempty"`
	DeletedAt    *time.Time `json:"deletedAt,omitempty"`
}

func (j *JobTask) ValidateCreate() error {
	if j.JobID == 0 {
		return fmt.Errorf("job_id is required")
	}
	if j.UserID == 0 {
		return fmt.Errorf("user_id is required")
	}
	if j.Task == "" {
		return fmt.Errorf("task is required")
	}
	if j.TaskStatusID == 0 {
		return fmt.Errorf("task_status_id is required")
	}
	return nil
}

func (j *JobTask) ValidateUpdate() error {
	if j.ID == 0 {
		return fmt.Errorf("id is required")
	}
	if j.UserID == 0 {
		return fmt.Errorf("user_id is required")
	}
	if j.Task == "" {
		return fmt.Errorf("task is required")
	}
	if j.TaskStatusID == 0 {
		return fmt.Errorf("task_status_id is required")
	}
	return nil
}

func (j *JobTask) IsDeleted() bool {
	return j.DeletedAt != nil
}
