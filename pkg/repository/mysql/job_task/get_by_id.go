package job_task

import (
	"context"
	"database/sql"

	"github.com/your-org/jvairv2/pkg/domain/job_task"
)

func (r *repository) GetByID(ctx context.Context, id int64) (*job_task.JobTask, error) {
	query := `
		SELECT id, job_id, user_id, due_date, task, task_status_id, created_at, updated_at, deleted_at
		FROM job_tasks
		WHERE id = ?
	`

	task := &job_task.JobTask{}
	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&task.ID,
		&task.JobID,
		&task.UserID,
		&task.DueDate,
		&task.Task,
		&task.TaskStatusID,
		&task.CreatedAt,
		&task.UpdatedAt,
		&task.DeletedAt,
	)

	if err == sql.ErrNoRows {
		return nil, job_task.ErrNotFound
	}
	if err != nil {
		return nil, err
	}

	return task, nil
}
