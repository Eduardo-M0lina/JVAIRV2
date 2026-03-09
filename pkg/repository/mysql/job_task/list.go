package job_task

import (
	"context"

	"github.com/your-org/jvairv2/pkg/domain/job_task"
)

func (r *repository) List(ctx context.Context, jobID int64, limit, offset int) ([]*job_task.JobTask, int64, error) {
	query := `
		SELECT id, job_id, user_id, due_date, task, task_status_id, created_at, updated_at, deleted_at
		FROM job_tasks
		WHERE job_id = ? AND deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, jobID, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer func() { _ = rows.Close() }()

	var tasks []*job_task.JobTask
	for rows.Next() {
		task := &job_task.JobTask{}
		err := rows.Scan(
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
		if err != nil {
			return nil, 0, err
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	countQuery := "SELECT COUNT(*) FROM job_tasks WHERE job_id = ? AND deleted_at IS NULL"
	var total int64
	err = r.db.QueryRowContext(ctx, countQuery, jobID).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}

func (r *repository) ListAll(ctx context.Context, limit, offset int) ([]*job_task.JobTask, int64, error) {
	query := `
		SELECT id, job_id, user_id, due_date, task, task_status_id, created_at, updated_at, deleted_at
		FROM job_tasks
		WHERE deleted_at IS NULL
		ORDER BY created_at DESC
		LIMIT ? OFFSET ?
	`

	rows, err := r.db.QueryContext(ctx, query, limit, offset)
	if err != nil {
		return nil, 0, err
	}
	defer func() { _ = rows.Close() }()

	var tasks []*job_task.JobTask
	for rows.Next() {
		task := &job_task.JobTask{}
		err := rows.Scan(
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
		if err != nil {
			return nil, 0, err
		}
		tasks = append(tasks, task)
	}

	if err = rows.Err(); err != nil {
		return nil, 0, err
	}

	countQuery := "SELECT COUNT(*) FROM job_tasks WHERE deleted_at IS NULL"
	var total int64
	err = r.db.QueryRowContext(ctx, countQuery).Scan(&total)
	if err != nil {
		return nil, 0, err
	}

	return tasks, total, nil
}
