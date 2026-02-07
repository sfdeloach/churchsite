# Architectural Decisions

This document records major architectural decisions and their rationale. Reference this when onboarding new developers or revisiting past choices.

---

## 1. Go + Chi Router

**Decision:** Use Go with Chi router for backend

**Why:**

- Built on standard library net/http (minimal abstraction)
- Single binary deployment (no runtime dependencies)
- Excellent AI/Claude Code assistance (extensive training data)
- Better performance and lower resource usage on VPS
- Strong type safety reduces bugs
- Idiomatic Go with standard patterns
- Future-proof: easy migration to pure net/http if needed

**Alternatives Considered:**

- Fiber: Uses fasthttp (less compatible), more abstractions
- Echo: Good but Chi is more minimal
- Gin: Less active maintenance, more "magical" API
- net/http alone: Too much routing boilerplate
- Node.js/Express: Requires runtime, more dependencies
- Python/FastAPI: Slower performance, requires runtime

---

## 2. Templ + Vanilla CSS

**Decision:** Use Templ for templates and vanilla CSS for styling

**Why:**

- **Templ**: Type-safe templates with compile-time error checking
- **Templ**: Component-based architecture, better IDE support
- **Templ**: Excellent AI/Claude Code generation support
- **Vanilla CSS**: Zero Node.js dependency (aligns with self-contained philosophy)
- **Vanilla CSS**: Full control over styling with maintainable design system
- Better SEO and initial page load than SPA
- Progressive enhancement for accessibility
- Simpler architecture, less tooling complexity

**Alternatives Considered:**

- html/template: No type safety, runtime errors
- Tailwind CSS: Requires Node.js, build tools
- React/Vue SPA: Overkill for content site, worse SEO

---

## 3. Self-Hosted Philosophy

**Decision:** Minimize external dependencies where feasible

**Why:**

- Long-term sustainability (services can disappear)
- Cost control (eliminate monthly subscriptions)
- Future-proofing (admin may not be available for urgent fixes)
- Data ownership and privacy

**Exceptions:**

- Microsoft 365 SMTP (already paid for, reliable)
- Let's Encrypt (industry standard, free)
- Potential fallback email service (Resend/SES)

---

## 4. PostgreSQL over SQLite

**Decision:** PostgreSQL for production database

**Why:**

- Better concurrency for multiple staff updating content
- Robust user/permission management
- JSONB for flexible form data storage
- Full-text search capabilities
- Better suited for 1000+ user database

---

## 5. Docker Volumes for File Storage

**Decision:** Docker volumes on VPS for MVP

**Why:**

- Simpler architecture (fewer external dependencies)
- Lower cost (no S3 fees)
- Easier backup with simple scripts
- Sufficient for expected file volume

**Future Consideration:**

- Migrate to S3-compatible storage (Backblaze B2, Wasabi) if file volume grows

---

## 6. Custom Authentication

**Decision:** Custom JWT-based authentication instead of OAuth/Auth0

**Why:**

- No external service dependency
- Full control over user data
- No ongoing costs
- Simpler for church members (email/password vs OAuth flow)

**Trade-offs:**

- More development work upfront
- Responsibility for security implementation

---

## 7. Role-Based Access Control (RBAC)

**Decision:** Use RBAC instead of hierarchical roles

**Why:**

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

---

## 8. Docker Compose v2

**Decision:** Use Docker Compose v2 (`docker compose`) instead of v1

**Why:**

- V2 is integrated into Docker CLI
- `docker compose` (no hyphen) is the current standard
- File naming: `compose.yml` instead of `docker-compose.yml`
- Better integration with Docker ecosystem
- V1 is deprecated

---

## 9. HTMX + Alpine.js for Interactivity

**Decision:** Server-side rendering with progressive enhancement

**Why:**

- Minimal JavaScript bundle size
- Works without JavaScript (graceful degradation)
- Simpler mental model than SPA
- Better SEO out of the box
- Alpine.js handles complex client-side state when needed
- HTMX handles dynamic updates without full page reloads

---

## 10. JWT with Redis Blacklist

**Decision:** Hybrid JWT + Redis approach for sessions

**Why:**

- JWT: Stateless verification (no DB lookup per request)
- Redis blacklist: Enables immediate token revocation
- 24-hour token lifetime limits exposure
- Blacklist only checked when token would otherwise be valid

**Trade-offs:**

- More complex than pure sessions
- Requires Redis availability

---

## 11. Password Reset Token Hashing

**Decision:** Hash reset tokens with SHA-256 before storage

**Why:**

- Protects against database breach
- Attacker with DB access cannot use stored hashes to reset passwords
- Plain token sent via email, hash stored in DB
- Industry best practice

---

## 12. WebP Image Conversion

**Decision:** Convert all uploaded images to WebP

**Why:**

- Smaller file sizes (typically 25-35% smaller than JPEG)
- 97%+ browser support
- Single format simplifies serving logic
- Quality 85 provides good balance of size/quality

---

## Git Commit Convention Used

```
<type>[optional scope]: <description>
```

**Type**: A prefix indicating the nature of the change. Common types include:

- `feat`: A new feature.
- `deps`: A new project dependency or library.
- `fix`: A bug fix.
- `docs`: Documentation only changes.
- `style`: Changes that do not affect the meaning of the code (e.g., formatting, missing semi-colons).
- `refactor`: A code change that neither fixes a bug nor adds a feature.
- `perf`: A code change that improves performance.
- `test`: Adding missing tests or correcting existing tests.
- `build`: Changes that affect the build system or external dependencies.
- `ci`: Changes to CI configuration files and scripts.
- `chore`: Other changes that don't modify src or test files.
- `revert`: Reverts a previous commit.

**Scope** (optional): In parentheses after the type, specifies the part of the codebase affected (e.g., feat(parser): add ability to parse arrays).

**Description**: A short, imperative summary of the change (start with a verb like "add", "update", "remove"). Keep it under 50 characters if possible.

---

## Glossary

- **MVP:** Minimum Viable Product - initial release with core features
- **VPS:** Virtual Private Server - cloud-based server instance
- **JWT:** JSON Web Token - token-based authentication standard
- **RBAC:** Role-Based Access Control - authorization model using roles
- **Chi:** Lightweight router built on Go's net/http
- **Templ:** Type-safe templating language for Go
- **HTMX:** Library for AJAX with HTML attributes
- **GORM:** Go Object-Relational Mapping library
- **WCAG:** Web Content Accessibility Guidelines
