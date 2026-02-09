# Saint Andrew's Chapel Website - Technical Specification

**Version:** 2.1 (Development Reference)
**Last Updated:** February 7, 2026

- For information on deployment, CI/CD, and operations, see @DEPLOYMENT.md
- For architectural decisions and rationale, see @DECISIONS.md
- For implementation progress and what to work on next, see @PROGRESS.md

---

## Feature Prioritization

### MVP - Phase 1

1. Homepage with service times and upcoming events
2. About section (Who We Are, Beliefs, History, Staff, Sanctuary)
3. Ministry pages with staff-editable content
4. Events calendar with CRUD, visibility scheduling, recurrence
5. Bulletins (2-week retention, morning/evening PDFs)
6. Announcements system
7. User registration with email verification
8. Password reset with hashed tokens
9. Member directory (opt-in)
10. Role-based access control
11. Event registration with capacity limits
12. JSON-schema based forms (developer-defined)
13. Small group directory (members-only)

### Phase 2

14. Visual form builder UI (drag-and-drop)
15. Prayer request system (elders/pastors only)
16. Volunteer scheduling
17. Photo gallery
18. Self-hosted analytics (Umami)
19. Audit logging

### Phase 3

20. Self-hosted sermon video
21. Email newsletter system
22. PWA/mobile capabilities

---

## Technical Stack

| Component  | Technology                          |
| ---------- | ----------------------------------- |
| Backend    | Go 1.25+ with Chi router            |
| Templates  | Templ (type-safe, compiled)         |
| Frontend   | HTMX (v2.0.8) + Alpine.js (v3.15.8) |
| CSS        | Vanilla CSS with design system      |
| Database   | PostgreSQL 17                       |
| Cache      | Redis 7                             |
| ORM        | GORM (queries only, no AutoMigrate) |
| Migrations | golang-migrate                      |
| Auth       | JWT (golang-jwt) + bcrypt           |
| Email      | net/smtp (M365 SMTP relay)          |
| Images     | disintegration/imaging → WebP       |
| Logging    | Go slog (structured JSON)           |
| Rich Text  | Quill editor + bluemonday sanitizer |

---

## URL Routing Structure

### Public Routes

```
/                                    # Homepage
/about/history                       # History & Identity
/about/beliefs                       # Doctrine & Beliefs
/about/worship                       # Theology of Worship
/about/gospel                        # What Is the Gospel?
/about/staff                         # Pastors & Staff
/about/sanctuary                      # Our Place of Worship

/ministries                          # Ministry overview
/ministries/{slug}                   # Individual ministry page

/calendar/events                     # Public events calendar
/calendar/events/:id                 # Event details

/resources/bulletins                 # Current bulletins (2 weeks)
/resources/supported-ministries      # Supported Ministries

/new-members                         # New Member's Class info
/employment                          # Employment Opportunities
/internship                          # Internship Information
/visit                               # Visit Us / Directions
/give                                # Giving (Breeze redirect)

/login                               # Login page
/register                            # Member registration
/forgot-password                     # Password reset request
/reset-password/:token               # Password reset form
/verify-email/:token                 # Email verification
```

### Member Routes (requires authentication + member role)

```
/member/dashboard                    # Member dashboard
/member/directory                    # Member directory (opt-in)
/member/profile                      # Edit own profile
/member/budget                       # Current year operating budget
/member/bible-studies                # Small group directory
/member/events/register/:id          # Event registration forms
```

### Staff Routes (requires staff or admin role)

```
/staff/dashboard                     # Staff dashboard
/staff/announcements                 # Manage announcements
/staff/announcements/create          # Create announcement
/staff/announcements/:id/edit        # Edit announcement
/staff/bulletins                     # Manage bulletins
/staff/bulletins/upload              # Upload bulletin
/staff/events                        # Manage events
/staff/events/create                 # Create event
/staff/events/:id/edit               # Edit event
/staff/events/:id/registrations      # View registrations
/staff/forms/:id/submissions         # View form submissions
/staff/forms/:id/export              # Export as CSV
/staff/ministry/:slug                # Edit ministry page
```

### Elder/Pastor Routes

```
/elder/dashboard                     # Elder/Pastor dashboard
/elder/prayer-requests               # Prayer requests
/elder/prayer-requests/:id           # Individual request
```

### Admin Routes

```
/admin/dashboard                     # Admin dashboard
/admin/users                         # Manage users
/admin/users/:id                     # Edit user / assign roles
/admin/users/pending                 # Pending verifications
/admin/users/invitations             # Send invitations
/admin/photos                        # Photo gallery management
/admin/photos/upload                 # Upload photos
/admin/forms                         # Form builder (Phase 2)
/admin/analytics                     # Analytics dashboard
/admin/audit-log                     # Login attempt logs
/admin/settings                      # Site settings
```

### API Routes (JSON)

```
/api/v1/auth/login                   # POST: Authenticate, returns JWT
/api/v1/auth/logout                  # POST: Invalidate token
/api/v1/auth/refresh                 # POST: Refresh JWT token
/api/v1/events                       # GET: List events
/api/v1/events/:id                   # GET: Event details
/api/v1/events/:id/register          # POST: Register for event
/api/v1/bulletins/current            # GET: Current bulletins
/api/v1/forms/:id/submit             # POST: Submit form
/api/v1/member/directory/search      # GET: Search directory (auth)
```

### Health Routes

```
/health                              # 200 OK if running
/health/ready                        # 200 if PostgreSQL + Redis connected
```

---

## Database Schema

### Model Patterns

**Soft Delete (embed gorm.Model):** users, events, announcements, bulletins, ministries, forms, staff_members

**Hard Delete (manual fields):** user_roles, event_registrations, form_submissions, audit_log, site_settings, member_profiles, small_groups, photos, ministry_assignments

### Core Tables

#### users

```sql
CREATE TABLE users (
    id                  BIGSERIAL PRIMARY KEY,
    email               VARCHAR(255) UNIQUE NOT NULL,
    password_hash       VARCHAR(255) NOT NULL,
    first_name          VARCHAR(100) NOT NULL,
    last_name           VARCHAR(100) NOT NULL,
    phone               VARCHAR(20),
    is_verified         BOOLEAN DEFAULT FALSE,
    verification_token  VARCHAR(255),    -- SHA-256 hash
    reset_token         VARCHAR(255),    -- SHA-256 hash
    reset_token_expires TIMESTAMP,
    last_login          TIMESTAMP,
    failed_login_count  INTEGER DEFAULT 0,
    locked_until        TIMESTAMP,
    created_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    deleted_at          TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_deleted_at ON users(deleted_at);
```

#### roles

```sql
CREATE TABLE roles (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Seed: public, member, deacon (future), elder, staff, musician, pastor, volunteer, admin
```

#### user_roles

```sql
CREATE TABLE user_roles (
    user_id     BIGINT REFERENCES users(id) ON DELETE CASCADE,
    role_id     INTEGER REFERENCES roles(id) ON DELETE CASCADE,
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    assigned_by BIGINT REFERENCES users(id),
    PRIMARY KEY (user_id, role_id)
);
```

#### member_profiles

```sql
CREATE TABLE member_profiles (
    user_id              BIGINT PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    photo_url            VARCHAR(500),
    address_line1        VARCHAR(255),
    address_line2        VARCHAR(255),
    city                 VARCHAR(100),
    state                VARCHAR(2),
    zip_code             VARCHAR(10),
    membership_date      DATE,
    directory_opt_in     BOOLEAN DEFAULT FALSE,
    show_email           BOOLEAN DEFAULT FALSE,
    show_phone           BOOLEAN DEFAULT FALSE,
    show_address         BOOLEAN DEFAULT FALSE,
    emergency_contact    VARCHAR(255),
    emergency_phone      VARCHAR(20),
    created_at           TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at           TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Content Tables

#### announcements

```sql
CREATE TABLE announcements (
    id              BIGSERIAL PRIMARY KEY,
    title           VARCHAR(255) NOT NULL,
    content         TEXT NOT NULL,
    author_id       BIGINT REFERENCES users(id),
    visible_from    TIMESTAMP,
    visible_until   TIMESTAMP,
    is_active       BOOLEAN DEFAULT TRUE,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

#### events

```sql
CREATE TABLE events (
    id                    BIGSERIAL PRIMARY KEY,
    title                 VARCHAR(255) NOT NULL,
    description           TEXT,
    event_date            TIMESTAMP NOT NULL,
    end_date              TIMESTAMP,
    location              VARCHAR(255),
    location_details      TEXT,
    is_recurring          BOOLEAN DEFAULT FALSE,
    recurrence_rule       VARCHAR(20),  -- none, daily, weekly, biweekly, monthly, yearly
    recurrence_end        DATE,
    registration_enabled  BOOLEAN DEFAULT FALSE,
    capacity_limit        INTEGER,
    registration_deadline TIMESTAMP,
    visible_from          TIMESTAMP,
    visible_until         TIMESTAMP,
    is_public             BOOLEAN DEFAULT TRUE,
    ministry_id           INTEGER REFERENCES ministries(id),
    created_by            BIGINT REFERENCES users(id),
    created_at            TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at            TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

#### event_registrations

```sql
CREATE TABLE event_registrations (
    id               BIGSERIAL PRIMARY KEY,
    event_id         BIGINT REFERENCES events(id) ON DELETE CASCADE,
    user_id          BIGINT REFERENCES users(id),
    guest_name       VARCHAR(255),
    guest_email      VARCHAR(255),
    guest_phone      VARCHAR(20),
    number_attending INTEGER DEFAULT 1,
    special_needs    TEXT,
    form_data        JSONB,
    registered_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_registration UNIQUE(event_id, user_id)
);
```

#### bulletins

```sql
CREATE TABLE bulletins (
    id            BIGSERIAL PRIMARY KEY,
    bulletin_date DATE NOT NULL,
    service_type  VARCHAR(20) NOT NULL,  -- 'morning' or 'evening'
    file_path     VARCHAR(500) NOT NULL,
    file_size     INTEGER,
    uploaded_by   BIGINT REFERENCES users(id),
    uploaded_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_bulletin UNIQUE(bulletin_date, service_type)
);
```

#### ministries

```sql
CREATE TABLE ministries (
    id            SERIAL PRIMARY KEY,
    name          VARCHAR(255) NOT NULL,
    slug          VARCHAR(255) UNIQUE NOT NULL,
    description   TEXT,
    leader_id     BIGINT REFERENCES users(id),
    contact_email VARCHAR(255),
    meeting_time  VARCHAR(255),
    location      VARCHAR(255),
    is_active     BOOLEAN DEFAULT TRUE,
    sort_order    INTEGER DEFAULT 0,
    page_content  TEXT,
    created_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

#### ministry_assignments

```sql
CREATE TABLE ministry_assignments (
    id          SERIAL PRIMARY KEY,
    user_id     BIGINT NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    ministry_id INTEGER NOT NULL REFERENCES ministries(id) ON DELETE CASCADE,
    can_edit    BOOLEAN DEFAULT TRUE,
    assigned_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    assigned_by BIGINT REFERENCES users(id),
    CONSTRAINT unique_ministry_assignment UNIQUE(user_id, ministry_id)
);
```

#### small_groups

```sql
CREATE TABLE small_groups (
    id           SERIAL PRIMARY KEY,
    name         VARCHAR(255) NOT NULL,
    description  TEXT,
    leader_id    BIGINT REFERENCES users(id),
    meeting_day  VARCHAR(20),
    meeting_time TIME,
    location     VARCHAR(255),
    is_active    BOOLEAN DEFAULT TRUE,
    created_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Form System

#### forms

```sql
CREATE TABLE forms (
    id         BIGSERIAL PRIMARY KEY,
    title      VARCHAR(255) NOT NULL,
    description TEXT,
    schema     JSONB NOT NULL,
    is_active  BOOLEAN DEFAULT TRUE,
    created_by BIGINT REFERENCES users(id),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

#### form_submissions

```sql
CREATE TABLE form_submissions (
    id           BIGSERIAL PRIMARY KEY,
    form_id      BIGINT REFERENCES forms(id) ON DELETE CASCADE,
    user_id      BIGINT REFERENCES users(id),
    data         JSONB NOT NULL,
    ip_address   VARCHAR(45),
    user_agent   TEXT,
    submitted_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Staff & Resources

#### staff_members

```sql
CREATE TABLE staff_members (
    id              SERIAL PRIMARY KEY,
    user_id         BIGINT REFERENCES users(id),
    title           VARCHAR(255) NOT NULL,
    department      VARCHAR(100),
    bio             TEXT,
    photo_url       VARCHAR(500),
    email           VARCHAR(255),
    phone           VARCHAR(20),
    office_location VARCHAR(255),
    display_order   INTEGER DEFAULT 0,
    is_pastor       BOOLEAN DEFAULT FALSE,
    is_active       BOOLEAN DEFAULT TRUE,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

#### photos

```sql
CREATE TABLE photos (
    id             BIGSERIAL PRIMARY KEY,
    title          VARCHAR(255),
    description    TEXT,
    file_path      VARCHAR(500) NOT NULL,
    thumbnail_path VARCHAR(500),
    file_size      INTEGER,
    uploaded_by    BIGINT REFERENCES users(id),
    uploaded_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Audit & System

#### audit_log

```sql
CREATE TABLE audit_log (
    id          BIGSERIAL PRIMARY KEY,
    user_id     BIGINT REFERENCES users(id),
    action      VARCHAR(100) NOT NULL,
    entity_type VARCHAR(100),
    entity_id   BIGINT,
    ip_address  VARCHAR(45),
    user_agent  TEXT,
    success     BOOLEAN DEFAULT TRUE,
    details     JSONB,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

#### site_settings

```sql
CREATE TABLE site_settings (
    key         VARCHAR(100) PRIMARY KEY,
    value       TEXT,
    description TEXT,
    updated_by  BIGINT REFERENCES users(id),
    updated_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

---

## Authentication & Authorization

### Authentication Flows

#### Registration

1. User submits email, password, name
2. Generate verification token (UUID), hash with SHA-256, store hash
3. Create user with `is_verified = false`
4. Send verification email with plain token
5. User clicks link → hash token, compare → set `is_verified = true`
6. User can log in but has NO roles until admin assigns "member"

#### Login

1. Validate email/password (bcrypt)
2. Check `is_verified` and `locked_until`
3. On success: generate JWT, store session in Redis, set HTTP-only cookie
4. On failure: increment `failed_login_count`, lock after 5 failures (15 min)

#### Password Reset

1. Generate reset token (UUID), hash with SHA-256, store hash with 1-hour expiry
2. Send email with plain token
3. User submits new password → hash URL token, compare to stored hash
4. Update password, clear reset token

### JWT Structure

```json
{
  "jti": "unique-token-id",
  "user_id": 123,
  "email": "member@example.com",
  "roles": ["member", "volunteer"],
  "exp": 1234567890,
  "iat": 1234567890
}
```

### Redis Keys

```
blacklist:{jti}     # Revoked tokens (24h TTL)
csrf:{jti}          # CSRF tokens (24h TTL)
rate:login:{ip}     # Login rate limiting
rate:api:{ip}       # API rate limiting
```

### Permission Matrix

| Resource           | Public | Member | Volunteer | Musician | Staff | Elder | Pastor | Admin |
| ------------------ | ------ | ------ | --------- | -------- | ----- | ----- | ------ | ----- |
| Public pages       | ✓      | ✓      | ✓         | ✓        | ✓     | ✓     | ✓      | ✓     |
| Member directory   | -      | ✓      | -         | -        | ✓     | ✓     | ✓      | ✓     |
| Operating budget   | -      | ✓      | -         | -        | ✓     | ✓     | ✓      | ✓     |
| Small groups list  | -      | ✓      | -         | -        | ✓     | ✓     | ✓      | ✓     |
| Event registration | -      | ✓      | ✓         | ✓        | ✓     | ✓     | ✓      | ✓     |
| Volunteer schedule | -      | -      | ✓         | -        | ✓     | -     | -      | ✓     |
| Music schedule     | -      | -      | -         | ✓        | ✓     | -     | -      | ✓     |
| Manage content     | -      | -      | -         | -        | ✓     | -     | -      | ✓     |
| Ministry pages     | -      | ¹      | -         | -        | ✓     | -     | -      | ✓     |
| Prayer requests    | -      | -      | -         | -        | -     | ✓     | ✓      | ✓     |
| User management    | -      | -      | -         | -        | -     | -     | -      | ✓     |
| System settings    | -      | -      | -         | -        | -     | -     | -      | ✓     |

¹ Members can edit ministry pages only if assigned via `ministry_assignments`

### Middleware Examples

```go
// RequireAuth - must be logged in
r.Use(middleware.RequireAuth)

// RequireAnyRole - must have at least one role
r.Use(middleware.RequireAnyRole("staff", "admin"))

// RequireAllRoles - must have all roles
r.Use(middleware.RequireAllRoles("volunteer", "musician"))
```

---

## Code Organization

```
churchsite/
├── cmd/
│   ├── server/main.go           # Entry point
│   └── seed/main.go             # Database seeding
├── internal/
│   ├── config/config.go         # Configuration
│   ├── database/database.go     # DB connection
│   ├── models/                  # GORM models
│   ├── handlers/                # HTTP handlers
│   ├── middleware/              # Auth, rate limit, logging
│   ├── services/                # Business logic
│   └── utils/                   # Helpers
├── migrations/                  # SQL migration files
├── static/
│   ├── css/                     # base, layout, components, utilities, print
|   ├── fonts/
│   ├── js/                      # htmx.min.js, alpine.min.js
│   └── images/
├── templates/
│   ├── components/              # nav, footer, button, card
│   └── errors/                  # 404, 403, 500, maintenance
│   ├── layouts/                 # base.templ, member.templ
│   ├── pages/                   # home.templ, about.templ, ...
├── compose.yml                  # Base Docker Compose
├── compose.dev.yml
├── compose.staging.yml
├── compose.prod.yml
├── Dockerfile
├── Makefile
├── go.mod
└── go.sum
```

---

## Security Requirements

### Password Rules

- 8+ characters, uppercase, lowercase, number, special character
- bcrypt cost factor: 12

### Rate Limits

- Login: 5 attempts / 15 min / IP
- Password reset: 3 / hour / email
- Form submission: 10 / hour / IP
- API: 60 / min / IP
- Account lockout: 15 min after 5 failures

### File Uploads

- Allowed: .pdf, .jpg, .png, .webp
- Max size: 10MB
- Convert images to WebP (quality 85)
- Store outside web root, serve via X-Accel-Redirect

### CSRF Protection

- Token stored in Redis as `csrf:{jti}`
- Rendered in `<meta>` tag
- HTMX sends via `hx-headers` as `X-CSRF-Token`
- Validated on POST/PUT/PATCH/DELETE

---

## Content Workflows

### Bulletins

- Upload: `/staff/bulletins/upload` (PDF, max 10MB)
- Storage: `/storage/bulletins/{morning|evening}/`
- Retention: 14 days (auto-cleanup via `./sachapel cleanup-bulletins`)

### Events

- Recurrence patterns: none, daily, weekly, biweekly, monthly, yearly
- Edit/delete affects all occurrences (MVP)
- Visibility scheduling via `visible_from` / `visible_until`

### Member Directory

- Default: `directory_opt_in = FALSE`
- Separate toggles: show_email, show_phone, show_address
- If opted in but all toggles off: shows name and photo only

### Image Processing

1. Validate type (JPEG, PNG, WebP) and size (≤10MB)
2. Convert to WebP (lossy, quality 85)
3. Resize if >2000px on longest side
4. Generate thumbnail (300x300)
5. Store paths in database

---

## Environment Variables

```bash
APP_ENV=production|staging|development
APP_URL=https://sachapel.com
APP_PORT=3000

DATABASE_URL=postgres://user:pass@host:5432/db?sslmode=disable
REDIS_URL=redis://host:6379/0

JWT_SECRET=<64-char-random>
JWT_EXPIRATION=24h

SMTP_HOST=smtp.office365.com
SMTP_PORT=587
SMTP_USER=noreply@sachapel.com
SMTP_PASS=<password>
FROM_EMAIL=noreply@sachapel.com
FROM_NAME=Saint Andrew's Chapel

MAX_UPLOAD_SIZE=10485760
```

---

## Development Commands

```bash
make generate          # Generate Templ templates
make dev-up            # Start Docker containers
make dev-down          # Stop containers
make dev-logs          # Follow app logs
make migrate-up        # Run migrations
make migrate-down      # Rollback last migration
make migrate-create name=<name>  # Create migration
make seed              # Seed development data
make test              # Run tests
make lint              # Run golangci-lint
make build             # Build binary
make watch             # Hot reload (templ + air)
```

### Seed Data (Development)

- Admin: `admin@sachapel.test` / `AdminPass123!`
- Staff: `staff@sachapel.test` / `StaffPass123!`
- 5 sample members, 2 elders/pastors
- 10 ministries, 8 events, 5 announcements, 2 bulletins
- Event registration + new member form schemas

---

## Health Checks

| Endpoint        | Checks                  | Response                                                   |
| --------------- | ----------------------- | ---------------------------------------------------------- |
| `/health`       | App running             | `{"status": "ok"}`                                         |
| `/health/ready` | PostgreSQL + Redis ping | `{"status": "ok", "postgres": "ok", "redis": "ok"}` or 503 |
