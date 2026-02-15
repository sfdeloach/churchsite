---
name: handler
description: Create HTTP handlers, services, and route registration following project patterns. Use when adding new routes or endpoints.
argument-hint: "[resource-name]"
disable-model-invocation: true
---

# Handler + Route Scaffolding

Create a complete handler for the `$ARGUMENTS` resource, including the service layer and route registration.

## Process

### 1. Create or Update the Service

File: `internal/services/<resource>.go`

```go
package services

import (
    "github.com/sfdeloach/churchsite/internal/models"
    "gorm.io/gorm"
)

type ResourceService struct {
    db *gorm.DB
}

func NewResourceService(db *gorm.DB) *ResourceService {
    return &ResourceService{db: db}
}

// Query methods here — each returns (result, error)
func (s *ResourceService) GetAll() ([]models.Resource, error) {
    var items []models.Resource
    err := s.db.Find(&items).Error
    return items, err
}
```

**Conventions:**
- One service per resource/domain
- Constructor: `New<Name>Service(db *gorm.DB)`
- All methods return `(result, error)`
- Use GORM query builder — never raw SQL in services
- Filter soft-deleted records automatically (GORM handles this for `gorm.Model` embeds)

### 2. Create the Handler

File: `internal/handlers/<resource>.go`

```go
package handlers

import (
    "log/slog"
    "net/http"

    "github.com/sfdeloach/churchsite/internal/services"
    "github.com/sfdeloach/churchsite/templates/pages"
)

type ResourceHandler struct {
    resourceSvc *services.ResourceService
}

func NewResourceHandler(resourceSvc *services.ResourceService) *ResourceHandler {
    return &ResourceHandler{
        resourceSvc: resourceSvc,
    }
}

// Index renders the resource listing page.
func (h *ResourceHandler) Index(w http.ResponseWriter, r *http.Request) {
    items, err := h.resourceSvc.GetAll()
    if err != nil {
        slog.Error("failed to load resources", "error", err)
        http.Error(w, "Internal Server Error", http.StatusInternalServerError)
        return
    }

    component := pages.ResourceIndex(items)
    if err := component.Render(r.Context(), w); err != nil {
        slog.Error("failed to render resource page", "error", err)
    }
}
```

**Conventions:**
- One handler struct per resource
- Constructor: `New<Name>Handler(deps...)` — accepts service pointers
- Methods have signature `(w http.ResponseWriter, r *http.Request)`
- Log errors with `slog.Error` including context fields
- For pages: render Templ component via `.Render(r.Context(), w)`
- For API endpoints: use `encoding/json` to write JSON responses
- For error pages: use `http.Error()` for now (error templates come later)
- Extract URL params with `chi.URLParam(r, "paramName")`

### 3. Register Routes

In `cmd/server/main.go`:

1. **Initialize the service** in the services section:
   ```go
   resourceSvc := services.NewResourceService(db.Postgres)
   ```

2. **Initialize the handler** in the handlers section:
   ```go
   resourceHandler := handlers.NewResourceHandler(resourceSvc)
   ```

3. **Register routes** in the appropriate section:
   ```go
   // Public routes
   r.Get("/resource", resourceHandler.Index)
   r.Get("/resource/{id}", resourceHandler.Show)

   // Authenticated routes (when middleware exists)
   r.Group(func(r chi.Router) {
       // r.Use(middleware.RequireAuth)
       r.Get("/staff/resource", resourceHandler.Manage)
       r.Post("/staff/resource/create", resourceHandler.Create)
   })
   ```

**Route conventions:**
- Public pages: `GET /resource`, `GET /resource/{slug}` or `GET /resource/{id}`
- Staff CRUD: `GET /staff/resource`, `POST /staff/resource/create`, `GET /staff/resource/{id}/edit`, `POST /staff/resource/{id}/edit`, `DELETE /staff/resource/{id}`
- API: `GET /api/v1/resource`, `POST /api/v1/resource`
- Use `{slug}` for user-facing URLs, `{id}` for admin/API URLs
- Group protected routes with `r.Group()` + middleware

### 4. Verify

After creating all files, remind the user to:
- Run `make generate` if new `.templ` files were created
- Run `make dev-up` to test the routes
- Check that routes appear in the Chi router debug output
