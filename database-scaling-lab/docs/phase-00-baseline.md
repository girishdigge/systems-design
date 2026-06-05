# Phase 0 — Baseline Dataset & Benchmarking

## Objective

Establish a reproducible performance baseline before introducing any database or infrastructure optimizations.

Future phases will compare their results directly against this benchmark.

---

# Environment

| Component      | Configuration  |
| -------------- | -------------- |
| API            | Go + pgxpool   |
| Database       | PostgreSQL 18  |
| Deployment     | Docker Compose |
| Load Generator | Grafana k6     |
| Dataset Size   | ~8.1M Rows     |

---

# Dataset Design

The database models a simplified SaaS environment.

| Table       | Approximate Rows |
| ----------- | ---------------- |
| users       | 100,000          |
| products    | 10,000           |
| orders      | 1,000,000        |
| order_items | 2,000,000+       |
| events      | 5,000,000        |
| **Total**   | **~8.1M**        |

---

# Connection Strategy

The API communicates with PostgreSQL using `pgxpool`.

Baseline testing was intentionally performed without:

- PgBouncer
- Read replicas
- Partitioning
- Materialized views
- Query optimization

This establishes a clean reference point for future experiments.

---

# Benchmark Configuration

```text
Duration: 3 Minutes

Virtual Users: 100

Workload:
- Read Operations
- Write Operations

Target:
Go REST API
      ↓
PostgreSQL
```

Benchmark command:

```bash
k6 run k6/mixed-load.js
```

---

# Results

## Request Throughput

```text
http_reqs........: 812179
rate.............: 4878.20 req/sec
```

---

## Latency

```text
avg..............: 3.31ms
p90..............: 7.09ms
p95..............: 8.01ms
max..............: 133ms
```

---

## Reliability

```text
checks_succeeded.: 100%
checks_failed....: 0%

http_req_failed..: 0%
```

---

# Observations

## Primary Key Lookups Scale Well

Queries against:

```text
/users/:id
```

remain consistently fast despite millions of records.

The primary-key B-tree index provides logarithmic lookup performance and prevents full-table scans.

---

## PostgreSQL Cache Efficiency

Repeated access to a limited subset of records produced a high cache-hit ratio.

As a result, most requests avoided physical disk access and were served from memory.

---

## Connection Pool Stability

`pgxpool` successfully sustained more than 5,000 requests per second without exhausting database connections.

No connection-related errors were observed during testing.

---

## Bulk Seeding Strategy

Dataset generation used PostgreSQL's `COPY` protocol through `pgx.CopyFrom`.

This enabled efficient insertion of millions of rows while avoiding row-by-row insertion overhead.

---

# Baseline Summary

| Metric          | Value         |
| --------------- | ------------- |
| Dataset         | ~8.1M Rows    |
| Throughput      | 5,383 req/sec |
| Average Latency | 2.03ms        |
| p95 Latency     | 4.82ms        |
| Error Rate      | 0%            |

This baseline serves as the reference point for all future optimization phases.

Any improvement introduced in later phases must demonstrate measurable gains relative to these results.
