package repository

import (
	"context"
	"time"

	"database-scaling-lab/internal/metrics"
	"database-scaling-lab/internal/models"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresRepository struct {
	Pool *pgxpool.Pool
}

func (r *PostgresRepository) GetUser(ctx context.Context, id uuid.UUID) (*models.User, error) {
	start := time.Now()
	defer func() {
		metrics.DBQueryDuration.
			WithLabelValues(metrics.QueryGetUser).
			Observe(time.Since(start).Seconds())
	}()

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

func (r *PostgresRepository) GetOrder(ctx context.Context, id uuid.UUID) (*models.Order, error) {
	start := time.Now()
	defer func() {
		metrics.DBQueryDuration.WithLabelValues(metrics.QueryGetOrder).Observe(time.Since(start).Seconds())
	}()

	order := &models.Order{}
	err := r.Pool.QueryRow(ctx, `SELECT id, user_id, total, created_at FROM orders WHERE id = $1`, id).Scan(
		&order.ID,
		&order.UserID,
		&order.Total,
		&order.CreatedAt,
	)

	return order, err
}

func (r *PostgresRepository) GetUserOrders(ctx context.Context, userID uuid.UUID) ([]models.Order, error) {

	start := time.Now()
	defer func() {
		metrics.DBQueryDuration.
			WithLabelValues(metrics.QueryGetUserOrders).
			Observe(time.Since(start).Seconds())
	}()

	rows, err := r.Pool.Query(ctx, `SELECT id, user_id, total, created_at FROM orders WHERE user_id = $1 ORDER BY created_at DESC`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var orders []models.Order

	for rows.Next() {
		var order models.Order

		err := rows.Scan(
			&order.ID,
			&order.UserID,
			&order.Total,
			&order.CreatedAt,
		)
		if err != nil {
			return nil, err
		}

		orders = append(orders, order)
	}

	return orders, rows.Err()
}

func (r *PostgresRepository) GetUserOrderSummary(ctx context.Context, userID uuid.UUID) (*models.UserOrderSummary, error) {

	start := time.Now()
	defer func() {
		metrics.DBQueryDuration.
			WithLabelValues(metrics.QueryGetUserOrderSummary).
			Observe(time.Since(start).Seconds())
	}()

	summary := &models.UserOrderSummary{}

	err := r.Pool.QueryRow(ctx, `SELECT id, total_orders FROM user_order_summary WHERE id = $1`, userID).Scan(
		&summary.UserID,
		&summary.TotalOrders,
	)

	return summary, err
}
