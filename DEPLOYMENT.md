# Deployment & Operations

This document covers deployment configuration, CI/CD pipelines, backup procedures, and infrastructure details.

---

## Infrastructure Overview

```
┌─────────────────────────────────────────────────────────┐
│                    Internet / CDN                       │
└────────────────────┬────────────────────────────────────┘
                     │
                     ▼
┌─────────────────────────────────────────────────────────┐
│                  VPS (Single Server)                    │
│  ┌───────────────────────────────────────────────────┐  │
│  │  Docker Host (Ubuntu 22.04 LTS)                   │  │
│  │                                                   │  │
│  │  Nginx (Container) ─────► Go/Chi App (Container)  │  │
│  │    - Port 80/443              - Port 3000         │  │
│  │    - TLS (Let's Encrypt)           │              │  │
│  │    - Reverse Proxy                 ▼              │  │
│  │                          ┌─────────────────────┐  │  │
│  │                          │ PostgreSQL │ Redis  │  │  │
│  │                          │   5432     │  6379  │  │  │
│  │                          └─────────────────────┘  │  │
│  │                                    │              │  │
│  │                          Docker Volumes           │  │
│  └───────────────────────────────────────────────────┘  │
└─────────────────────────────────────────────────────────┘
```

### Production Environment

- **VPS:** Hetzner Cloud CX21 (2 vCPU, 4GB RAM, 40GB SSD, ~€5/month)
- **OS:** Ubuntu 22.04 LTS
- **Domain:** sachapel.com
- **SSL:** Let's Encrypt (automated via certbot)

### Staging Environment

- Same VPS, separate Docker Compose stack
- Ports: 8080 (HTTP), 8443 (HTTPS)
- Domain: staging.sachapel.com

---

## Docker Configuration

### compose.yml (Base)

```yaml
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
      - APP_URL=${APP_URL}
      - APP_PORT=${APP_PORT}
      - DATABASE_URL=${DATABASE_URL}
      - REDIS_URL=${REDIS_URL}
      - JWT_SECRET=${JWT_SECRET}
      - JWT_EXPIRATION=${JWT_EXPIRATION}
      - SMTP_HOST=${SMTP_HOST}
      - SMTP_PORT=${SMTP_PORT}
      - SMTP_USER=${SMTP_USER}
      - SMTP_PASS=${SMTP_PASS}
      - FROM_EMAIL=${FROM_EMAIL}
      - FROM_NAME=${FROM_NAME}
      - MAX_UPLOAD_SIZE=${MAX_UPLOAD_SIZE}
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

### compose.prod.yml

```yaml
services:
  app:
    build:
      target: production
    environment:
      - APP_ENV=production
```

### compose.staging.yml

```yaml
services:
  nginx:
    ports:
      - "8080:80"
      - "8443:443"

  app:
    environment:
      - APP_ENV=staging
```

---

## Dockerfile

```dockerfile
# Build stage
FROM golang:1.25-alpine AS builder

WORKDIR /build

RUN apk add --no-cache git make
RUN go install github.com/a-h/templ/cmd/templ@latest
RUN go install github.com/air-verse/air@latest

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN templ generate
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags="-w -s" -o /app/sachapel ./cmd/server

# Production stage
FROM alpine:latest AS production

WORKDIR /app

RUN apk add --no-cache ca-certificates tzdata

COPY --from=builder /app/sachapel .
COPY --from=builder /build/static ./static
COPY --from=builder /build/migrations ./migrations

RUN mkdir -p /app/storage/bulletins/morning \
             /app/storage/bulletins/evening \
             /app/storage/photos \
             /app/storage/documents \
             /app/storage/uploads

RUN addgroup -g 1000 appuser && \
    adduser -D -u 1000 -G appuser appuser && \
    chown -R appuser:appuser /app

USER appuser

EXPOSE 3000

CMD ["./sachapel"]
```

---

## Nginx Configuration

```nginx
# Rate limiting zones
limit_req_zone $binary_remote_addr zone=login_limit:10m rate=5r/m;
limit_req_zone $binary_remote_addr zone=api_limit:10m rate=60r/m;

# HTTP → HTTPS redirect
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

# HTTPS server
server {
    listen 443 ssl http2;
    server_name sachapel.com www.sachapel.com;

    ssl_certificate /etc/letsencrypt/live/sachapel.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/sachapel.com/privkey.pem;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers HIGH:!aNULL:!MD5;
    ssl_prefer_server_ciphers on;

    # Security headers
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    add_header X-Frame-Options "SAMEORIGIN" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header X-XSS-Protection "1; mode=block" always;
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;

    # Static files
    location /static/ {
        proxy_pass http://app:3000;
        expires 1y;
        add_header Cache-Control "public, immutable";
    }

    # Protected uploads (X-Accel-Redirect)
    location /uploads/ {
        internal;
        alias /var/www/uploads/;
    }

    # Rate-limited login
    location /login {
        limit_req zone=login_limit burst=2 nodelay;
        proxy_pass http://app:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # Rate-limited API
    location /api/ {
        limit_req zone=api_limit burst=10 nodelay;
        proxy_pass http://app:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # Default proxy
    location / {
        proxy_pass http://app:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
        proxy_http_version 1.1;
        proxy_set_header Upgrade $http_upgrade;
        proxy_set_header Connection "upgrade";
    }

    client_max_body_size 10M;
}
```

---

## CI/CD Pipelines

### Test & Lint (.github/workflows/test.yml)

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

### Deploy Staging (.github/workflows/deploy-staging.yml)

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

### Deploy Production (.github/workflows/deploy-production.yml)

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
```

---

## Backup Procedures

### Database Backup

```bash
#!/bin/bash
# /opt/scripts/backup-database.sh

DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/opt/backups/database"
RETENTION_DAYS=30

docker compose exec -T postgres pg_dump -U sacuser sachapel | gzip > "$BACKUP_DIR/sachapel_$DATE.sql.gz"

find "$BACKUP_DIR" -name "sachapel_*.sql.gz" -mtime +$RETENTION_DAYS -delete
```

### File Backup

```bash
#!/bin/bash
# /opt/scripts/backup-files.sh

DATE=$(date +%Y%m%d_%H%M%S)
BACKUP_DIR="/opt/backups/files"
SOURCE_DIR="/var/lib/docker/volumes/sachapel_app_uploads/_data"
RETENTION_DAYS=30

tar -czf "$BACKUP_DIR/uploads_$DATE.tar.gz" -C "$SOURCE_DIR" .

find "$BACKUP_DIR" -name "uploads_*.tar.gz" -mtime +$RETENTION_DAYS -delete
```

### Bulletin Cleanup

```bash
#!/bin/bash
# /opt/scripts/cleanup-bulletins.sh

# Atomic cleanup: queries DB, deletes files, removes records in transaction
docker compose exec -T app ./sachapel cleanup-bulletins --days=14
```

### Cron Schedule

```cron
# Daily backups at 2 AM
0 2 * * * /opt/scripts/backup-database.sh
15 2 * * * /opt/scripts/backup-files.sh

# Daily cleanup at 3-4 AM
0 3 * * * /opt/scripts/cleanup-bulletins.sh
0 4 * * * docker compose exec -T app ./sachapel cleanup-audit-log --days=90

# Weekly full backup (Sunday 5 AM)
0 5 * * 0 /opt/scripts/backup-full.sh
```

---

## Disaster Recovery

### Recovery Objectives

- **RTO (Recovery Time):** 4 hours
- **RPO (Recovery Point):** 24 hours

### Complete Server Failure

1. Provision new VPS (same specs)
2. Install Docker and Docker Compose
3. Clone repository
4. Restore environment variables
5. Restore latest database backup
6. Restore latest file backup
7. Start Docker containers
8. Verify functionality
9. Update DNS if IP changed

### Database Corruption

1. Stop application container
2. Restore database from latest backup
3. Verify data integrity
4. Restart application
5. Test critical functions

### Testing Schedule

- Quarterly disaster recovery drill
- Annual full failover test

---

## Security Hardening

### Server

- Automatic security updates
- UFW firewall: allow 22, 80, 443
- Fail2ban for SSH protection
- Non-root user for application
- Docker containers run as non-root

### SSL/TLS

- TLS 1.2 and 1.3 only
- Strong cipher suites
- HSTS with includeSubDomains
- Automated certificate renewal

### Database

- Not exposed to public internet
- Strong passwords
- Encrypted connections
- Regular automated backups

---

## Monitoring

### Application Logging

- Structured JSON logs (Go slog)
- Levels: DEBUG, INFO, WARN, ERROR
- PII excluded
- 7-day rotation

### Security Monitoring

- Failed login tracking
- Unusual activity detection
- File upload anomalies
- Alert on repeated security events

### Performance Metrics

- Response times (p50, p95, p99)
- Error rates
- Database query performance
- Memory/CPU usage
- Disk I/O

---

## Environment Variables (.env.production)

```bash
APP_ENV=production
APP_URL=https://sachapel.com
APP_PORT=3000

DATABASE_URL=postgres://sacuser:STRONG_PASSWORD@postgres:5432/sachapel?sslmode=disable
POSTGRES_DB=sachapel
POSTGRES_USER=sacuser
POSTGRES_PASSWORD=STRONG_PASSWORD

REDIS_URL=redis://redis:6379/0

JWT_SECRET=RANDOM_64_CHAR_STRING
JWT_EXPIRATION=24h

SMTP_HOST=smtp.office365.com
SMTP_PORT=587
SMTP_USER=noreply@sachapel.com
SMTP_PASS=SMTP_PASSWORD
FROM_EMAIL=noreply@sachapel.com
FROM_NAME=Saint Andrew's Chapel

MAX_UPLOAD_SIZE=10485760
LOGIN_RATE_LIMIT=5
LOGIN_RATE_WINDOW=15m
ACCOUNT_LOCKOUT_DURATION=15m
SESSION_DURATION=24h
COOKIE_SECURE=true
COOKIE_HTTPONLY=true
COOKIE_SAMESITE=strict
```
