package job_task

import (
	"context"

	"github.com/angumol/jvairv2/pkg/domain/job_task"
)

func (r *repository) Create(ctx context.Context, task *job_task.JobTask) error {
	query := `
		INSERT INTO job_tasks (job_id, user_id, due_date, task, task_status_id, created_at, updated_at)
		VALUES (?, ?, ?, ?, ?, NOW(), NOW())
	`

	result, err := r.db.ExecContext(ctx, query,
		task.JobID,
		task.UserID,
		task.DueDate,
		task.Task,
		task.TaskStatusID,
	)
	if err != nil {
		return err
	}

	id, err := result.LastInsertId()
	if err != nil {
		return err
	}

	task.ID = id
	return nil
}
