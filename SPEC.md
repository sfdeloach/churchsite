# Saint Andrew's Chapel Website Redesign - Technical Specification

**Project:** Saint Andrew's Chapel Public Website Redesign  
**Version:** 1.2  
**Last Updated:** February 5, 2026  
**Repository:** https://github.com/sfdeloach/churchsite

---

## Table of Contents

1. [Project Overview](#project-overview)
2. [Feature Prioritization](#feature-prioritization)
3. [Technical Stack](#technical-stack)
4. [URL Routing Structure](#url-routing-structure)
5. [Database Schema](#database-schema)
6. [Authentication & Authorization](#authentication--authorization)
7. [Deployment Architecture](#deployment-architecture)
8. [Development Workflow](#development-workflow)
9. [Security Considerations](#security-considerations)
10. [Backup & Disaster Recovery](#backup--disaster-recovery)
11. [Content Management](#content-management)
12. [Third-Party Integrations](#third-party-integrations)
13. [Performance & Scalability](#performance--scalability)
14. [Accessibility Requirements](#accessibility-requirements)
15. [Future Enhancements](#future-enhancements)

---

## Project Overview

### Purpose

Complete redesign and rebuild of Saint Andrew's Chapel's public website (https://sachapel.com) with modern architecture, improved user experience, and self-contained functionality to minimize external dependencies.

### Core Principles

- **Self-contained**: Minimize external service dependencies for long-term sustainability
- **Future-proof**: Use stable, well-supported technologies
- **Mobile-first**: Responsive design prioritizing mobile experience
- **Performance**: Fast page loads and efficient resource usage
- **Maintainability**: Clear code structure for future developers
- **Accessibility**: WCAG 2.1 AA compliance for elderly and disabled users

### Target Audiences

1. **First-time visitors** - Seeking information about beliefs, services, and location
2. **Regular attenders** - Finding weekly bulletins, events, and announcements
3. **Members** - Accessing directory, private content, and registration forms
4. **Staff & ministry leaders** - Managing content, events, and forms

### Current Site Limitations

- Static HTML requiring manual file editing
- Poor mobile experience
- No content management capabilities
- External dependencies (Wufoo, MailChimp, Vimeo)
- Limited navigation and information architecture
- No member portal or authentication

---

## Feature Prioritization

### MVP (Minimum Viable Product) - Phase 1

**Priority 1: Public Content & Navigation**

1. Modern, responsive homepage with service times and upcoming events
2. About section (Who We Are, What We Believe, History & Identity)
3. Pastors & Staff directory with photos and bios
4. Place of Worship information with Google Maps integration
5. Improved navigation and information architecture
6. Contact information and directions

**Priority 2: Weekly Content Management** 7. Bulletin upload and management (2-week retention, morning/evening) 8. Events calendar with CRUD operations 9. Announcements system (text-based updates) 10. Event visibility scheduling (auto-publish/expire)

**Priority 3: Authentication & Member Portal** 11. User registration with email verification 12. Password reset (automated for verified emails) 13. Member directory (opt-in contact sharing) 14. Role-based access control (public, member, staff, elder, pastor, admin) 15. Current-year operating budget (members-only)

**Priority 4: Ministry Management** 16. Ministry pages with staff update capabilities 17. Ministry-specific event management 18. Small group directory (members-only with off-site locations)

**Priority 5: Forms & Registration** 19. Event registration system with capacity limits 20. Form submission storage and reporting 21. Email notifications for form submissions 22. New member application forms (digital, stored in database)

### Phase 2 - Post-Launch Enhancements

**Priority 6: Advanced Features** 23. Form builder UI for staff to create custom forms 24. Prayer request system (elders/pastors only) 25. Volunteer scheduling for children's ministry and musicians 26. Photo gallery with admin upload capabilities 27. Self-hosted analytics (Umami or Plausible) 28. Audit logging for login attempts

**Priority 7: Integration Features** 29. Breeze ChMS integration for giving (redirect or embed) 30. Enhanced Google Calendar integration 31. Email newsletter system (replace MailChimp)

### Phase 3 - Future Enhancements

**Priority 8: Advanced Ministry Tools** 32. Self-hosted sermon video with member authentication (replace Vimeo) 33. Sermon archive with search/filtering 34. Saint Andrew's Conservatory integration 35. Mobile app or PWA capabilities 36. Advanced reporting and analytics dashboards

---

## Technical Stack

### Backend Framework

**Go + Chi Router**

**Rationale:**

- Built entirely on Go's standard library net/http
- Minimal dependency with maximum compatibility
- Excellent AI/Claude Code assistance (extensive training data)
- Idiomatic Go with standard http.Handler patterns
- High performance with low resource usage on VPS
- Single binary deployment (no runtime dependencies)
- Future-proof: easy migration to pure net/http if needed
- Active maintenance with stable API

**Why Chi over alternatives:**

- **vs net/http alone**: Adds routing without abstraction overhead
- **vs Fiber**: More idiomatic Go, better net/http compatibility, no fasthttp dependency
- **vs Echo**: Simpler, more minimal, uses standard handlers
- **vs Gin**: More active maintenance, cleaner API, better AI code generation

**Key Features:**

- URL parameter extraction
- Middleware chaining
- Route grouping
- Standard http.Request/ResponseWriter
- Compatible with any net/http middleware

### Frontend Architecture

**Server-Side Rendered with Progressive Enhancement**

**Stack:**

- **Templ** for type-safe, compiled Go templates
- **HTMX** for dynamic interactions without full page reloads
- **Alpine.js** for lightweight client-side interactivity
- **Vanilla CSS** with design system (hybrid approach)
- **Vanilla JavaScript** for custom interactions when needed

**Rationale:**

- Type-safe templates with compile-time error checking
- Component-based architecture similar to React/JSX
- Better IDE support and refactoring safety
- Zero Node.js dependency for CSS (self-contained)
- Full control over styling with maintainable CSS
- Better SEO and initial page load performance
- Progressive enhancement for accessibility
- Minimal JavaScript bundle size

**Templ Benefits:**

- Templates compile to Go functions (faster rendering)
- Type-safe props and component composition
- IDE autocomplete and error detection
- Clean integration with HTMX/Alpine.js
- Excellent AI/Claude Code generation support

**CSS Approach - Hybrid Vanilla + Utilities:**

- CSS custom properties for design system
- Minimal utility classes for common patterns
- Semantic component classes
- No build tools required (optional minification)
- Church brand colors deeply integrated
- Easy maintenance by any developer
- Excellent AI code generation

### Database

**PostgreSQL 17**

**Rationale:**

- Latest stable version with performance improvements
- Robust user authentication and permission management
- JSONB support for flexible form data storage
- Full-text search capabilities for future sermon search
- Mature, stable, well-documented
- Excellent Docker support
- Strong data integrity and ACID compliance

**Schema Management:**

- **golang-migrate** for database migrations (explicit SQL files)
- **GORM** for ORM queries only (AutoMigrate disabled)
- Migration files in version control
- Never use GORM AutoMigrate in production

### Caching Layer

**Redis 7**

**Usage:**

- Session storage
- Page fragment caching (bulletins, events)
- Rate limiting for login attempts and form submissions

### File Storage

**Local Docker Volumes**

**Structure:**

```
/app/storage/
├── bulletins/
│   ├── morning/
│   └── evening/
├── photos/
├── documents/
└── uploads/
```

**Backup Strategy:**

- Daily automated backups via cron + bash script
- Database dumps (pg_dump)
- File system snapshots
- Secure copy (scp) to backup server or S3-compatible storage
- 30-day retention policy

### Image Processing

**disintegration/imaging (Go Library)**

**Rationale:**

- Pure Go, no CGO dependencies
- Simple API for common operations
- Aligns with self-contained philosophy

**Processing Pipeline:**

1. Upload received → validate type (JPEG, PNG, WebP) and size (max 10MB)
2. Generate thumbnail (300x300, fit within bounds, maintain aspect ratio)
3. Optionally resize original if larger than 2000px on longest side
4. Save both to storage directory
5. Store paths in database

**Thumbnail Dimensions:**

- Standard: 300x300 (fit within, maintain aspect ratio)
- Gallery grid: 200x200 (square crop for uniform grid)

### Email Service

**Self-Hosted SMTP with Fallback**

**Primary:** Microsoft 365 SMTP relay (if available)
**Fallback:** Resend or AWS SES for transactional emails
**Use Cases:**

- Email verification
- Password reset links
- Event registration confirmations
- Form submission notifications

**Implementation:**

- Go standard library `net/smtp`
- Template-based email generation
- Synchronous sending with goroutine retry on failure

### Logging

**Go slog (Standard Library)**

**Rationale:**

- Built into Go 1.21+ (no external dependency)
- Structured JSON logging out of the box
- Aligns with self-contained philosophy
- Easy integration with log aggregation if needed later

**Configuration:**

- Development: Text format with colors
- Production: JSON format for parsing
- Log levels: DEBUG, INFO, WARN, ERROR

### Analytics

**Umami Analytics (Self-Hosted)**

**Rationale:**

- Privacy-focused, GDPR compliant
- Self-hosted, no external dependencies
- Lightweight JavaScript tracker
- Simple Docker deployment
- Open-source with active development

**Alternative:** Plausible Analytics (also self-hosted)

### Authentication

**Custom JWT-based Authentication**

**Stack:**

- **golang-jwt** for token generation/validation
- **bcrypt** for password hashing (cost factor: 12)
- **Redis** for session storage
- **Secure HTTP-only cookies** for token storage

**Security Measures:**

- CSRF protection via SameSite cookies and CSRF tokens
- Rate limiting on login attempts (5 attempts per 15 minutes)
- Email verification required before access
- Secure password reset with time-limited tokens
- Account lockout after repeated failed attempts

### Form Builder (Phase 2)

**Custom JSON-Schema Based System**

**Implementation:**

- JSON schema definition for form structure
- Dynamic form rendering from schema
- Validation based on schema rules
- Form submissions stored as JSONB in PostgreSQL
- Staff UI to build forms with drag-and-drop interface

### CSS Architecture

**Hybrid Vanilla CSS + Minimal Utilities**

**Approach:**

- CSS custom properties (variables) for design system
- Semantic component classes for reusable UI elements
- Minimal utility classes for common patterns
- No build tools required (zero Node.js dependency)
- Optional minification for production

**Design System Structure:**

```css
/* CSS Variables for Consistency */
:root {
  --color-firebrick: #b22222;
  --color-gray-50: #f9fafb;
  --color-gray-900: #111827;

  --spacing-xs: 0.5rem;
  --spacing-sm: 1rem;
  --spacing-md: 1.5rem;
  --spacing-lg: 2rem;
  --spacing-xl: 3rem;

  --font-sans: -apple-system, BlinkMacSystemFont, "Segoe UI", sans-serif;

  --shadow-sm: 0 1px 2px rgba(0, 0, 0, 0.05);
  --shadow-md: 0 4px 6px rgba(0, 0, 0, 0.1);
}
```

**File Organization:**

```
static/css/
├── base.css          # Reset, typography, variables
├── layout.css        # Containers, grid, flex utilities
├── components.css    # Buttons, cards, forms, nav
├── utilities.css     # Spacing, colors, responsive helpers
└── print.css         # Print styles for bulletins
```

**Rationale:**

- Zero external dependencies (aligns with project philosophy)
- Full control over every style
- Excellent AI/Claude Code generation support
- Easy maintenance by any web developer
- Custom church branding naturally integrated
- Small file size (15-25KB total)
- Future-proof (CSS is stable and universal)

### Build Tools

**Minimal Tooling**

- **Templ CLI** - Template generation and watching
- **Air** - Hot reload for Go development
- **Make** - Build automation and common tasks
- **(Optional) CSS minifier** - For production optimization

### CI/CD Pipeline

**GitHub Actions**

**Workflows:**

1. **Test & Lint** (on pull request)
   - Install Templ CLI
   - Generate Templ templates
   - Run Go tests
   - Run linters (golangci-lint)
   - Check migrations
2. **Build & Deploy Staging** (on push to `develop` branch)
   - Generate Templ templates
   - Build Docker images
   - Deploy to staging.sachapel.com
   - Run smoke tests
3. **Deploy Production** (on push to `main` branch)
   - Generate Templ templates
   - Build production Docker images
   - Tag release
   - Deploy to sachapel.com
   - Send deployment notification

### Containerization

**Docker + Docker Compose**

**Services:**

```yaml
services:
  app: # Go/Chi application
  postgres: # PostgreSQL 18 database
  redis: # Redis 8 cache
  nginx: # Reverse proxy / TLS termination
  umami: # Analytics (optional, Phase 2)
```

**Environment-Specific Compose Files:**

- `compose.yml` - Base configuration
- `compose.dev.yml` - Development overrides
- `compose.staging.yml` - Staging configuration
- `compose.prod.yml` - Production configuration

---

## URL Routing Structure

### Public Routes (No Authentication)

```
/                                    # Homepage
/about/history                       # History & Identity
/about/beliefs                       # Doctrine & Beliefs
/about/worship                       # Theology of Worship
/about/gospel                        # What Is the Gospel?
/about/staff                         # Pastors & Staff
/about/building                      # Our Place of Worship

/ministries                          # Ministry overview
/ministries/sunday-school            # Sunday School
/ministries/children                 # Children's Ministry
/ministries/music                    # Music Ministry
/ministries/youth                    # Youth Ministry
/ministries/college-career           # College & Career
/ministries/men                      # Men's Ministry
/ministries/women                    # Women's Ministry
/ministries/bible-studies            # Weekly Bible Studies (public info)
/ministries/disaster-response        # Disaster Response
/ministries/prayer                   # Prayer Ministry
/ministries/salt-light               # Salt & Light
/ministries/sanctity-of-life         # Sanctity of Life

/calendar/events                     # Public events calendar
/calendar/events/:id                 # Individual event details

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
/verify-email/:token                 # Email verification confirmation
```

### Member-Only Routes (Authentication Required)

```
/member/dashboard                    # Member dashboard
/member/directory                    # Member directory (opt-in)
/member/profile                      # Edit own profile
/member/budget                       # Current year operating budget
/member/bible-studies                # Small group directory with locations

/member/events/register/:id          # Event registration forms
```

### Staff Routes (Staff+ Authentication)

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
/staff/events/:id/registrations      # View event registrations

/staff/forms/:id/submissions         # View form submissions
/staff/forms/:id/export              # Export submissions as CSV

/staff/ministry/:slug                # Edit specific ministry page
```

### Elder/Pastor Routes (Elder+ Authentication)

```
/elder/dashboard                     # Elder/Pastor dashboard
/elder/prayer-requests               # Prayer requests
/elder/prayer-requests/:id           # Individual request
```

### Admin Routes (Admin Authentication)

```
/admin/dashboard                     # Admin dashboard
/admin/users                         # Manage users
/admin/users/:id                     # Edit user / assign roles
/admin/users/pending                 # Pending email verifications
/admin/users/invitations             # Send member invitations

/admin/photos                        # Photo gallery management
/admin/photos/upload                 # Upload photos

/admin/forms                         # Form builder (Phase 2)
/admin/forms/create                  # Create new form
/admin/forms/:id/edit                # Edit form

/admin/analytics                     # Analytics dashboard
/admin/audit-log                     # Login attempt logs

/admin/settings                      # Site settings
/admin/settings/email                # Email configuration
/admin/settings/maintenance          # Maintenance mode
```

### API Routes (JSON Responses)

```
/api/v1/events                       # List events (public)
/api/v1/events/:id                   # Event details
/api/v1/events/:id/register          # POST: Register for event

/api/v1/bulletins/current            # Get current bulletins

/api/v1/forms/:id/submit             # POST: Submit form

/api/v1/member/directory/search      # Search member directory (auth)
```

### Health & Monitoring Routes

```
/health                              # Basic health check (returns 200 OK)
/health/ready                        # Readiness check (DB + Redis connected)
```

### Static Assets

```
/static/css/                         # Compiled CSS
/static/js/                          # JavaScript files
/static/images/                      # Images (logos, icons)
/static/fonts/                       # Custom fonts (if any)

/uploads/bulletins/                  # Uploaded bulletins (protected)
/uploads/photos/                     # Photo gallery
```

---

## Database Schema

### GORM Model Convention

All tables will use GORM's `gorm.Model` which automatically includes:

```go
type Model struct {
    ID        uint           `gorm:"primarykey"`
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}
```

This means:

- `id` is automatically `BIGINT UNSIGNED` (or `BIGSERIAL` equivalent)
- `created_at` and `updated_at` are automatically managed
- `deleted_at` enables soft deletes (records marked as deleted but not removed)
- Indexes are automatically created for primary keys and DeletedAt

**Note:** The SQL schemas below show the explicit columns for clarity, but when using GORM structs, these timestamp fields will be inherited from `gorm.Model` and don't need to be explicitly defined in the migration files.

### Soft Delete Policy

**Tables with Soft Deletes (include deleted_at column):**

- `users` - Account recovery, audit trail
- `events` - May need to restore cancelled events
- `announcements` - Historical reference
- `bulletins` - Historical reference
- `ministries` - Preserve relationships
- `forms` - Preserve submission references
- `staff_members` - Historical reference

**Tables with Hard Deletes:**

- `user_roles` - Junction table, no history needed
- `event_registrations` - Can be recreated
- `form_submissions` - Optional: keep for compliance
- `audit_log` - Never delete (or time-based retention)
- `site_settings` - Overwrite only, no delete
- `member_profiles` - Cascade from users
- `small_groups` - Hard delete acceptable
- `photos` - Hard delete with file cleanup

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
    verification_token  VARCHAR(255),
    reset_token         VARCHAR(255),
    reset_token_expires TIMESTAMP,
    last_login          TIMESTAMP,
    failed_login_count  INTEGER DEFAULT 0,
    locked_until        TIMESTAMP,
    created_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_verification_token ON users(verification_token);
CREATE INDEX idx_users_reset_token ON users(reset_token);
```

#### roles

```sql
CREATE TABLE roles (
    id          SERIAL PRIMARY KEY,
    name        VARCHAR(50) UNIQUE NOT NULL,
    description TEXT,
    created_at  TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Seed data
INSERT INTO roles (name, description) VALUES
    ('public', 'General public - no authentication'),
    ('member', 'Verified church member'),
    ('deacon', 'Deacon'),
    ('elder', 'Ruling elder'),
    ('staff', 'Church staff member'),
    ('musician', 'Music ministry staff'),
    ('pastor', 'Pastor'),
    ('volunteer', 'Ministry volunteer'),
    ('admin', 'System administrator');
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

CREATE INDEX idx_user_roles_user_id ON user_roles(user_id);
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

CREATE INDEX idx_announcements_visible ON announcements(visible_from, visible_until);
CREATE INDEX idx_announcements_active ON announcements(is_active);
```

#### events

```sql
CREATE TABLE events (
    id                  BIGSERIAL PRIMARY KEY,
    title               VARCHAR(255) NOT NULL,
    description         TEXT,
    event_date          TIMESTAMP NOT NULL,
    end_date            TIMESTAMP,
    location            VARCHAR(255),
    location_details    TEXT,
    is_recurring        BOOLEAN DEFAULT FALSE,
    recurrence_rule     VARCHAR(20),  -- none, daily, weekly, biweekly, monthly, yearly
    recurrence_end      DATE,         -- When recurring events stop
    registration_enabled BOOLEAN DEFAULT FALSE,
    capacity_limit      INTEGER,
    registration_deadline TIMESTAMP,
    visible_from        TIMESTAMP,
    visible_until       TIMESTAMP,
    is_public           BOOLEAN DEFAULT TRUE,
    ministry_id         INTEGER REFERENCES ministries(id),
    created_by          BIGINT REFERENCES users(id),
    created_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at          TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_events_date ON events(event_date);
CREATE INDEX idx_events_visible ON events(visible_from, visible_until);
CREATE INDEX idx_events_public ON events(is_public);
```

#### event_registrations

```sql
CREATE TABLE event_registrations (
    id              BIGSERIAL PRIMARY KEY,
    event_id        BIGINT REFERENCES events(id) ON DELETE CASCADE,
    user_id         BIGINT REFERENCES users(id),
    guest_name      VARCHAR(255),
    guest_email     VARCHAR(255),
    guest_phone     VARCHAR(20),
    number_attending INTEGER DEFAULT 1,
    special_needs   TEXT,
    form_data       JSONB, -- Additional custom form fields
    registered_at   TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_registration UNIQUE(event_id, user_id)
);

CREATE INDEX idx_registrations_event ON event_registrations(event_id);
CREATE INDEX idx_registrations_user ON event_registrations(user_id);
```

#### bulletins

```sql
CREATE TABLE bulletins (
    id              BIGSERIAL PRIMARY KEY,
    bulletin_date   DATE NOT NULL,
    service_type    VARCHAR(20) NOT NULL, -- 'morning' or 'evening'
    file_path       VARCHAR(500) NOT NULL,
    file_size       INTEGER,
    uploaded_by     BIGINT REFERENCES users(id),
    uploaded_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    CONSTRAINT unique_bulletin UNIQUE(bulletin_date, service_type)
);

CREATE INDEX idx_bulletins_date ON bulletins(bulletin_date DESC);
```

#### ministries

```sql
CREATE TABLE ministries (
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(255) NOT NULL,
    slug            VARCHAR(255) UNIQUE NOT NULL,
    description     TEXT,
    leader_id       BIGINT REFERENCES users(id),
    contact_email   VARCHAR(255),
    meeting_time    VARCHAR(255),
    location        VARCHAR(255),
    is_active       BOOLEAN DEFAULT TRUE,
    sort_order      INTEGER DEFAULT 0,
    page_content    TEXT, -- HTML content for ministry page
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

#### small_groups

```sql
CREATE TABLE small_groups (
    id              SERIAL PRIMARY KEY,
    name            VARCHAR(255) NOT NULL,
    description     TEXT,
    leader_id       BIGINT REFERENCES users(id),
    meeting_day     VARCHAR(20),
    meeting_time    TIME,
    location        VARCHAR(255),
    is_active       BOOLEAN DEFAULT TRUE,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Form System Tables

#### forms

```sql
CREATE TABLE forms (
    id              BIGSERIAL PRIMARY KEY,
    title           VARCHAR(255) NOT NULL,
    description     TEXT,
    schema          JSONB NOT NULL, -- Form field definitions
    is_active       BOOLEAN DEFAULT TRUE,
    created_by      BIGINT REFERENCES users(id),
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

#### form_submissions

```sql
CREATE TABLE form_submissions (
    id              BIGSERIAL PRIMARY KEY,
    form_id         BIGINT REFERENCES forms(id) ON DELETE CASCADE,
    user_id         BIGINT REFERENCES users(id),
    data            JSONB NOT NULL, -- Form submission data
    ip_address      VARCHAR(45),
    user_agent      TEXT,
    submitted_at    TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_submissions_form ON form_submissions(form_id);
CREATE INDEX idx_submissions_user ON form_submissions(user_id);
CREATE INDEX idx_submissions_date ON form_submissions(submitted_at DESC);
```

### Staff & Resources Tables

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
    id              BIGSERIAL PRIMARY KEY,
    title           VARCHAR(255),
    description     TEXT,
    file_path       VARCHAR(500) NOT NULL,
    thumbnail_path  VARCHAR(500),
    file_size       INTEGER,
    uploaded_by     BIGINT REFERENCES users(id),
    uploaded_at     TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Audit & System Tables

#### audit_log

```sql
CREATE TABLE audit_log (
    id              BIGSERIAL PRIMARY KEY,
    user_id         BIGINT REFERENCES users(id),
    action          VARCHAR(100) NOT NULL,
    entity_type     VARCHAR(100),
    entity_id       BIGINT,
    ip_address      VARCHAR(45),
    user_agent      TEXT,
    success         BOOLEAN DEFAULT TRUE,
    details         JSONB,
    created_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX idx_audit_user ON audit_log(user_id);
CREATE INDEX idx_audit_action ON audit_log(action);
CREATE INDEX idx_audit_date ON audit_log(created_at DESC);
```

**Audit Log Policy:**

_Actions Logged:_

- **Authentication:** Login success/failure, logout, password reset request/completion
- **User Management:** Create, update, delete, role changes
- **Content Changes:** Create/update/delete events, announcements, bulletins
- **Form Submissions:** Submission received (not the data itself)
- **Administrative:** Settings changes, user role assignments

_Not Logged (to reduce volume):_

- Page views (use analytics for this)
- Read-only API calls
- Session refreshes

_Retention:_

- Production: 90 days
- Cleanup via daily cron job at 4 AM

#### site_settings

```sql
CREATE TABLE site_settings (
    key             VARCHAR(100) PRIMARY KEY,
    value           TEXT,
    description     TEXT,
    updated_by      BIGINT REFERENCES users(id),
    updated_at      TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);
```

### Phase 2 Tables (To Be Designed)

The following tables will be designed during Phase 2 implementation:

- `prayer_requests` - Prayer request submissions (elder/pastor access only)
- `volunteer_schedules` - Children's ministry and musician scheduling
- `newsletter_subscribers` - Email newsletter subscription management

---

## Authentication & Authorization

### Authentication Flow

#### Registration & Email Verification

1. User submits registration form with email, password, name
2. System generates verification token (UUID)
3. User record created with `is_verified = false`
4. Verification email sent with link: `/verify-email/:token`
5. User clicks link, token validated, `is_verified` set to `true`
6. User can now log in

#### Login Flow

1. User submits email and password
2. System checks if account exists and is verified
3. Password compared using bcrypt
4. If failed login count >= 5, check if account locked
5. On success:
   - Generate JWT token with claims (user_id, email, roles)
   - Store session in Redis with 24-hour TTL
   - Set HTTP-only, Secure, SameSite cookie
   - Reset failed login count
   - Log successful login
6. On failure:
   - Increment failed login count
   - Lock account for 15 minutes after 5 failures
   - Log failed attempt

#### Password Reset Flow

1. User requests password reset with email
2. System generates reset token (UUID) with 1-hour expiration
3. Reset email sent with link: `/reset-password/:token`
4. User clicks link, enters new password
5. Token validated (not expired, matches user)
6. Password updated, reset token cleared
7. User redirected to login

### Member Onboarding Flow

**Registration Types:**

1. **Self-Registration**
   - User registers with email/password
   - Email verification required
   - Starts with NO roles (verified user, not yet member)
   - Admin must assign "member" role after verification

2. **Admin Invitation**
   - Admin sends invitation email from `/admin/users/invitations`
   - User clicks link, sets password
   - Automatically assigned "member" role
   - Email pre-verified

**Directory Opt-In:**

- Default: `directory_opt_in = FALSE` (privacy by default)
- User can enable via profile settings at `/member/profile`
- Separate toggles for: email, phone, address visibility
- Admin cannot force opt-in (user consent required)

### Authorization Model

#### Role-Based Access Control (RBAC)

Instead of a simple hierarchy, this application uses **Role-Based Access Control** where roles grant specific permissions. Users can have multiple roles, and each role provides access to specific resources.

**Why RBAC over Hierarchy:**

- Musicians may not be church members but need access to music ministry features
- Volunteers may need specific access without member privileges
- Staff members may have limited administrative access
- More flexible as church needs evolve

#### Available Roles

```go
const (
    RolePublic    = "public"    // No authentication required
    RoleMember    = "member"    // Verified church member
    RoleVolunteer = "volunteer" // Ministry volunteer
    RoleMusician  = "musician"  // Music ministry (may not be member)
    RoleStaff     = "staff"     // Church staff
    RoleDeacon    = "deacon"    // Deacon
    RoleElder     = "elder"     // Ruling elder
    RolePastor    = "pastor"    // Pastor
    RoleAdmin     = "admin"     // Full system access
)
```

#### Permission Matrix

| Resource             | Public | Member | Volunteer | Musician | Staff | Deacon | Elder | Pastor | Admin |
| -------------------- | ------ | ------ | --------- | -------- | ----- | ------ | ----- | ------ | ----- |
| Public pages         | ✓      | ✓      | ✓         | ✓        | ✓     | ✓      | ✓     | ✓      | ✓     |
| Member directory     | -      | ✓      | -         | -        | ✓     | ✓      | ✓     | ✓      | ✓     |
| Operating budget     | -      | ✓      | -         | -        | ✓     | ✓      | ✓     | ✓      | ✓     |
| Small groups list    | -      | ✓      | -         | -        | ✓     | ✓      | ✓     | ✓      | ✓     |
| Event registration   | -      | ✓      | ✓         | ✓        | ✓     | ✓      | ✓     | ✓      | ✓     |
| Volunteer schedule   | -      | -      | ✓         | -        | ✓     | -      | -     | -      | ✓     |
| Music schedule       | -      | -      | -         | ✓        | ✓     | -      | -     | -      | ✓     |
| Manage events        | -      | -      | -         | -        | ✓     | -      | -     | -      | ✓     |
| Manage bulletins     | -      | -      | -         | -        | ✓     | -      | -     | -      | ✓     |
| Manage announcements | -      | -      | -         | -        | ✓     | -      | -     | -      | ✓     |
| Ministry page edits  | -      | -      | -         | -        | ✓     | -      | -     | -      | ✓     |
| Prayer requests      | -      | -      | -         | -        | -     | -      | ✓     | ✓      | ✓     |
| User management      | -      | -      | -         | -        | -     | -      | -     | -      | ✓     |
| System settings      | -      | -      | -         | -        | -     | -      | -     | -      | ✓     |

**Notes:**

- Users can have multiple roles (e.g., a member who is also a volunteer and musician)
- Admin role has access to all resources
- Role combinations are additive (having both "volunteer" and "musician" grants both sets of permissions)

#### Route Protection Middleware

```go
// RequireAuth - User must be logged in
func RequireAuth(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Check for valid session/JWT
        // If not authenticated, redirect to login or return 401
        next.ServeHTTP(w, r)
    })
}

// RequireAnyRole - User must have at least one of the specified roles
func RequireAnyRole(roles ...string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Check if user has any of the required roles
            // If not, return 403 Forbidden
            next.ServeHTTP(w, r)
        })
    }
}

// RequireAllRoles - User must have all of the specified roles
func RequireAllRoles(roles ...string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            // Check if user has all required roles
            // If not, return 403 Forbidden
            next.ServeHTTP(w, r)
        })
    }
}

// RequireAdmin - User must have admin role
func RequireAdmin(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        // Check if user has admin role
        // If not, return 403 Forbidden
        next.ServeHTTP(w, r)
    })
}
```

#### Example Route Protection with Chi

```go
// Member-only route
r.Group(func(r chi.Router) {
    r.Use(middleware.RequireAuth)
    r.Use(middleware.RequireAnyRole("member"))
    r.Get("/member/directory", handlers.MemberDirectory)
})

// Staff can edit events
r.Group(func(r chi.Router) {
    r.Use(middleware.RequireAuth)
    r.Use(middleware.RequireAnyRole("staff", "admin"))
    r.Post("/staff/events", handlers.CreateEvent)
})

// Elder OR pastor can access prayer requests
r.Group(func(r chi.Router) {
    r.Use(middleware.RequireAuth)
    r.Use(middleware.RequireAnyRole("elder", "pastor", "admin"))
    r.Get("/elder/prayer-requests", handlers.PrayerRequests)
})

// Volunteer AND musician scheduling (user needs both roles)
r.Group(func(r chi.Router) {
    r.Use(middleware.RequireAuth)
    r.Use(middleware.RequireAllRoles("volunteer", "musician"))
    r.Get("/schedule/music-volunteer", handlers.MusicVolunteerSchedule)
})

// Admin-only user management
r.Group(func(r chi.Router) {
    r.Use(middleware.RequireAdmin)
    r.Get("/admin/users", handlers.ManageUsers)
})
```

### Session Management

#### Hybrid JWT + Redis Approach

- **JWT Token:** Contains user claims (id, email, roles), verified by signature
- **Redis Blacklist:** Stores revoked token IDs for immediate invalidation
- **Token Lifetime:** 24 hours (short enough to limit exposure)

#### Flow

1. User logs in → JWT issued, stored in HTTP-only cookie
2. Each request → JWT signature verified (no Redis lookup)
3. User logs out → Token ID added to Redis blacklist (24h TTL)
4. Password change → All user tokens added to blacklist
5. Blacklist checked only when token would otherwise be valid

#### JWT Token Structure

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

#### Redis Keys

```
blacklist:{jti}          # Revoked token (TTL matches token expiry)
rate:login:{ip}          # Login rate limiting
rate:api:{ip}            # API rate limiting
```

#### Security Measures

- JWT secret stored in environment variable (32+ random bytes)
- Tokens expire after 24 hours
- CSRF token validation on state-changing requests
- Rate limiting on login endpoint (5 attempts per 15 minutes per IP)
- Account lockout after repeated failures

#### API Authentication

**Dual Authentication Support:**
The API accepts authentication via two methods (checked in order):

1. **HTTP-only Cookie** (web requests)
   - Automatically sent by browser
   - Used by HTMX and standard web requests
   - CSRF token required for mutations

2. **Bearer Token** (API/mobile requests)
   - Header: `Authorization: Bearer <jwt-token>`
   - Used by mobile apps, third-party integrations
   - No CSRF required (token proves intent)

**Token Acquisition:**

- Web: Cookie set automatically on login
- API: POST to `/api/v1/auth/login` returns token in response body

**Middleware Priority:**

1. Check for valid cookie → use if present
2. Check for Bearer header → use if present
3. No auth → treat as anonymous (public routes only)

---

## Deployment Architecture

### Infrastructure Overview

```
┌─────────────────────────────────────────────────────────┐
│                    Internet / CDN                        │
└────────────────────┬────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────┐
│                  VPS (Single Server)                     │
│  ┌──────────────────────────────────────────────────┐  │
│  │  Docker Host (Ubuntu 22.04 LTS)                   │  │
│  │  ┌────────────────────────────────────────────┐  │  │
│  │  │  Nginx (Container)                          │  │  │
│  │  │  - Port 80/443                              │  │  │
│  │  │  - TLS Termination (Let's Encrypt)          │  │  │
│  │  │  - Reverse Proxy                            │  │  │
│  │  │  - Static File Serving                      │  │  │
│  │  │  - Rate Limiting                            │  │  │
│  │  └────────────────────────────────────────────┘  │  │
│  │               │                                    │  │
│  │               ▼                                    │  │
│  │  ┌────────────────────────────────────────────┐  │  │
│  │  │  Go/Chi App (Container)                   │  │  │
│  │  │  - Port 3000 (internal)                     │  │  │
│  │  │  - Application Logic                        │  │  │
│  │  │  - Template Rendering                       │  │  │
│  │  │  - API Endpoints                            │  │  │
│  │  └────────────────────────────────────────────┘  │  │
│  │       │               │                           │  │
│  │       ▼               ▼                           │  │
│  │  ┌─────────┐     ┌──────────┐                    │  │
│  │  │PostgreSQL│     │  Redis   │                    │  │
│  │  │Container │     │Container │                    │  │
│  │  │Port 5432 │     │Port 6379 │                    │  │
│  │  └─────────┘     └──────────┘                    │  │
│  │       │               │                           │  │
│  │       ▼               ▼                           │  │
│  │  ┌────────────────────────────┐                  │  │
│  │  │    Docker Volumes           │                  │  │
│  │  │  - postgres_data            │                  │  │
│  │  │  - redis_data               │                  │  │
│  │  │  - app_uploads              │                  │  │
│  │  └────────────────────────────┘                  │  │
│  └──────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────┘
```

### Environment Setup

#### Production (sachapel.com)

- **VPS Provider:** Hetzner Cloud (CX21: 2 vCPU, 4GB RAM, 40GB SSD - ~€5/month)
- **Operating System:** Ubuntu 22.04 LTS
- **Docker:** Docker Engine + Docker Compose
- **Domain:** sachapel.com (DNS managed via registrar)
- **SSL/TLS:** Let's Encrypt (automated renewal via certbot)

#### Staging (staging.sachapel.com)

- Same VPS as production, separate containers
- Separate Docker Compose stack
- Separate PostgreSQL database
- Separate Redis instance
- Separate volumes

### Docker Compose Configuration

#### compose.yml (Base)

```yaml
version: "3.8"

services:
  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf:ro
      - ./nginx/sites-enabled:/etc/nginx/sites-enabled:ro
      - ./certbot/conf:/etc/letsencrypt:ro
      - ./certbot/www:/var/www/certbot:ro
      - app_uploads:/var/www/uploads:ro
    depends_on:
      - app
    restart: unless-stopped

  app:
    build:
      context: .
      dockerfile: Dockerfile
    environment:
      - APP_ENV=${APP_ENV}
      - DATABASE_URL=${DATABASE_URL}
      - REDIS_URL=${REDIS_URL}
      - JWT_SECRET=${JWT_SECRET}
      - SMTP_HOST=${SMTP_HOST}
      - SMTP_PORT=${SMTP_PORT}
      - SMTP_USER=${SMTP_USER}
      - SMTP_PASS=${SMTP_PASS}
    volumes:
      - app_uploads:/app/storage
    depends_on:
      - postgres
      - redis
    healthcheck:
      test: ["CMD", "wget", "-q", "--spider", "http://localhost:3000/health"]
      interval: 30s
      timeout: 10s
      retries: 3
      start_period: 10s
    restart: unless-stopped

  postgres:
    image: postgres:17-alpine
    environment:
      - POSTGRES_DB=${POSTGRES_DB}
      - POSTGRES_USER=${POSTGRES_USER}
      - POSTGRES_PASSWORD=${POSTGRES_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    restart: unless-stopped

  redis:
    image: redis:7-alpine
    command: redis-server --appendonly yes
    volumes:
      - redis_data:/data
    restart: unless-stopped

volumes:
  postgres_data:
  redis_data:
  app_uploads:
```

#### compose.prod.yml (Production Overrides)

```yaml
version: "3.8"

services:
  app:
    build:
      context: .
      dockerfile: Dockerfile
      target: production
    environment:
      - APP_ENV=production
```

#### compose.staging.yml (Staging Overrides)

```yaml
version: "3.8"

services:
  nginx:
    ports:
      - "8080:80"
      - "8443:443"

  app:
    environment:
      - APP_ENV=staging
```

### Dockerfile (Multi-stage Build)

```dockerfile
# Build stage
FROM golang:1.22-alpine AS builder

WORKDIR /build

# Install build dependencies including Templ
RUN apk add --no-cache git make

# Install Templ CLI
RUN go install github.com/a-h/templ/cmd/templ@latest

# Copy go mod files
COPY go.mod go.sum ./
RUN go mod download

# Copy source code
COPY . .

# Generate Templ templates
RUN templ generate

# Build binary
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/sachapel ./cmd/server

# Production stage
FROM alpine:latest AS production

WORKDIR /app

# Install runtime dependencies
RUN apk add --no-cache ca-certificates tzdata

# Copy binary from builder
COPY --from=builder /app/sachapel .

# Copy static assets (CSS, JS, images)
COPY --from=builder /build/static ./static

# Copy migrations
COPY --from=builder /build/migrations ./migrations

# Note: Generated _templ.go files are compiled into binary
# No need to copy .templ files to production

# Create storage directories
RUN mkdir -p /app/storage/bulletins/morning \
             /app/storage/bulletins/evening \
             /app/storage/photos \
             /app/storage/documents \
             /app/storage/uploads

# Non-root user
RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser && \
    chown -R appuser:appuser /app

USER appuser

EXPOSE 3000

CMD ["./sachapel"]
```

### Nginx Configuration

#### /etc/nginx/sites-enabled/sachapel.com

```nginx
# Rate limiting
limit_req_zone $binary_remote_addr zone=login_limit:10m rate=5r/m;
limit_req_zone $binary_remote_addr zone=api_limit:10m rate=60r/m;

# Redirect HTTP to HTTPS
server {
    listen 80;
    server_name sachapel.com www.sachapel.com;

    location /.well-known/acme-challenge/ {
        root /var/www/certbot;
    }

    location / {
        return 301 https://$server_name$request_uri;
    }
}

# HTTPS Server
server {
    listen 443 ssl http2;
    server_name sachapel.com www.sachapel.com;

    # SSL Configuration
    ssl_certificate /etc/letsencrypt/live/sachapel.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/sachapel.com/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;

    # Security Headers
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;

    # Static files - proxied to app with cache headers
    location /static/ {
        proxy_pass http://app:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        expires 1y;
        add_header Cache-Control "public, immutable";
    }

    # Uploaded files (protected)
    location /uploads/ {
        internal;
        alias /var/www/uploads/;
    }

    # Rate limit login
    location /login {
        limit_req zone=login_limit burst=2 nodelay;
        proxy_pass http://app:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # Rate limit API
    location /api/ {
        limit_req zone=api_limit burst=10 nodelay;
        proxy_pass http://app:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # Proxy all other requests to app
    location / {
        proxy_pass http://app:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;

        # WebSocket support (for future features)
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }

    # File upload size limit
    client_max_body_size 10M;
}
```

### Environment Variables

#### .env.production

```bash
# Application
APP_ENV=production
APP_URL=https://sachapel.com
APP_PORT=3000

# Database
DATABASE_URL=postgres://sacuser:STRONG_PASSWORD@postgres:5432/sachapel?sslmode=disable
POSTGRES_DB=sachapel
POSTGRES_USER=sacuser
POSTGRES_PASSWORD=STRONG_PASSWORD

# Redis
REDIS_URL=redis://redis:6379/0

# JWT
JWT_SECRET=RANDOM_64_CHAR_STRING
JWT_EXPIRATION=24h

# Email
SMTP_HOST=smtp.office365.com
SMTP_PORT=587
SMTP_USER=noreply@sachapel.com
SMTP_PASS=SMTP_PASSWORD
FROM_EMAIL=noreply@sachapel.com
FROM_NAME=Saint Andrew's Chapel

# File Upload
MAX_UPLOAD_SIZE=10485760  # 10MB in bytes

# Rate Limiting
LOGIN_RATE_LIMIT=5
LOGIN_RATE_WINDOW=15m
ACCOUNT_LOCKOUT_DURATION=15m

# Session
SESSION_DURATION=24h
COOKIE_SECURE=true
COOKIE_HTTPONLY=true
COOKIE_SAMESITE=strict
```

---

## Development Workflow

### Local Development Setup

#### Prerequisites

- Docker & Docker Compose
- Go 1.22+
- Git
- Make

#### Initial Setup

```bash
# Clone repository
git clone https://github.com/sfdeloach/churchsite.git
cd churchsite

# Copy environment template
cp .env.example .env.development

# Start Docker containers
make dev-up

# Run database migrations
make migrate-up

# Seed development data
make seed

# Access application
open http://localhost:3000
```

#### Development Commands (Makefile)

```makefile
.PHONY: dev-up dev-down dev-logs migrate-up migrate-down migrate-create seed test lint build generate watch

generate:
	templ generate

dev-up:
	docker compose -f compose.yml -f compose.dev.yml up -d

dev-down:
	docker compose -f compose.yml -f compose.dev.yml down

dev-logs:
	docker compose logs -f app

migrate-up:
	migrate -path ./migrations -database "$(DATABASE_URL)" up

migrate-down:
	migrate -path ./migrations -database "$(DATABASE_URL)" down 1

migrate-create:
	migrate create -ext sql -dir ./migrations -seq $(name)

seed:
	go run cmd/seed/main.go

test:
	templ generate
	go test -v ./...

test-coverage:
	templ generate
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out

lint:
	golangci-lint run

build:
	templ generate
	go build -o bin/sachapel cmd/server/main.go

watch:
	templ generate --watch & air
```

### Git Workflow

#### Branch Strategy

```
main              # Production-ready code
├── develop       # Integration branch
    ├── feature/* # Feature branches
    ├── bugfix/*  # Bug fix branches
    └── hotfix/*  # Emergency fixes
```

#### Workflow

1. Create feature branch from `develop`
2. Make changes, commit with conventional commits
3. Push branch, open Pull Request to `develop`
4. CI runs tests and lints
5. Code review and approval
6. Merge to `develop`, deploys to staging
7. Test on staging
8. Merge `develop` to `main`, deploys to production

#### Conventional Commits

```
feat: Add event registration capacity limits
fix: Correct bulletin upload validation
docs: Update API documentation
style: Format code with gofmt
refactor: Simplify authentication middleware
test: Add tests for user registration
chore: Update dependencies
```

### CI/CD Pipeline (GitHub Actions)

#### .github/workflows/test.yml

```yaml
name: Test & Lint

on:
  pull_request:
    branches: [develop, main]

jobs:
  test:
    runs-on: ubuntu-latest

    services:
      postgres:
        image: postgres:17-alpine
        env:
          POSTGRES_PASSWORD: postgres
          POSTGRES_DB: sachapel_test
        options: >-
          --health-cmd pg_isready
          --health-interval 10s
          --health-timeout 5s
          --health-retries 5

    steps:
      - uses: actions/checkout@v4

      - name: Set up Go
        uses: actions/setup-go@v5
        with:
          go-version: "1.22"

      - name: Install Templ
        run: go install github.com/a-h/templ/cmd/templ@latest

      - name: Generate templates
        run: templ generate

      - name: Install dependencies
        run: go mod download

      - name: Run tests
        run: go test -v -race -coverprofile=coverage.txt ./...
        env:
          DATABASE_URL: postgres://postgres:postgres@localhost:5432/sachapel_test?sslmode=disable

      - name: Upload coverage
        uses: codecov/codecov-action@v3
        with:
          file: ./coverage.txt

  lint:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: "1.22"
      - name: golangci-lint
        uses: golangci/golangci-lint-action@v4
```

#### .github/workflows/deploy-staging.yml

```yaml
name: Deploy to Staging

on:
  push:
    branches: [develop]

jobs:
  deploy:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - name: Set up Docker Buildx
        uses: docker/setup-buildx-action@v3

      - name: Build Docker image
        run: docker build -t sachapel-staging .

      - name: Deploy to staging
        uses: appleboy/ssh-action@v1.0.0
        with:
          host: ${{ secrets.STAGING_HOST }}
          username: ${{ secrets.STAGING_USER }}
          key: ${{ secrets.STAGING_SSH_KEY }}
          script: |
            cd /opt/sachapel-staging
            git pull origin develop
            docker compose -f compose.yml -f compose.staging.yml up -d --build
            docker compose exec -T app ./sachapel migrate

      - name: Smoke test
        run: |
          sleep 10
          curl --fail https://staging.sachapel.com || exit 1
```

#### .github/workflows/deploy-production.yml

```yaml
name: Deploy to Production

on:
  push:
    branches: [main]

jobs:
  deploy:
    runs-on: ubuntu-latest
    environment: production
    steps:
      - uses: actions/checkout@v4

      - name: Create release tag
        run: |
          VERSION=$(date +%Y.%m.%d-%H%M)
          git tag -a "v$VERSION" -m "Release $VERSION"
          git push origin "v$VERSION"

      - name: Deploy to production
        uses: appleboy/ssh-action@v1.0.0
        with:
          host: ${{ secrets.PRODUCTION_HOST }}
          username: ${{ secrets.PRODUCTION_USER }}
          key: ${{ secrets.PRODUCTION_SSH_KEY }}
          script: |
            cd /opt/sachapel
            git pull origin main
            docker compose -f compose.yml -f compose.prod.yml up -d --build
            docker compose exec -T app ./sachapel migrate

      - name: Verify deployment
        run: |
          sleep 15
          curl --fail https://sachapel.com || exit 1

      - name: Notify deployment
        run: echo "Deployment successful"
```

### Database Migrations

#### Creating Migrations

```bash
# Create new migration
make migrate-create name=add_events_table

# This creates two files:
# migrations/000001_add_events_table.up.sql
# migrations/000001_add_events_table.down.sql
```

#### Migration Files

```sql
-- 000001_add_events_table.up.sql
CREATE TABLE events (
    id BIGSERIAL PRIMARY KEY,
    title VARCHAR(255) NOT NULL,
    -- ... rest of schema
);

-- 000001_add_events_table.down.sql
DROP TABLE IF EXISTS events;
```

#### Running Migrations

```bash
# Apply all pending migrations
make migrate-up

# Rollback last migration
make migrate-down

# Check migration status
migrate -path ./migrations -database "$DATABASE_URL" version
```

### Code Organization

```
sachapel/
├── cmd/
│   ├── server/
│   │   └── main.go              # Application entry point
│   └── seed/
│       └── main.go              # Database seeding
├── internal/
│   ├── config/
│   │   └── config.go            # Configuration management
│   ├── database/
│   │   └── database.go          # Database connection
│   ├── models/
│   │   ├── user.go
│   │   ├── event.go
│   │   ├── bulletin.go
│   │   └── ...
│   ├── handlers/
│   │   ├── auth.go
│   │   ├── events.go
│   │   ├── bulletins.go
│   │   └── ...
│   ├── middleware/
│   │   ├── auth.go
│   │   ├── rate_limit.go
│   │   └── logging.go
│   ├── services/
│   │   ├── email.go
│   │   ├── auth.go
│   │   └── storage.go
│   └── utils/
│       ├── validator.go
│       └── helpers.go
├── migrations/
│   ├── 000001_create_users.up.sql
│   ├── 000001_create_users.down.sql
│   └── ...
├── static/
│   ├── css/
│   │   ├── base.css
│   │   ├── layout.css
│   │   ├── components.css
│   │   ├── utilities.css
│   │   └── print.css
│   ├── js/
│   │   ├── htmx.min.js
│   │   └── alpine.min.js
│   └── images/
│       └── arch-logo.svg
├── templates/
│   ├── layouts/
│   │   ├── base.templ
│   │   └── member.templ
│   ├── pages/
│   │   ├── home.templ
│   │   ├── about.templ
│   │   └── ...
│   ├── components/
│   │   ├── nav.templ
│   │   ├── footer.templ
│   │   ├── button.templ
│   │   └── card.templ
│   └── errors/
│       ├── 404.templ
│       ├── 403.templ
│       ├── 500.templ
│       └── maintenance.templ
├── compose.yml
├── compose.dev.yml
├── compose.staging.yml
├── compose.prod.yml
├── Dockerfile
├── Makefile
├── go.mod
├── go.sum
├── .env.example
├── .gitignore
└── README.md
```

---

## Security Considerations

### Application Security

#### Input Validation

- All user inputs validated server-side
- GORM ORM prevents SQL injection
- HTML template escaping prevents XSS
- File upload validation (type, size, content)
- CSRF tokens on all state-changing forms

#### Password Security

- bcrypt hashing with cost factor 12
- Minimum password requirements:
  - 8+ characters
  - At least one uppercase letter
  - At least one lowercase letter
  - At least one number
  - At least one special character
- Password strength meter on registration
- No password hints or recovery questions

#### Session Security

- JWT tokens with 24-hour expiration
- HTTP-only, Secure, SameSite=Strict cookies
- Session invalidation on logout
- Concurrent session detection (optional)
- CSRF protection on all mutations

#### Rate Limiting

- Login: 5 attempts per 15 minutes per IP
- Password reset: 3 requests per hour per email
- Form submission: 10 per hour per IP
- API endpoints: 60 requests per minute per IP
- Account lockout: 15 minutes after 5 failed logins

#### File Upload Security

- Whitelist allowed extensions (.pdf, .jpg, .png)
- MIME type validation
- File size limits (10MB max)
- Virus scanning (ClamAV in Phase 2)
- Files stored outside web root
- Nginx X-Accel-Redirect for protected files

### Infrastructure Security

#### Server Hardening

- Automatic security updates enabled
- Firewall (UFW) configured:
  - Allow 22 (SSH, key-only)
  - Allow 80 (HTTP)
  - Allow 443 (HTTPS)
  - Deny all other inbound
- Fail2ban for SSH brute force protection
- Non-root user for application
- Docker containers run as non-root user

#### SSL/TLS Configuration

- Let's Encrypt certificates (automated renewal)
- TLS 1.2 and 1.3 only
- Strong cipher suites
- HSTS with includeSubDomains
- OCSP stapling
- Certificate pinning (optional)

#### Database Security

- PostgreSQL not exposed to public internet
- Strong passwords for database users
- Encrypted connections between app and database
- Regular automated backups
- Point-in-time recovery capability

#### Environment Variables

- Secrets stored in .env files (not committed)
- Production secrets managed via server environment
- Rotation schedule for JWT secret (annually)
- Different secrets per environment

#### Security Headers

```
Strict-Transport-Security: max-age=31536000; includeSubDomains
X-Frame-Options: SAMEORIGIN
X-Content-Type-Options: nosniff
X-XSS-Protection: 1; mode=block
Referrer-Policy: strict-origin-when-cross-origin
Content-Security-Policy: default-src 'self'; script-src 'self' 'unsafe-inline'; style-src 'self' 'unsafe-inline'
```

#### CORS Policy

- **Default:** Same-origin only (no CORS headers needed)
- **API endpoints:** Restrict to same origin for MVP
- **Future consideration:** If mobile app or third-party integrations needed, implement allowlist-based CORS with:
  - Explicit allowed origins
  - Credentials support for authenticated requests
  - Preflight caching (Access-Control-Max-Age)

#### CSRF Protection with HTMX

**Implementation:**

1. Server generates CSRF token and stores in session
2. Token rendered in page `<meta>` tag and HTMX headers
3. HTMX configured globally to include token in all requests

**Base Template Setup (Templ):**

```go
// templates/layouts/base.templ

templ Base(csrfToken string) {
  <!DOCTYPE html>
  <html>
    <head>
      <meta name="csrf-token" content={ csrfToken }/>
    </head>
    <body hx-headers={ fmt.Sprintf(`{"X-CSRF-Token": "%s"}`, csrfToken) }>
      { children... }
    </body>
  </html>
}
```

**Server Validation:**

- Middleware checks `X-CSRF-Token` header on POST/PUT/PATCH/DELETE
- Token compared against session-stored value
- Mismatch returns 403 Forbidden
- GET/HEAD/OPTIONS exempt (safe methods)

### Monitoring & Logging

#### Application Logging

- Structured JSON logs
- Log levels: DEBUG, INFO, WARN, ERROR
- PII excluded from logs
- Failed login attempts logged
- All administrative actions logged
- Log rotation (7-day retention)

#### Security Monitoring

- Failed login attempt tracking
- Unusual activity detection (many failed logins, rapid requests)
- File upload anomalies
- Database query monitoring
- Alert on repeated security events

---

## Backup & Disaster Recovery

### Backup Strategy

#### Daily Automated Backups

**Database Backup**

```bash
#!/bin/bash
# /opt/scripts/backup-database.sh

DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/opt/backups/database"
RETENTION_DAYS=30

# Create backup
docker compose exec -T postgres pg_dump -U sacuser sachapel | gzip > "$BACKUP_DIR/sachapel_$DATE.sql.gz"

# Delete backups older than retention period
find "$BACKUP_DIR" -name "sachapel_*.sql.gz" -mtime +$RETENTION_DAYS -delete

# Upload to remote storage (optional)
# rclone copy "$BACKUP_DIR/sachapel_$DATE.sql.gz" remote:sachapel-backups/database/
```

**File System Backup**

```bash
#!/bin/bash
# /opt/scripts/backup-files.sh

DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/opt/backups/files"
SOURCE_DIR="/var/lib/docker/volumes/sachapel_app_uploads/_data"
RETENTION_DAYS=30

# Create tarball of uploads
tar -czf "$BACKUP_DIR/uploads_$DATE.tar.gz" -C "$SOURCE_DIR" .

# Delete old backups
find "$BACKUP_DIR" -name "uploads_*.tar.gz" -mtime +$RETENTION_DAYS -delete

# Upload to remote storage (optional)
# rclone copy "$BACKUP_DIR/uploads_$DATE.tar.gz" remote:sachapel-backups/files/
```

**Bulletin Cleanup Script**

```bash
#!/bin/bash
# /opt/scripts/cleanup-bulletins.sh

BULLETIN_DIR="/var/lib/docker/volumes/sachapel_app_uploads/_data/bulletins"
RETENTION_DAYS=14

# Delete bulletin files older than retention period
find "$BULLETIN_DIR" -name "*.pdf" -mtime +$RETENTION_DAYS -delete

# Also clean up database records (run via app CLI)
docker compose exec -T app ./sachapel cleanup-bulletins --days=$RETENTION_DAYS
```

**Cron Schedule**

```cron
# Daily at 2 AM - Backups
0 2 * * * /opt/scripts/backup-database.sh
15 2 * * * /opt/scripts/backup-files.sh

# Daily at 3 AM - Bulletin cleanup
0 3 * * * /opt/scripts/cleanup-bulletins.sh

# Daily at 4 AM - Audit log cleanup (90-day retention)
0 4 * * * docker compose exec -T app ./sachapel cleanup-audit-log --days=90

# Weekly full system backup (Sunday 5 AM)
0 5 * * 0 /opt/scripts/backup-full.sh
```

#### Backup Verification

- Weekly automated restore test to staging environment
- Checksum verification of backup files
- Backup size monitoring (alert on unexpected changes)

### Disaster Recovery Plan

#### Recovery Time Objective (RTO)

- **Target:** 4 hours
- Time to restore service after complete failure

#### Recovery Point Objective (RPO)

- **Target:** 24 hours
- Maximum acceptable data loss

#### Recovery Procedure

**Complete Server Failure:**

1. Provision new VPS (same specs)
2. Install Docker and Docker Compose
3. Clone repository
4. Restore environment variables
5. Restore latest database backup
6. Restore latest file backup
7. Start Docker containers
8. Verify functionality
9. Update DNS if IP changed

**Database Corruption:**

1. Stop application container
2. Restore database from latest backup
3. Verify data integrity
4. Restart application
5. Test critical functions

**File System Loss:**

1. Restore files from latest backup
2. Verify file integrity
3. Update file permissions
4. Test file access

#### Testing Schedule

- Quarterly disaster recovery drill
- Annual full failover test
- Document lessons learned

---

## Content Management

### Content Update Workflows

#### Weekly Bulletin Upload

1. Staff member logs in with staff credentials
2. Navigate to `/staff/bulletins/upload`
3. Select bulletin type (morning/evening)
4. Select Sunday date
5. Upload PDF file (max 10MB)
6. System validates file, stores in `/storage/bulletins/{type}/`
7. Old bulletins (>2 weeks) automatically deleted by daily cron job
8. Bulletin appears on `/resources/bulletins`

#### Event Creation

1. Staff member navigates to `/staff/events/create`
2. Fill out form:
   - Title, description
   - Date/time (start and optional end)
   - Location (on-site or off-site with address)
   - Recurring event settings (if applicable)
   - Visibility schedule (when to show/hide)
   - Registration settings (enable/disable, capacity limit, deadline)
   - Public or members-only
   - Ministry association
3. Submit form
4. Event auto-publishes based on visibility schedule
5. Event auto-hides after visibility end date

#### Event Recurrence (MVP)

**Simple Patterns Only:**

- `none` - Single occurrence (default)
- `daily` - Every day
- `weekly` - Same day each week
- `biweekly` - Every two weeks
- `monthly` - Same date each month
- `yearly` - Same date each year

**Implementation:**

- Store pattern type in `recurrence_rule` column
- Store end date in `recurrence_end` column
- Generate occurrences on-the-fly when displaying calendar
- No exception dates for MVP

**Phase 2 Enhancement:**

- Full RRULE (RFC 5545) support using teambition/rrule-go library
- Exception dates for skipping occurrences
- Single occurrence modifications

#### Announcement Posting

1. Staff member navigates to `/staff/announcements/create`
2. Write announcement:
   - Title
   - Content (rich text editor)
   - Visibility dates (from/until)
3. Save or publish immediately
4. Announcement displays on homepage
5. Auto-expires based on visibility until date

#### Ministry Page Updates

1. Ministry leader logs in
2. Navigate to `/staff/ministry/{ministry-slug}`
3. Edit page content (WYSIWYG editor)
4. Update meeting times, contact info
5. Save changes
6. Changes immediately reflected on public site

### Content Scheduling

**Auto-Publish/Expire Logic:**

- Content with `visible_from` in future: hidden until that date/time
- Content with `visible_until` in past: automatically hidden
- Daily cron job (6 AM) processes visibility changes
- Manual override available for admin

### Rich Text Editor

**Quill (Lightweight WYSIWYG)**

**Rationale:**

- Lightweight (~40KB vs 300KB+ for alternatives)
- BSD licensed (no restrictions)
- Simple API, easy integration with HTMX
- Sufficient for announcements and ministry content
- HTML output easy to sanitize server-side

**Features enabled:**

- Headers (H2, H3)
- Bold, italic, underline
- Ordered and unordered lists
- Links
- Image embedding (uploaded to server)

**Server-side sanitization:**

- Use bluemonday (Go HTML sanitizer)
- Whitelist allowed tags and attributes
- Strip all JavaScript and event handlers

---

## Third-Party Integrations

### Current External Services

#### Breeze ChMS (Giving)

- **Integration:** Simple redirect or iframe embed
- **URL:** https://sachapel.breezechms.com/give/online
- **Implementation:** Link on `/give` page
- **Future:** Replace with self-hosted payment processing (Phase 3)

#### Google Calendar (Events)

- **Integration:** Embed iframe or API integration
- **Current:** Embedded Google Calendar
- **Future:** Migrate to self-hosted calendar with iCal export
- **Phase 2:** Two-way sync (create events in app, sync to Google)

#### Microsoft 365 (Email)

- **Integration:** SMTP relay for transactional emails
- **Configuration:**
  - Host: smtp.office365.com
  - Port: 587 (STARTTLS)
  - Authentication: noreply@sachapel.com
- **Usage:** Account verification, password resets, form notifications

### Services to Replace

#### Wufoo (Forms) → Self-Hosted Forms

- **Replacement:** Custom form builder (Phase 2)
- **Timeline:** Phase 2 completion
- **Cost Savings:** ~$20/month

#### MailChimp (Newsletters) → Self-Hosted Email

- **Replacement:** Custom newsletter system (Phase 3)
- **Timeline:** Phase 3
- **Cost Savings:** ~$30/month

#### Vimeo (Sermon Streaming) → Self-Hosted Video

- **Replacement:** Self-hosted video with HLS streaming (Phase 3)
- **Timeline:** Phase 3
- **Cost Savings:** ~$75/month
- **Considerations:** Increased bandwidth usage, CDN may be needed

### APIs to Consider

#### Address Validation (Optional)

- **Purpose:** Validate member addresses
- **Options:** USPS API (free), SmartyStreets
- **Timeline:** Phase 2

#### Maps Integration

- **Current:** Google Maps link
- **Future:** OpenStreetMap for self-hosted alternative
- **Timeline:** Phase 3

---

## Performance & Scalability

### Performance Targets

- **Time to First Byte (TTFB):** < 200ms
- **First Contentful Paint (FCP):** < 1.5s
- **Largest Contentful Paint (LCP):** < 2.5s
- **Time to Interactive (TTI):** < 3.5s
- **Cumulative Layout Shift (CLS):** < 0.1

### Optimization Strategies

#### Backend Optimization

- Efficient database queries with proper indexes
- Connection pooling (GORM default)
- Redis caching for frequently accessed data
- Lazy loading of images and PDFs
- GZIP compression in Nginx

#### Frontend Optimization

- Minified CSS and JavaScript
- Image optimization (WebP format with fallbacks)
- Lazy loading images below the fold
- Preload critical resources
- HTTP/2 server push for critical assets

#### Caching Strategy

- **Static assets:** 1 year cache (versioned filenames)
- **HTML pages:** No cache for dynamic, 5 minutes for static
- **API responses:** Vary by endpoint (events: 5 min, bulletins: 1 hour)
- **Redis cache:** Page fragments, session data

#### Database Optimization

- Indexes on frequently queried columns
- Materialized views for complex reports (Phase 2)
- Query optimization with EXPLAIN ANALYZE
- Connection pooling (max 25 connections)
- Read replicas (if needed in Phase 3)

### Scalability Considerations

#### Current Architecture (Single VPS)

- **Capacity:** 500-1000 concurrent users
- **Monthly traffic:** ~50,000 page views
- **Storage:** 40GB SSD (upgradable)

#### Scaling Path (If Needed)

1. **Vertical scaling:** Upgrade VPS resources
2. **Caching:** Add CDN (Cloudflare free tier)
3. **Database:** Separate database server
4. **Load balancing:** Add second app server + Nginx load balancer
5. **Object storage:** Move files to S3-compatible storage

#### Monitoring Metrics

- Response times (p50, p95, p99)
- Error rates
- Database query performance
- Memory and CPU usage
- Disk I/O
- Network bandwidth

---

## Accessibility Requirements

### WCAG 2.1 AA Compliance

#### Visual Accessibility

- Color contrast ratio: minimum 4.5:1 for normal text, 3:1 for large text
- Text resizable up to 200% without loss of functionality
- No reliance on color alone to convey information
- Focus indicators clearly visible on all interactive elements

#### Keyboard Navigation

- All interactive elements accessible via keyboard
- Logical tab order
- Skip navigation link
- Keyboard shortcuts documented

#### Screen Reader Support

- Semantic HTML5 elements
- ARIA labels where needed
- Alt text for all images
- Form labels properly associated
- Heading hierarchy (H1 → H2 → H3)

#### Content Structure

- Clear heading hierarchy
- Lists for grouped content
- Tables with proper headers
- Descriptive link text (no "click here")

#### Forms

- Associated labels for all inputs
- Clear error messages
- Inline validation feedback
- Error summary at top of form

#### Media

- Captions for videos (future sermon archives)
- Transcripts for audio content
- Audio descriptions where applicable

### Testing

- Automated testing with axe-core
- Manual testing with screen readers (NVDA, JAWS, VoiceOver)
- Keyboard-only navigation testing
- Color contrast validation

---

## Future Enhancements

### Phase 3 Features (Beyond MVP)

#### Advanced Member Features

- Member photo directory with search/filter
- Direct messaging between members (opt-in)
- Member-to-member prayer requests
- Small group discussion forums
- Personalized dashboards

#### Self-Hosted Sermon Management

- Video upload and transcoding
- HLS streaming for adaptive bitrate
- Sermon notes and slides synchronized with video
- Search by speaker, date, scripture, topic
- Playlists and series
- Podcast feed generation

#### Saint Andrew's Conservatory Integration

- Separate section of site: `/conservatory/`
- Music lesson scheduling
- Student progress tracking
- Recital calendar and registrations
- Teacher profiles and bios
- Tuition payment processing

#### Advanced Form Builder

- Drag-and-drop form designer
- Conditional logic (show field based on previous answer)
- Multi-page forms
- File uploads in forms
- Payment integration (Stripe)
- Form analytics (completion rates, drop-off points)

#### Email Newsletter System

- Visual email template builder
- Subscriber management with opt-in/opt-out
- Segmentation (send to specific groups)
- Scheduled sending
- Analytics (open rates, click rates)
- A/B testing

#### Advanced Reporting

- Member engagement metrics
- Event attendance trends
- Form submission analytics
- Financial giving reports (Breeze integration)
- Custom report builder

#### Mobile App (PWA)

- Progressive Web App for offline access
- Push notifications for events and announcements
- Native install prompt
- Offline bulletin access
- Calendar sync to device calendar

#### API for Third-Party Integrations

- Public REST API for events, bulletins
- Webhook support for external systems
- OAuth 2.0 for secure integrations
- Rate limiting and API keys
- Developer documentation

---

## Decision Log

### Major Architectural Decisions

#### 1. Go + Chi Router Instead of Other Frameworks

**Date:** 2026-02-05  
**Decision:** Use Go with Chi router for backend  
**Rationale:**

- Built on standard library net/http (minimal abstraction)
- Single binary deployment (no runtime dependencies)
- Excellent AI/Claude Code assistance (extensive training data)
- Better performance and lower resource usage on VPS
- Strong type safety reduces bugs
- Idiomatic Go with standard patterns
- Future-proof: easy migration to pure net/http if needed
- Active maintenance with stable API

**Alternatives Considered:**

- Fiber: Good performance but uses fasthttp (less compatible), more abstractions
- Echo: Good option but Chi is more minimal and standard
- Gin: Less active maintenance, more "magical" API
- net/http alone: Too much routing boilerplate
- Node.js/Express: Requires Node runtime, more dependencies
- Python/FastAPI: Requires Python runtime, slower performance

---

#### 2. Templ + Vanilla CSS vs html/template + Tailwind

**Date:** 2026-02-05  
**Decision:** Use Templ for templates and vanilla CSS for styling  
**Rationale:**

- **Templ**: Type-safe templates with compile-time error checking
- **Templ**: Component-based architecture, better IDE support
- **Templ**: Excellent AI/Claude Code generation support
- **Vanilla CSS**: Zero Node.js dependency (aligns with self-contained philosophy)
- **Vanilla CSS**: Full control over styling with maintainable design system
- **Vanilla CSS**: Excellent AI code generation, easy for any developer
- **Both**: Better SEO and initial page load than SPA
- **Both**: Progressive enhancement for accessibility
- **Both**: Simpler architecture, less tooling complexity

**Alternatives Considered:**

- html/template: No type safety, harder to refactor, runtime errors
- Tailwind CSS: Requires Node.js, build tools, another dependency
- React/Vue SPA: Overkill for content site, worse SEO, complex tooling

---

#### 3. Self-Hosted vs External Services

**Date:** 2026-02-05  
**Decision:** Minimize external dependencies where feasible  
**Rationale:**

- Long-term sustainability (services can disappear)
- Cost control (eliminate monthly subscriptions)
- Future-proofing (admin may not be available for urgent fixes)
- Data ownership and privacy

**Exceptions:**

- Microsoft 365 SMTP (already paid for, reliable)
- Let's Encrypt (industry standard, free)
- Potential fallback email service (Resend/SES)

---

#### 4. PostgreSQL vs SQLite

**Date:** 2026-02-05  
**Decision:** PostgreSQL for production database  
**Rationale:**

- Better concurrency for multiple staff updating content
- Robust user/permission management
- JSONB for flexible form data storage
- Full-text search capabilities
- Better suited for 1000+ user database

**Alternatives Considered:**

- SQLite: Simpler but limited concurrency, less suitable for web app

---

#### 5. Docker Volumes vs S3 for File Storage

**Date:** 2026-02-05  
**Decision:** Docker volumes on VPS for MVP  
**Rationale:**

- Simpler architecture (fewer external dependencies)
- Lower cost (no S3 fees)
- Easier backup with simple scripts
- Sufficient for expected file volume

**Future Consideration:**

- Migrate to S3-compatible storage (Backblaze B2, Wasabi) if file volume grows significantly

---

#### 6. Custom Auth vs OAuth/Auth0

**Date:** 2026-02-05  
**Decision:** Custom JWT-based authentication  
**Rationale:**

- No external service dependency
- Full control over user data
- No ongoing costs
- Simpler for church members (email/password vs OAuth flow)

**Trade-offs:**

- More development work upfront
- Responsibility for security implementation

---

#### 7. RBAC vs Simple Role Hierarchy

**Date:** 2026-02-05  
**Decision:** Use Role-Based Access Control (RBAC) instead of hierarchical roles  
**Rationale:**

- Musicians may not be church members but need specific access
- Volunteers need limited access without full member privileges
- Users can have multiple roles (e.g., member + volunteer + musician)
- More flexible as church organizational needs evolve
- Prevents permission escalation issues
- Better separation of concerns

**Implementation:**

- Users can have multiple roles simultaneously
- Each role grants specific permissions
- Permissions are checked individually, not inherited
- Admin role has access to all resources
- Role combinations are additive

**Alternatives Considered:**

- Simple hierarchy (public → member → staff → admin): Too rigid, doesn't accommodate musicians who aren't members or volunteers with limited access

---

#### 8. Docker Compose v2 Convention

**Date:** 2026-02-05  
**Decision:** Use Docker Compose v2 (`docker compose`) instead of v1 (`docker-compose`)  
**Rationale:**

- Docker Compose v2 is now integrated into Docker CLI
- `docker compose` (without hyphen) is the current standard
- File naming convention changed: `compose.yml` instead of `docker-compose.yml`
- Better integration with Docker ecosystem
- Improved performance and features
- V1 is deprecated and will be removed

**Implementation:**

- Use `docker compose` command (no hyphen)
- Name files: `compose.yml`, `compose.dev.yml`, `compose.staging.yml`, `compose.prod.yml`
- All documentation and scripts use v2 conventions

---

## Appendix

### Glossary

- **MVP:** Minimum Viable Product - initial release with core features
- **VPS:** Virtual Private Server - cloud-based server instance
- **JWT:** JSON Web Token - token-based authentication standard
- **RBAC:** Role-Based Access Control - authorization model using roles and permissions
- **Chi:** Lightweight router built on Go's net/http standard library
- **Templ:** Type-safe templating language for Go that compiles to Go code
- **HTMX:** Library for AJAX, CSS Transitions, WebSockets with HTML attributes
- **GORM:** Go Object-Relational Mapping library
- **WCAG:** Web Content Accessibility Guidelines
- **TLS:** Transport Layer Security (SSL successor)
- **CSRF:** Cross-Site Request Forgery attack
- **XSS:** Cross-Site Scripting attack

### Reference Links

- **Go Documentation:** https://go.dev/doc/
- **Chi Router:** https://github.com/go-chi/chi
- **Templ:** https://templ.guide/
- **HTMX:** https://htmx.org/
- **Alpine.js:** https://alpinejs.dev/
- **PostgreSQL:** https://www.postgresql.org/docs/
- **Redis:** https://redis.io/docs/
- **Docker:** https://docs.docker.com/
- **GORM:** https://gorm.io/docs/
- **WCAG Guidelines:** https://www.w3.org/WAI/WCAG21/quickref/
- **Let's Encrypt:** https://letsencrypt.org/docs/
- **MDN CSS Reference:** https://developer.mozilla.org/en-US/docs/Web/CSS

### Contact Information

**Project Lead:** Steven, Church Administrator  
**Church:** Saint Andrew's Chapel, Sanford, FL  
**Repository:** https://github.com/sfdeloach/churchsite

---

**END OF SPECIFICATION**

_Last Updated: February 5, 2026_  
_Version: 1.2_
