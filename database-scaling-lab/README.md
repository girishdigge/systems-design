# Database Scaling Laboratory

A backend engineering laboratory exploring PostgreSQL scaling, connection pooling, indexing strategies, replication, partitioning, and observability through measurable performance experiments.

The objective is simple:

> Build a realistic high-volume backend system, introduce one optimization at a time, benchmark the results, and document the trade-offs.

---

# Architecture

```text
                               +----------------------+
                               |      k6 Load Test    |
                               +----------+-----------+
                                          |
                                          v
                      +--------------------------------------+
                      |         Go REST API Container        |
                      |      pgxpool Connection Pool         |
                      +------------------+-------------------+
                                         |
                                         v
                      +--------------------------------------+
                      |      PostgreSQL Primary Instance     |
                      +------------------+-------------------+
                                         ^
                                         |
                +------------------------+------------------------+
                |                                                 |
                v                                                 |
     +--------------------------+                      +--------------------------+
     |    Migration Container   |                      |      Seeder Container    |
     +--------------------------+                      +--------------------------+

     +--------------------------+         Queries         +------------------------+
     |        Grafana           | <---------------------- |      Prometheus        |
     +--------------------------+                         +------------------------+
```

---

# Technology Stack

| Component     | Technology     |
| ------------- | -------------- |
| Language      | Go             |
| Database      | PostgreSQL     |
| Driver        | pgx/v5         |
| Pooling       | pgxpool        |
| Migrations    | golang-migrate |
| Load Testing  | k6             |
| Metrics       | Prometheus     |
| Dashboards    | Grafana        |
| Containers    | Docker         |
| Orchestration | Docker Compose |

---

# Current Metrics

| Metric          | Phase 0 Baseline |
| --------------- | ---------------- |
| Dataset Size    | ~8.1M Rows       |
| Throughput      | ~5,383 req/sec   |
| p95 Latency     | 4.82ms           |
| Failed Requests | 0%               |
| Database        | PostgreSQL 18    |

---

# Progress Roadmap

## Foundation

- [x] Phase 0 — Baseline Dataset, Traffic Generation & Benchmarking

### Goals

- Generate realistic SaaS-scale datasets
- Establish baseline latency and throughput
- Instrument the platform with Prometheus and Grafana
- Create reproducible benchmarks using k6

---

## Single Node Optimization

- [ ] Phase 1 — Single PostgreSQL Node Analysis

### Topics

- Query optimization
- `EXPLAIN ANALYZE`
- Indexing strategies
- Connection exhaustion
- Large joins
- N+1 query detection

### Experiments

- Remove indexes
- Create slow joins
- Saturate connections
- Measure CPU, memory and latency impact

---

## Replication

- [ ] Phase 2 — Read Replicas

### Topics

- Read/write routing

- Round-robin replica balancing

- Read-after-write consistency

- Replication lag

- Replica failure recovery

- [ ] Phase 3 — Synchronous vs Asynchronous Replication

### Topics

- Strong consistency
- Eventual consistency
- CAP tradeoffs
- Replication latency
- Data-loss scenarios

---

## Scaling Up

- [ ] Phase 4 — Vertical Scaling

### Topics

- CPU scaling
- Memory scaling
- Cost vs performance analysis
- Throughput ceilings

---

## Data Distribution

- [ ] Phase 5 — Table Partitioning

### Topics

- Range partitioning

- Hash partitioning

- List partitioning

- Composite partitioning

- [ ] Phase 6 — Database Sharding

### Topics

- Hash sharding
- Range sharding
- Geo sharding
- Directory-based sharding
- Shard routing
- Fan-out queries

---

## High Availability

- [ ] Phase 7 — Multi-Master Replication

### Topics

- Conflict resolution

- Vector clocks

- Last-write-wins strategies

- [ ] Phase 8 — Failover & Recovery

### Topics

- Primary promotion
- Health checks
- Leader election
- RTO/RPO analysis

---

## Connection Scaling

- [ ] Phase 9 — PgBouncer & Connection Pooling

### Topics

- Connection storms
- Pool sizing
- Backend connection reduction

---

## Performance Acceleration

- [ ] Phase 10 — Redis Caching Layer

### Topics

- Cache-aside
- Write-through
- Write-back
- Read-through
- Cache hit ratios

---

## Distributed Systems

- [ ] Phase 11 — Distributed Transactions

### Topics

- Two-Phase Commit (2PC)

- Saga Pattern

- Failure handling

- [ ] Phase 12 — Multi-Region Architecture

### Topics

- Cross-region replication
- Latency analysis
- Geo-distribution

---

## Observability

- [ ] Phase 13 — Advanced Observability

### Topics

- Prometheus metrics
- Grafana dashboards
- Replication lag monitoring
- Connection tracking
- P50 / P95 / P99 analysis

---

## Final Objective

By the completion of this laboratory, the platform will support:

- Read replicas
- Sync replication
- Async replication
- Partitioning
- Sharding
- Failover
- PgBouncer
- Redis caching
- Distributed transactions
- Multi-region deployments
- Advanced observability

while providing measurable benchmarks and tradeoff analysis for every architectural decision.

# Quick Start

```bash
docker compose up -d
```

Verify API:

```bash
curl http://localhost:8080/users/c1a3055f-defc-414c-9fc1-2944be267f7d
```

Run load test:

```bash
k6 run k6/mixed-load.js
```

---

# Experiment Log

| Phase                                | Description                     |
| ------------------------------------ | ------------------------------- |
| [Phase 0](docs/phase-00-baseline.md) | Baseline Dataset & Benchmarking |

Additional phases will be added as the laboratory evolves.

---

# Goals

This project exists to answer practical backend engineering questions:

- How does PostgreSQL behave at millions of rows?
- How much impact do indexes have on latency?
- When should replicas be introduced?
- What are the limits of connection pooling?
- Which optimizations produce measurable gains?
- What trade-offs emerge as systems scale?

Every phase introduces a single architectural change and measures its impact against a reproducible baseline.
