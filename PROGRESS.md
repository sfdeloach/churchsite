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

### Step 2: About Section — COMPLETE

**Database:**
- `staff_members` table (soft-delete) with `name`, `title`, `bio`, `email`, `phone`, `photo_url`, `display_order`, `is_active`, `category`, nullable `user_id` (FK deferred to Step 7)
- Migrations: `20250101000004` (create table), `20250101000005` (seed 5 staff members), `20250101000007` (add `category VARCHAR(50)` column with index, backfill pastors)

**Backend:**
- Model: `StaffMember` (`internal/models/staff_member.go`) — soft-delete, embeds `gorm.Model`, `StaffCategory` typed string with constants (`CategoryPastor`, `CategoryStaff`), `StaffCategories` map with display labels/ordering, `OrderedStaffCategories()` helper
- Service: `StaffMemberService.GetActive()`, `GroupByCategory()` (`internal/services/staff_member.go`)
- Handler: `AboutHandler` (`internal/handlers/about.go`) — 6 page methods + `Index` redirect, `Staff()` groups members by category

**Routes:**
- `GET /about` → 301 redirect to `/about/history`
- `GET /about/history` — church history page
- `GET /about/beliefs` — doctrine and confessional standards
- `GET /about/worship` — theology of worship, regulative principle
- `GET /about/gospel` — the gospel explained for visitors
- `GET /about/staff` — pastors and staff (data-driven from `staff_members` table)
- `GET /about/sanctuary` — sanctuary, fellowship hall, grounds

**Templates:**
- Components: `page_header.templ` (reusable banner), `staff_card.templ` (photo placeholder + info)
- Pages: 6 about page templates in `templates/pages/about_*.templ`
- Nav updated: "About" link replaced with Alpine.js dropdown submenu (6 sub-links)

**CSS:**
- `layout.css` — nav dropdown styles (desktop absolute, mobile inline accordion)
- `components.css` — page-header, about-content, content-section, staff-grid, staff-card styles
- `print.css` — page-header hidden, staff-card print-friendly

**Seed data** updated in `cmd/seed/main.go` to include 5 staff members

### Step 3: Ministry Pages — NOT STARTED

### Step 4: Events Calendar with CRUD — NOT STARTED

**Testing prerequisite:** Step 4 will introduce the first handler tests. Consider setting up test helpers (test DB, fixtures) as part of this step.

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
