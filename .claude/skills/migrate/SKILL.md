---
name: migrate
description: Create and manage database migrations following project conventions. Use when adding or modifying database tables.
argument-hint: "[create name|up|down]"
disable-model-invocation: true
---

# Database Migration Helper

## Commands

- `/migrate create <name>` — Generate a new migration pair
- `/migrate up` — Run pending migrations via `make migrate-up`
- `/migrate down` — Roll back the last migration via `make migrate-down`

## Creating Migrations

When the argument starts with `create`:

### 1. Determine the Next Number

Read the existing files in `migrations/` to find the highest numbered migration. The next migration number should increment by 1 from the highest existing number. The project uses the format `YYYYMMDDHHMMSS` but the current migrations use a sequential scheme starting at `20250101000001`.

### 2. Reference the Schema

Look up the target table in `SPEC.md` under "Database Schema" for the canonical column definitions, types, constraints, and indexes.

### 3. Write Both Files

Create both files in the `migrations/` directory:

**Up migration** (`YYYYMMDDHHMMSS_<name>.up.sql`):
- Use `CREATE TABLE` for new tables, `ALTER TABLE` for modifications
- Include all indexes defined in the spec
- Include `NOT NULL`, `DEFAULT`, `UNIQUE`, and `REFERENCES` constraints
- For soft-delete tables: include `created_at`, `updated_at`, `deleted_at` columns
- For hard-delete tables: include `created_at`, `updated_at` (no `deleted_at`)
- Always use `BIGSERIAL` for primary keys on soft-delete tables, `SERIAL` for hard-delete
- Add `CREATE INDEX` statements for commonly queried columns

**Down migration** (`YYYYMMDDHHMMSS_<name>.down.sql`):
- Must fully reverse the up migration
- For new tables: `DROP TABLE IF EXISTS <name>;`
- For added columns: `ALTER TABLE <name> DROP COLUMN IF EXISTS <column>;`
- Drop indexes before dropping tables if needed

### 4. Corresponding GORM Model

After creating the migration, remind the user to create or update the corresponding GORM model in `internal/models/`. Key patterns:
- Soft-delete models embed `gorm.Model` (provides `ID`, `CreatedAt`, `UpdatedAt`, `DeletedAt`)
- Hard-delete models define `ID`, `CreatedAt`, `UpdatedAt` manually
- Always define `TableName()` method
- Use GORM struct tags: `gorm:"column:<name>;type:<type>;not null"` etc.
- Use JSON struct tags: `json:"<snake_case>"`
- Nullable fields use pointer types (`*uint`, `*time.Time`, `*string`)

## Example Migration Style

```sql
-- Up: 20250101000004_create_staff_members.up.sql
CREATE TABLE staff_members (
    id             BIGSERIAL PRIMARY KEY,
    user_id        BIGINT,
    name           VARCHAR(255) NOT NULL,
    title          VARCHAR(255) NOT NULL,
    bio            TEXT,
    email          VARCHAR(255),
    phone          VARCHAR(20),
    photo_url      VARCHAR(500),
    display_order  INTEGER DEFAULT 0,
    is_active      BOOLEAN DEFAULT TRUE,
    created_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at     TIMESTAMP
);

CREATE INDEX idx_staff_members_is_active ON staff_members(is_active);
CREATE INDEX idx_staff_members_display_order ON staff_members(display_order);
CREATE INDEX idx_staff_members_deleted_at ON staff_members(deleted_at);
```

```sql
-- Down: 20250101000004_create_staff_members.down.sql
DROP TABLE IF EXISTS staff_members;
```

## Running Migrations

When the argument is `up` or `down`, run the corresponding make target:
- `make migrate-up` — applies all pending migrations
- `make migrate-down` — rolls back the last applied migration
