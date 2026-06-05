package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

const (
	QueryGetUser             = "get_user"
	QueryInsertEvent         = "insert_event"
	QueryGetOrder            = "get_order"
	QueryGetUserOrders       = "get_user_orders"
	QueryGetUserOrderSummary = "get_user_order_summary"
)

var (
	DBQueryDuration = promauto.NewHistogramVec(
		prometheus.HistogramOpts{
			Name: "db_query_duration_seconds",
			Help: "Database query latency",
		},
		[]string{"query"},
	)

	DBPoolAcquired = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_pool_acquired_connections",
			Help: "Connections currently acquired",
		},
	)

	DBPoolIdle = promauto.NewGauge(
		prometheus.GaugeOpts{
			Name: "db_pool_idle_connections",
			Help: "Idle connections",
		},
	)
)
