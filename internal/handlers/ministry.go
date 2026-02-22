package handlers

import (
	"errors"
	"log/slog"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/sfdeloach/churchsite/internal/services"
	"github.com/sfdeloach/churchsite/templates/pages"
	"gorm.io/gorm"
)

// MinistryHandler handles ministry pages.
type MinistryHandler struct {
	ministries *services.MinistryService
}

// NewMinistryHandler creates a new MinistryHandler.
func NewMinistryHandler(ministries *services.MinistryService) *MinistryHandler {
	return &MinistryHandler{ministries: ministries}
}

// Index renders the ministries overview page.
func (h *MinistryHandler) Index(w http.ResponseWriter, r *http.Request) {
	ministries, err := h.ministries.GetActive()
	if err != nil {
		slog.Error("failed to load ministries", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	component := pages.MinistriesIndex(ministries)
	if err := component.Render(r.Context(), w); err != nil {
		slog.Error("failed to render ministries index page", "error", err)
	}
}

// Show renders a single ministry detail page.
func (h *MinistryHandler) Show(w http.ResponseWriter, r *http.Request) {
	slug := chi.URLParam(r, "slug")

	ministry, err := h.ministries.GetBySlug(slug)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			http.Error(w, "Ministry not found", http.StatusNotFound)
			return
		}
		slog.Error("failed to load ministry", "slug", slug, "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	component := pages.MinistryShow(*ministry)
	if err := component.Render(r.Context(), w); err != nil {
		slog.Error("failed to render ministry show page", "slug", slug, "error", err)
	}
}
