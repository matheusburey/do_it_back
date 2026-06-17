package auth

import (
	"context"
	"log/slog"

	"github.com/jackc/pgx/v5/pgxpool"
)

type UserRepository interface {
	Create(
		ctx context.Context,
		user User,
	) (*User, error)
	FindByEmail(
		ctx context.Context,
		email string,
	) (*User, error)
}

type PostgresRepository struct {
	db *pgxpool.Pool
}

func NewRepository(
	db *pgxpool.Pool,
) UserRepository {
	return &PostgresRepository{
		db: db,
	}
}

func (r *PostgresRepository) Create(
	ctx context.Context,
	user User,
) (*User, error) {
	query := `
        INSERT INTO users (
            name,
			email,
			password,
            created_at,
            updated_at
        )
        VALUES ($1, $2, $3, NOW(), NOW())
        RETURNING *;
    `

	err := r.db.QueryRow(
		ctx,
		query,
		user.Name,
		user.Email,
		user.Password,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		slog.Error("database", "error", err)
		return &user, err
	}

	return &user, nil
}

func (r *PostgresRepository) FindByEmail(
	ctx context.Context,
	email string,
) (*User, error) {
	user := User{}
	query := `SELECT * FROM users WHERE email = $1;`

	err := r.db.QueryRow(
		ctx,
		query,
		email,
	).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return &user, nil
}
