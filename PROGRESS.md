# Progress Tracker

This document tracks implementation progress against the [SPEC.md](SPEC.md) feature list. Update it as steps are completed so that any developer (or AI assistant) can quickly orient themselves.

---

## Infrastructure

### Preview Deployment (AWS EC2) — COMPLETE

Temporary preview environment for stakeholder review. See [DEPLOYMENT.md](DEPLOYMENT.md) for setup instructions.

- Nginx config: `nginx/nginx.conf`, `nginx/sites-enabled/default.conf`
- Let's Encrypt bootstrap: `scripts/init-letsencrypt.sh`
- Production compose override: `compose.prod.yml` (certbot, memory tuning)
- Production env template: `.env.production.example`
- Dockerfile: includes seed binary for demo data
- Makefile targets: `preview-up`, `preview-down`, `preview-logs`, `preview-deploy`

### CI/CD & Deployment — NOT STARTED

**Prerequisites:**
- Provision permanent VPS (Hetzner CX21 or alternative)
- Configure DNS for sachapel.com and staging.sachapel.com
- Add SSH keys and secrets to GitHub repository settings
- Install Docker and Docker Compose on VPS

**Pipeline definitions** are in [DEPLOYMENT.md](DEPLOYMENT.md) — need to create the actual workflow files:
- `.github/workflows/test.yml` — test & lint on PRs
- `.github/workflows/deploy-staging.yml` — deploy to staging on push to develop
- `.github/workflows/deploy-production.yml` — deploy to production on push to main

---

## Phase 1 — MVP

### Step 1: Homepage with Service Times and Upcoming Events — COMPLETE

**Infrastructure built:**
- Configuration package (`internal/config/config.go`) — loads env vars into typed struct
- Database package (`internal/database/database.go`) — PostgreSQL (GORM) + Redis connections
- Migration system (`internal/database/migrate.go`) — golang-migrate wrapper with `migrate up/down` subcommands
- GORM models for `Event` (soft-delete) and `SiteSetting` (hard-delete)
- Services: `SiteSettingsService.GetAll()`, `EventService.GetUpcoming(limit)`
- Health endpoints: `GET /health` (liveness), `GET /health/ready` (readiness)
- Seed command (`cmd/seed/main.go`) — 5 sample events
- Main entry point (`cmd/server/main.go`) — Chi router, middleware (RequestID, RealIP, Logger, Recoverer, Compress), static file serving, graceful shutdown

**Migrations:**
- `20250101000001` — `site_settings` table
- `20250101000002` — `events` table with indexes
- `20250101000003` — seed 14 site settings (church info, service times, hero text)

**Templates:**
- Base layout (`templates/layouts/base.templ`) — HTML shell with head, nav, footer, children slot
- Nav (`templates/components/nav.templ`) — sticky header, hamburger menu (Alpine.js), slide animation
- Footer (`templates/components/footer.templ`) — church info, service times, dynamic copyright year
- Service times (`templates/components/service_times.templ`) — Sunday + Wednesday cards
- Event card + grid (`templates/components/event_card.templ`)
- Homepage (`templates/pages/home.templ`) — hero, service times, upcoming events

**CSS design system (5 files in `static/css/`):**
- `base.css` — reset, CSS variables, typography, `@font-face` declarations
- `layout.css` — header, nav (responsive with slide animation), footer
- `components.css` — hero, service cards, event cards, buttons
- `utilities.css` — spacing, text, display helpers
- `print.css` — print-friendly overrides

**Deviations from original plan:**
- Color palette: changed from navy blue to brand crimson (`#89191C`) with warm gold accent (`#B8860B`) and warm-tinted neutral grays
- Fonts: self-hosted EB Garamond (serif) and Nunito (sans-serif) variable fonts instead of system font stacks
- Footer copyright year: dynamic via Go `time.Now().Year()` instead of hardcoded

### Step 2: About Section — NOT STARTED

### Step 3: Ministry Pages — NOT STARTED

### Step 4: Events Calendar with CRUD — NOT STARTED

### Step 5: Bulletins — NOT STARTED

### Step 6: Announcements System — NOT STARTED

### Step 7: User Registration with Email Verification — NOT STARTED

### Step 8: Password Reset — NOT STARTED

### Step 9: Member Directory — NOT STARTED

### Step 10: Role-Based Access Control — NOT STARTED

### Step 11: Event Registration with Capacity Limits — NOT STARTED

### Step 12: JSON-Schema Based Forms — NOT STARTED

### Step 13: Small Group Directory — NOT STARTED

---

## Phase 2

Not started.

---

## Phase 3

Not started.
