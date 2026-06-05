# Phase 00.1 — Observability & Query Workload Enhancements

## Overview

This document describes the first major enhancement phase of the Database Scaling Lab.

## Completed Changes

### 1. Additional Database Indexes
- Added `idx_orders_created_at`
- Added `idx_events_user_id`
- Verified index usage via `EXPLAIN ANALYZE`
- Confirmed `Index Scan using idx_orders_user_id`

### 2. pg_stat_statements Integration
- Enabled `shared_preload_libraries=pg_stat_statements`
- Installed extension with:
  ```sql
  CREATE EXTENSION IF NOT EXISTS pg_stat_statements;
  ```
- Verified query capture and execution statistics

### 3. Prometheus Database Metrics
Added:
- `db_query_duration_seconds`
- `db_pool_acquired_connections`
- `db_pool_idle_connections`

Instrumented:
- get_user
- get_order
- get_user_orders
- get_user_order_summary

### 4. Context Timeouts
Implemented:
```go
ctx, cancel := context.WithTimeout(
    r.Context(),
    2*time.Second,
)
defer cancel()
```

Applied to:
- GET /users/{id}
- POST /events
- GET /health
- Future handlers

### 5. New Query Workloads

Added endpoints:

#### GET /orders/{id}
Fetch a single order.

#### GET /user-orders/{userId}
Fetch all orders for a user.

### 6. View-Based Aggregation

Created:

```sql
CREATE VIEW user_order_summary AS
SELECT
    u.id,
    COUNT(o.id)::BIGINT AS total_orders
FROM users u
LEFT JOIN orders o
    ON u.id = o.user_id
GROUP BY u.id;
```

Added:

```http
GET /user-summary/{userId}
```

## Validation Results

Dataset:

- Users: 100,000
- Orders: 1,000,000
- Events: 5,466,441

Verified:

- User endpoint working
- Orders endpoint working
- User orders endpoint working
- User summary endpoint working
- Prometheus metrics exposed
- Pool metrics exposed
- View aggregation correct
- pg_stat_statements operational
- Index scans confirmed

## Outcome

The project now has:

- Database observability
- Query-level performance tracking
- Pool monitoring
- Timeout protection
- Aggregation benchmarking capability
- Realistic read workloads

This establishes the foundation for:

1. Query optimization
2. Composite indexes
3. Materialized views
4. Read replicas
5. Partitioning
6. Caching experiments
