# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Saint Andrew's Chapel website — a self-contained church website built with Go/Chi, Templ templates, HTMX + Alpine.js, PostgreSQL 17, and Redis 7. The project specification is described in @SPEC.md. Refer to it for feature prioritization, tech stack, URL routing structure, database schema, auth strategy, code organization, security requirements, and other general specifications to complete your tasks. Other documents of note:

- For information on deployment, CI/CD, and operations, see @DEPLOYMENT.md
- For architectural decisions and rationale, see @DECISIONS.md
- For implementation progress and what to work on next, see @PROGRESS.md

## Development Commands

```bash
# Start dev environment (Docker containers with hot reload)
make dev-up

# Stop dev environment
make dev-down

# Follow app container logs
make dev-logs

# Hot reload locally (watches .go and .templ files)
make watch

# Generate Templ templates
make generate

# Build production binary
make build

# Run all tests with race detection
make test

# Run a single test
go test -v -race ./internal/handlers/ -run TestHandlerName

# Run linter
make lint

# Create a new migration
make migrate-create name=create_users

# Run/rollback migrations
make migrate-up
make migrate-down

# Seed development data
make seed
```

## Architecture

**Request flow:** Nginx (prod only) → Chi router → Middleware chain → Handler → Service → GORM → PostgreSQL

**Key patterns:**
- Handlers in `internal/handlers/` receive HTTP requests and call services
- Services in `internal/services/` contain business logic
- Models in `internal/models/` are GORM structs — soft-delete models embed `gorm.Model`, hard-delete models use manual timestamp fields
- GORM is used for queries only — never use `AutoMigrate`; all schema changes go through SQL migration files in `migrations/`
- Templates are `.templ` files (type-safe, compiled to Go) — generated `*_templ.go` files are gitignored

**Auth system:** Custom JWT stored in HTTP-only cookies, with Redis blacklist for revocation. CSRF tokens stored in Redis as `csrf:{jti}`, rendered in `<meta>` tags, sent by HTMX via `hx-headers` as `X-CSRF-Token`. Middleware: `RequireAuth`, `RequireAnyRole(roles...)`, `RequireAllRoles(roles...)`.

**RBAC:** Users can hold multiple roles simultaneously (public, member, deacon, elder, staff, musician, pastor, volunteer, admin). Permissions are additive, not hierarchical. Ministry page editing requires explicit assignment via `ministry_assignments`.

## Docker Compose Structure

Development uses two compose files merged together: `docker compose -f compose.yml -f compose.dev.yml`. The dev override disables nginx, exposes ports (3000, 5432, 6379), mounts source code for hot reload, and runs `air` instead of the production binary. Environment variables come from `.env` (copy `.env.example` to create it).

## Conventions

- **Binary name:** `sachapel`
- **Go module:** `github.com/sfdeloach/churchsite`
- **Compose v2:** Use `docker compose` (no hyphen), compose files named `compose.yml` (not `docker-compose.yml`)
- **CSS:** Vanilla CSS with a design system — no Tailwind, no Node.js tooling
- **Images:** All uploads converted to WebP (quality 85) via `disintegration/imaging`
- **Logging:** Go `slog` with structured JSON output
- **Password hashing:** bcrypt with cost factor 12
- **Token hashing:** Verification and reset tokens hashed with SHA-256 before storage
