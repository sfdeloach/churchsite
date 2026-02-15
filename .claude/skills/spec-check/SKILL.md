---
name: spec-check
description: Verify current implementation matches SPEC.md requirements. Use after completing a feature step to catch drift.
argument-hint: "[step-number]"
context: fork
agent: Explore
---

# Spec Compliance Check

Verify that the implementation of Step $ARGUMENTS matches the requirements defined in `SPEC.md`.

## Checks to Perform

### 1. Route Coverage

Compare routes defined in `SPEC.md` under "URL Routing Structure" against routes registered in `cmd/server/main.go`. Report:
- Routes in spec but not implemented (expected for future steps)
- Routes implemented but not in spec (possible drift)
- Routes implemented with wrong HTTP method

### 2. Database Schema

Compare the migrations in `migrations/` against table definitions in `SPEC.md` under "Database Schema". Report:
- Missing columns that the spec defines
- Extra columns not in the spec (may be intentional — flag for review)
- Type mismatches (e.g., `VARCHAR(100)` vs `VARCHAR(255)`)
- Missing indexes
- Missing constraints (`NOT NULL`, `UNIQUE`, `REFERENCES`)

### 3. Model Alignment

Compare GORM models in `internal/models/` against migration files. Report:
- Fields in migration but not in model
- Fields in model but not in migration
- Soft-delete vs hard-delete correctness (does the model embed `gorm.Model` when spec says soft-delete?)
- Missing `TableName()` method

### 4. Handler Completeness

For the step being checked, verify:
- All routes have corresponding handler methods
- Handlers call services (not raw DB queries)
- Handlers render templates (not raw HTML strings)
- Error handling follows the pattern: `slog.Error()` + `http.Error()` or error page

### 5. Template Completeness

Verify:
- All handler methods have corresponding `.templ` files
- Templates use `@layouts.Base()` wrapper
- Data-driven templates receive the correct model types
- Templates use existing components where appropriate

### 6. Progress Documentation

Verify that `PROGRESS.md` accurately reflects:
- Step marked as COMPLETE (if it is)
- All implemented artifacts listed (migrations, models, services, handlers, templates, routes)
- Any deviations from spec are documented

## Output Format

Provide a summary organized by check category:

```
## Spec Check: Step N

### Routes: [PASS/ISSUES]
- ...

### Schema: [PASS/ISSUES]
- ...

### Models: [PASS/ISSUES]
- ...

### Handlers: [PASS/ISSUES]
- ...

### Templates: [PASS/ISSUES]
- ...

### Progress Doc: [PASS/ISSUES]
- ...

### Overall: [PASS/NEEDS ATTENTION]
```
