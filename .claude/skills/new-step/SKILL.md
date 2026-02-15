---
name: new-step
description: Scaffold the next implementation step from SPEC.md. Use when starting a new feature step from the Phase 1-3 roadmap.
argument-hint: "[step-number]"
disable-model-invocation: true
---

# Scaffold a New Implementation Step

You are scaffolding Step $ARGUMENTS from the project roadmap.

## Process

### 1. Gather Context

- Read `SPEC.md` to understand the full requirements for Step $ARGUMENTS
- Read `PROGRESS.md` to understand what's already been completed and what infrastructure exists
- Read `DECISIONS.md` for any relevant architectural decisions
- Identify all database tables, routes, models, services, handlers, templates, and CSS needed

### 2. Plan the Work

Present a checklist of everything needed for this step before writing any code. The checklist should cover:

- **Migrations**: What tables/columns need to be created? Reference the schema in SPEC.md.
- **Models**: What GORM structs are needed? Soft-delete (embed `gorm.Model`) or hard-delete (manual timestamp fields)?
- **Services**: What business logic queries are needed?
- **Handlers**: What HTTP handlers and route methods are needed?
- **Routes**: What routes need to be registered in `cmd/server/main.go`? Include any middleware groups.
- **Templates**: What page templates and reusable components are needed?
- **CSS**: What new styles are needed in the design system files?
- **Seed data**: What sample data should be added to `cmd/seed/main.go`?
- **Nav updates**: Does the navigation need new links or dropdowns?

Wait for user approval before proceeding.

### 3. Implement in Order

Execute the implementation in this sequence:

1. **Migration files** in `migrations/` (always both `.up.sql` and `.down.sql`)
2. **Model structs** in `internal/models/`
3. **Service layer** in `internal/services/`
4. **Handler** in `internal/handlers/`
5. **Templ templates** in `templates/pages/` and `templates/components/`
6. **CSS additions** to existing files in `static/css/`
7. **Route registration** in `cmd/server/main.go`
8. **Seed data updates** in `cmd/seed/main.go`
9. **Nav/layout updates** in `templates/components/nav.templ` or `templates/layouts/base.templ`

### 4. Update Progress

After implementation, update `PROGRESS.md` with a detailed record of what was built, following the format used for Steps 1 and 2.

## Key Conventions

- Migration naming: `YYYYMMDDHHMMSS_description.up.sql` / `.down.sql`
- Never use GORM `AutoMigrate` — all schema changes via SQL migration files
- Soft-delete models embed `gorm.Model`; hard-delete models use manual timestamp fields
- Handlers are structs with a `New*Handler()` constructor that accepts service dependencies
- Services are structs with a `New*Service(db *gorm.DB)` constructor
- Templates use Templ (`.templ` files), not `html/template`
- CSS uses the existing vanilla design system — no Tailwind
- Brand colors: crimson `#89191C`, warm gold `#B8860B`, warm-tinted neutral grays
- Fonts: EB Garamond (headings), Nunito (body)
