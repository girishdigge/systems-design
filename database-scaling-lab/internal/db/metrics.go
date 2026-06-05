package db

import (
	"time"

	"database-scaling-lab/internal/metrics"

	"github.com/jackc/pgx/v5/pgxpool"
)

func StartPoolMetrics(pool *pgxpool.Pool) {
	go func() {
		ticker := time.NewTicker(5 * time.Second)

		for range ticker.C {
			stats := pool.Stat()

			metrics.DBPoolAcquired.Set(
				float64(stats.AcquiredConns()),
			)

			metrics.DBPoolIdle.Set(
				float64(stats.IdleConns()),
			)
		}
	}()
}
