# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

Saint Andrew's Chapel website — a self-contained church website built with Go 1.25/Chi, Templ templates, HTMX + Alpine.js, PostgreSQL 17, and Redis 7. The project specification is described in @SPEC.md. Refer to it for feature prioritization, tech stack, URL routing structure, database schema, auth strategy, code organization, security requirements, and other general specifications to complete your tasks. Other documents of note:

- For information on deployment, CI/CD, and operations, see @DEPLOYMENT.md
- For architectural decisions and rationale, see @DECISIONS.md
- For implementation progress and what to work on next, see @PROGRESS.md

## Current Implementation Status

**Completed:** Steps 1–2 of Phase 1 (Homepage + About section). Steps 3–13 and Phases 2–3 are not started. See @PROGRESS.md for details.

**What exists now:**
- Homepage with hero, service times, upcoming events grid
- About section with 6 subpages (history, beliefs, worship, gospel, staff, sanctuary)
- Health endpoints (`/health`, `/health/ready`)
- 6 database migrations (events table, staff_members table, site_settings created then dropped)
- Seed command with 5 sample events and 5 staff members
- Preview deployment infrastructure (AWS EC2 + DuckDNS)

**Not yet implemented:** Authentication, RBAC, middleware (directory exists but is empty), tests, CI/CD workflows, image uploads, ministry pages, bulletins, announcements, member directory, event registration, forms, small groups.

## Development Commands

```bash
# Start dev environment (Docker containers with hot reload, attached to logs)
make dev-up

# Start dev environment (detached)
make dev-up-detached

# Stop dev environment
make dev-down

# Follow app container logs
make dev-logs

# Hot reload locally (watches .go and .templ files via air)
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

# Create a new migration (generates up + down SQL files)
make migrate-create name=create_users

# Run/rollback migrations (inside dev container)
make migrate-up
make migrate-down

# Seed development data (5 events, 5 staff members)
make seed

# Preview/production stack (uses compose.yml + compose.prod.yml)
make preview-up
make preview-down
make preview-logs
make preview-deploy    # git pull + rebuild + migrate

# Cleanup
make clean             # remove binary, coverage files, generated *_templ.go

# Show all targets
make help
```

## Directory Structure

```
cmd/
  server/main.go          # Entry point — Chi router, middleware, routes, graceful shutdown
  seed/main.go            # Seed command — sample events and staff members
internal/
  config/config.go        # Env var loading into typed Config struct
  database/
    database.go           # PostgreSQL (GORM) + Redis connections
    migrate.go            # golang-migrate wrapper (up/down subcommands)
  handlers/               # HTTP handlers — one file per resource
    health.go             # GET /health, GET /health/ready
    home.go               # GET /
    about.go              # GET /about/* (7 methods including redirect)
  middleware/              # (empty — placeholder for auth middleware)
  models/                 # GORM structs
    event.go              # Event (soft-delete, embeds gorm.Model)
    staff_member.go       # StaffMember (soft-delete, embeds gorm.Model)
  services/               # Business logic layer
    event.go              # EventService.GetUpcoming(limit)
    staff_member.go       # StaffMemberService.GetActive()
  utils/                  # (empty — placeholder)
migrations/               # SQL migration pairs (YYYYMMDDHHMMSS_name.up.sql/.down.sql)
templates/
  layouts/base.templ      # Root HTML shell (head, nav, main, footer)
  components/             # Reusable template components
    nav.templ             # Header + Alpine.js dropdown + responsive hamburger
    footer.templ          # Footer with church info, dynamic copyright year
    service_times.templ   # Sunday/Wednesday service time cards
    event_card.templ      # Event card with date/time formatting
    staff_card.templ      # Staff member card with photo placeholder
    page_header.templ     # Reusable banner component
  pages/                  # Full page templates
    home.templ            # Homepage (hero, service times, events)
    about_*.templ         # 6 about subpages
static/
  css/                    # Vanilla CSS design system (5 files)
  fonts/                  # Self-hosted EB Garamond + Nunito (variable, .ttf)
  js/                     # Vendored HTMX 2.0.8 + Alpine.js 3.15.8
  images/                 # (empty — placeholder)
nginx/                    # Nginx config (reverse proxy, TLS, rate limiting)
scripts/                  # init-letsencrypt.sh for preview deployment
```

## Architecture

**Request flow:** Nginx (prod only) → Chi router → Middleware chain → Handler → Service → GORM → PostgreSQL

**Key patterns:**
- Handlers in `internal/handlers/` receive HTTP requests and call services
- Services in `internal/services/` contain business logic
- Models in `internal/models/` are GORM structs — soft-delete models embed `gorm.Model`, hard-delete models use manual timestamp fields
- GORM is used for queries only — never use `AutoMigrate`; all schema changes go through SQL migration files in `migrations/`
- Templates are `.templ` files (type-safe, compiled to Go) — generated `*_templ.go` files are gitignored
- The server binary doubles as a CLI: `sachapel migrate up` and `sachapel migrate down` are subcommands handled in `cmd/server/main.go`

**Auth system (planned, not yet implemented):** Custom JWT stored in HTTP-only cookies, with Redis blacklist for revocation. CSRF tokens stored in Redis as `csrf:{jti}`, rendered in `<meta>` tags, sent by HTMX via `hx-headers` as `X-CSRF-Token`. Middleware: `RequireAuth`, `RequireAnyRole(roles...)`, `RequireAllRoles(roles...)`.

**RBAC (planned, not yet implemented):** Users can hold multiple roles simultaneously (public, member, deacon, elder, staff, musician, pastor, volunteer, admin). Permissions are additive, not hierarchical. Ministry page editing requires explicit assignment via `ministry_assignments`.

**Current routes:**

| Method | Path | Handler |
|--------|------|---------|
| GET | `/` | `homeHandler.Index` |
| GET | `/health` | `healthHandler.Liveness` |
| GET | `/health/ready` | `healthHandler.Readiness` |
| GET | `/about` | `aboutHandler.Index` (301 → `/about/history`) |
| GET | `/about/history` | `aboutHandler.History` |
| GET | `/about/beliefs` | `aboutHandler.Beliefs` |
| GET | `/about/worship` | `aboutHandler.Worship` |
| GET | `/about/gospel` | `aboutHandler.Gospel` |
| GET | `/about/staff` | `aboutHandler.Staff` |
| GET | `/about/sanctuary` | `aboutHandler.Sanctuary` |
| GET | `/static/*` | File server (static assets) |

**Middleware chain** (applied globally via Chi):
`RequestID` → `RealIP` → `Logger` → `Recoverer` → `Compress(5)`

## Docker Compose Structure

Development uses two compose files merged together: `docker compose -f compose.yml -f compose.dev.yml`. The dev override disables nginx, exposes ports (3000, 5432, 6379), mounts source code for hot reload, and runs `air` instead of the production binary. Environment variables come from `.env` (copy `.env.example` to create it).

Four compose files exist:
- `compose.yml` — base stack (nginx, app, postgres, redis, volumes)
- `compose.dev.yml` — development overrides (no nginx, exposed ports, source mount, air)
- `compose.prod.yml` — production overrides (certbot, memory limits for postgres/redis)
- `compose.staging.yml` — staging port overrides (8080, 8443)

## Database Migrations

Migrations live in `migrations/` as numbered SQL file pairs. The naming convention is `YYYYMMDDHHMMSS_description.up.sql` / `.down.sql`. Create new migrations with `make migrate-create name=description`.

**Current migrations:**

| Number | Description |
|--------|-------------|
| 20250101000001 | Create `site_settings` table |
| 20250101000002 | Create `events` table with indexes |
| 20250101000003 | Seed 14 site settings values |
| 20250101000004 | Create `staff_members` table with indexes |
| 20250101000005 | Seed 5 staff members |
| 20250101000006 | Drop `site_settings` table (inlined into templates per DECISIONS.md #15) |

## Design System

**CSS files** in `static/css/` (no Tailwind, no Node.js):
- `base.css` — `@font-face`, CSS custom properties, reset, typography
- `layout.css` — header, nav (responsive + hamburger + dropdown), footer, container
- `components.css` — hero, service cards, event cards, staff cards, buttons, page headers
- `utilities.css` — spacing, text alignment, display helpers
- `print.css` — print-friendly overrides

**Brand colors:**
- Primary: crimson `#89191C` (used for header, buttons, accents)
- Accent: warm gold `#B8860B`
- Neutrals: warm-tinted grays (not pure gray)

**Fonts:** Self-hosted variable fonts — EB Garamond (serif, headings) and Nunito (sans-serif, body). Loaded via `@font-face` in `base.css`.

## Conventions

- **Binary name:** `sachapel`
- **Go module:** `github.com/sfdeloach/churchsite`
- **Go version:** 1.25 (see `go.mod`)
- **Compose v2:** Use `docker compose` (no hyphen), compose files named `compose.yml` (not `docker-compose.yml`)
- **CSS:** Vanilla CSS with a design system — no Tailwind, no Node.js tooling
- **Images:** Planned — all uploads will be converted to WebP (quality 85) via `disintegration/imaging` (not yet implemented)
- **Logging:** Go `slog` — JSON output in production, text output in development (see `cmd/server/main.go`)
- **Password hashing:** bcrypt with cost factor 12 (planned)
- **Token hashing:** Verification and reset tokens hashed with SHA-256 before storage (planned)
- **Static content:** Church name, address, service times, and other rarely-changing text are inlined directly in Templ templates (not stored in DB — see DECISIONS.md #15)

## Git Commit Convention

```
<type>[optional scope]: <description>
```

Types: `feat`, `fix`, `docs`, `style`, `refactor`, `perf`, `test`, `build`, `ci`, `chore`, `deps`, `revert`

Description should be imperative ("add", "update", "remove"), under 50 characters when possible. See @DECISIONS.md for full details.

## Key Dependencies

| Package | Purpose |
|---------|---------|
| `github.com/go-chi/chi/v5` | HTTP router (built on net/http) |
| `github.com/a-h/templ` | Type-safe Go templates |
| `gorm.io/gorm` + `gorm.io/driver/postgres` | ORM for PostgreSQL |
| `github.com/golang-migrate/migrate/v4` | SQL migration runner |
| `github.com/redis/go-redis/v9` | Redis client |
