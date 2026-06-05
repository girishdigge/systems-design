package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	Pool *pgxpool.Pool
}

func (r *PostgresRepository) GetUser(ctx context.Context, id uuid.UUID) (string, string, error) {
	var name, email string
	err := r.Pool.QueryRow(ctx, "SELECT name, email FROM users WHERE id = $1", id).Scan(&name, &email)
	return name, email, err
}

func (r *PostgresRepository) InsertEvent(ctx context.Context, userID uuid.UUID, eventType string) error {
	_, err := r.Pool.Exec(ctx, "INSERT INTO events (user_id, event_type, created_at) VALUES ($1, $2, NOW())", userID, eventType)
	return err
}
