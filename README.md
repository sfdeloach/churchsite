# Saint Andrew's Chapel Website

A modern, self-contained church website built with Go, HTMX, and PostgreSQL.

## Quick Start

```bash
# Prerequisites: Docker, Go 1.22+, Make

# Clone and setup
git clone https://github.com/sfdeloach/churchsite.git
cd churchsite
cp .env.example .env.development

# Start services
make dev-up

# Run migrations and seed data
make migrate-up
make seed

# Access at http://localhost:3000
```

## Development

```bash
make watch          # Hot reload (templ + air)
make test           # Run tests
make lint           # Run linter
make generate       # Generate Templ templates
```

### Test Accounts (after seeding)

| Role  | Email               | Password      |
| ----- | ------------------- | ------------- |
| Admin | admin@sachapel.test | AdminPass123! |
| Staff | staff@sachapel.test | StaffPass123! |

## Documentation

- [SPEC.md](SPEC.md) - Technical specification (routes, schema, auth)
- [DECISIONS.md](DECISIONS.md) - Architectural decisions
- [DEPLOYMENT.md](DEPLOYMENT.md) - Docker, CI/CD, operations

## Tech Stack

- **Backend:** Go + Chi router
- **Templates:** Templ (type-safe)
- **Frontend:** HTMX + Alpine.js
- **Database:** PostgreSQL 17
- **Cache:** Redis 7

## Project Structure

```
├── cmd/server/          # Application entry point
├── internal/            # Application code
│   ├── handlers/        # HTTP handlers
│   ├── models/          # GORM models
│   ├── middleware/      # Auth, rate limiting
│   └── services/        # Business logic
├── migrations/          # SQL migrations
├── static/              # CSS, JS, images
└── templates/           # Templ templates
```

## License

This project is private and not licensed for external use.
