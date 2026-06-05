package repository

import (
	"context"

	"database-scaling-lab/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	Pool *pgxpool.Pool
}

func (r *PostgresRepository) GetUser(ctx context.Context, id uuid.UUID) (*models.User, error) {
	user := &models.User{}
	err := r.Pool.QueryRow(ctx, "SELECT id, name, email, created_at FROM users WHERE id = $1", id).Scan(&user.ID,
		&user.Name,
		&user.Email,
		&user.CreatedAt)
	return user, err
}

func (r *PostgresRepository) InsertEvent(ctx context.Context, userID uuid.UUID, eventType string) error {
	_, err := r.Pool.Exec(ctx, "INSERT INTO events (user_id, event_type, created_at) VALUES ($1, $2, NOW())", userID, eventType)
	return err
}
