package job_task

import (
	"context"

	"github.com/your-org/jvairv2/pkg/domain/job_task"
)

func (r *repository) Update(ctx context.Context, task *job_task.JobTask) error {
	query := `
		UPDATE job_tasks
		SET user_id = ?, due_date = ?, task = ?, task_status_id = ?, updated_at = NOW()
		WHERE id = ? AND deleted_at IS NULL
	`

	_, err := r.db.ExecContext(ctx, query,
		task.UserID,
		task.DueDate,
		task.Task,
		task.TaskStatusID,
		task.ID,
	)

	return err
}
