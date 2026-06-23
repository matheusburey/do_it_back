package task

import (
	"context"
	"errors"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type TasksRepository interface {
	GetTasks(ctx context.Context, user_id uuid.UUID) ([]Task, error)
	Create(ctx context.Context, task Task) (*Task, error)
	FindByID(ctx context.Context, id, user_id uuid.UUID) (*Task, error)
	Update(ctx context.Context, id, user_id uuid.UUID, task *Task) (*Task, error)
	Delete(ctx context.Context, id, user_id uuid.UUID) error
}

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewRepository(
	db *pgxpool.Pool,
) TasksRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) GetTasks(
	ctx context.Context, user_id uuid.UUID,
) ([]Task, error) {
	query := `
		SELECT id, user_id, title, description, completed, created_at, updated_at
		FROM tasks WHERE user_id = $1;
	`

	rows, err := r.db.Query(ctx, query, user_id)
	if err != nil {
		return nil, err
	}

	tasks, err := pgx.CollectRows(rows, pgx.RowToStructByPos[Task])
	if err != nil {
		return nil, err
	}

	return tasks, nil
}
func (r *PostgresRepository) Create(ctx context.Context, task Task) (*Task, error) {
	query := `
        INSERT INTO tasks (
			user_id,
			title,
			description,
			completed,
            created_at,
            updated_at
        )
        VALUES ($1, $2, $3, $4, NOW(), NOW())
        RETURNING id, user_id, title, description, completed, created_at, updated_at;
    `

	err := r.db.QueryRow(
		ctx,
		query,
		task.UserID,
		task.Title,
		task.Description,
		task.Completed,
	).Scan(
		&task.ID,
		&task.UserID,
		&task.Title,
		&task.Description,
		&task.Completed,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	return &task, err
}

func (r *PostgresRepository) FindByID(
	ctx context.Context,
	id, user_id uuid.UUID,
) (*Task, error) {
	task := Task{}
	query := `SELECT id, user_id, title, description, completed, created_at, updated_at FROM tasks WHERE id = $1 and user_id = $2;`

	err := r.db.QueryRow(
		ctx, query, id, user_id,
	).Scan(
		&task.ID,
		&task.UserID,
		&task.Title,
		&task.Description,
		&task.Completed,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	return &task, err
}

func (r *PostgresRepository) Update(
	ctx context.Context,
	id, user_id uuid.UUID,
	task *Task,
) (*Task, error) {
	query := `
        UPDATE tasks SET 
			title=$1, description=$2, completed=$3,updated_at=NOW()
		WHERE id = $4 AND user_id = $5
        RETURNING id, user_id, title, description, completed, created_at, updated_at;
    `

	err := r.db.QueryRow(
		ctx,
		query,
		task.Title,
		task.Description,
		task.Completed,
		id,
		user_id,
	).Scan(
		&task.ID,
		&task.UserID,
		&task.Title,
		&task.Description,
		&task.Completed,
		&task.CreatedAt,
		&task.UpdatedAt,
	)

	return task, err
}

func (r *PostgresRepository) Delete(ctx context.Context, id, user_id uuid.UUID) error {
	query := `DELETE FROM tasks WHERE id = $1 AND user_id = $2;`

	result, err := r.db.Exec(ctx, query, id, user_id)
	if err != nil {
		return err
	}

	if result.RowsAffected() == 0 {
		return errors.New("task not found")
	}

	return nil
}
