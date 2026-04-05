# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Overview
Monitoring & Observability microservice for Gondor platform. Go/Gin service handling alert rules, alerts, audit logging, and service health status tracking.

## Commands
- `make build` -- compile to bin/server
- `make run` -- run locally (needs PostgreSQL + Redis)
- `make test` -- run all tests with race detector
- `make lint` -- golangci-lint
- `make docker` -- build Docker image
- `make migrate-up` -- run database migrations
- `make migrate-down` -- rollback migrations

## Architecture
- `cmd/server/main.go` -- entry point, dependency injection, route registration
- `internal/config/` -- env-based configuration
- `internal/model/` -- GORM domain models (AlertRule, Alert, AuditLog, ServiceStatus)
- `internal/repository/` -- database access layer
- `internal/service/` -- business logic
- `internal/handler/` -- HTTP handlers (Gin)
- `internal/middleware/` -- JWT auth (validate-only), logging
- `internal/pkg/jwt/` -- JWT validation (tokens issued by gondor-users-security)

## Key Decisions
- JWT tokens are validated only (issued by gondor-users-security service)
- Port 8009
- Database: gondor_monitoring (PostgreSQL, database-per-service)
- Multi-tenancy via company_id on alert rules and audit logs
- Alert rule conditions: gt, lt, eq, gte, lte
- Alert rule severities: info, warning, critical
- Alert statuses: firing, resolved
- Service statuses: healthy, degraded, unhealthy
- POST /v1/monitoring/audit-logs and POST /v1/monitoring/services/status skip JWT auth (called by other services)
- All monitoring routes under /v1/monitoring/ prefix
- /health and /metrics skip JWT auth

## Database
PostgreSQL with GORM. Tables: alert_rules, alerts, audit_logs, service_statuses.

## Environment Variables
- `PORT` (default: 8009)
- `DATABASE_URL` (default: postgres://gondor:gondor_dev@localhost:5432/gondor_monitoring?sslmode=disable)
- `JWT_SECRET` (default: change-me-in-production)
- `REDIS_URL` (default: localhost:6379)
- `NATS_URL` (default: nats://localhost:4222)
- `LOG_LEVEL` (default: info)
- `ENVIRONMENT` (default: development)
