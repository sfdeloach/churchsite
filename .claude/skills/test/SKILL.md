---
name: test
description: Write and run Go tests for handlers, services, or models. Use when adding or running tests.
argument-hint: "[package-or-file]"
disable-model-invocation: true
allowed-tools: Bash, Read, Grep, Glob, Edit, Write
---

# Test Scaffolding and Runner

## Arguments

- `/test <package>` — Write tests for the specified package (e.g., `handlers`, `services/event`)
- `/test run` — Run the full test suite via `make test`
- `/test run <package>` — Run tests for a specific package

## Running Tests

When the argument is `run`:

```bash
# Full suite with race detection
make test

# Specific package
go test -v -race ./internal/handlers/ -run TestHandlerName

# With coverage
go test -v -race -coverprofile=coverage.txt ./...
```

## Writing Tests

### Handler Tests

File: `internal/handlers/<resource>_test.go`

```go
package handlers_test

import (
    "net/http"
    "net/http/httptest"
    "testing"

    "github.com/go-chi/chi/v5"
    "github.com/sfdeloach/churchsite/internal/handlers"
)

func TestResourceHandler_Index(t *testing.T) {
    // Setup
    // For handlers with service dependencies, use either:
    // 1. A test database with fixtures
    // 2. An interface-based mock

    handler := handlers.NewResourceHandler(svc)

    // Create request
    req := httptest.NewRequest(http.MethodGet, "/resource", nil)
    rec := httptest.NewRecorder()

    // Execute
    handler.Index(rec, req)

    // Assert
    if rec.Code != http.StatusOK {
        t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
    }
}
```

**For handlers that use `chi.URLParam`**, set up a Chi router in the test:

```go
func TestResourceHandler_Show(t *testing.T) {
    handler := handlers.NewResourceHandler(svc)

    r := chi.NewRouter()
    r.Get("/resource/{id}", handler.Show)

    req := httptest.NewRequest(http.MethodGet, "/resource/1", nil)
    rec := httptest.NewRecorder()

    r.ServeHTTP(rec, req)

    if rec.Code != http.StatusOK {
        t.Errorf("expected status %d, got %d", http.StatusOK, rec.Code)
    }
}
```

### Service Tests

File: `internal/services/<resource>_test.go`

Service tests need a real database connection. Use a test database (not mocks) for integration testing with GORM.

```go
package services_test

import (
    "testing"

    "github.com/sfdeloach/churchsite/internal/models"
    "github.com/sfdeloach/churchsite/internal/services"
    "gorm.io/driver/postgres"
    "gorm.io/gorm"
)

func setupTestDB(t *testing.T) *gorm.DB {
    t.Helper()
    dsn := "postgres://postgres:postgres@localhost:5432/sachapel_test?sslmode=disable"
    db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
    if err != nil {
        t.Fatalf("failed to connect to test database: %v", err)
    }
    return db
}

func TestEventService_GetUpcoming(t *testing.T) {
    db := setupTestDB(t)
    svc := services.NewEventService(db)

    events, err := svc.GetUpcoming(5)
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }

    // Assert on results
    if len(events) > 5 {
        t.Errorf("expected at most 5 events, got %d", len(events))
    }
}
```

### Test Conventions

- Use `_test` package suffix for black-box testing (e.g., `package handlers_test`)
- Use table-driven tests for multiple scenarios:

```go
func TestValidation(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        wantErr bool
    }{
        {"valid email", "user@example.com", false},
        {"empty email", "", true},
        {"no domain", "user@", true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            err := validate(tt.input)
            if (err != nil) != tt.wantErr {
                t.Errorf("validate(%q) error = %v, wantErr %v", tt.input, err, tt.wantErr)
            }
        })
    }
}
```

- Always run with `-race` flag
- Use `t.Helper()` in test helper functions
- Use `t.Fatalf()` for setup failures, `t.Errorf()` for assertion failures
- Use `t.Cleanup()` for teardown
- Name test functions: `Test<Type>_<Method>` (e.g., `TestEventService_GetUpcoming`)
