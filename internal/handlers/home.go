package handlers

import (
	"log/slog"
	"net/http"

	"github.com/sfdeloach/churchsite/internal/services"
	"github.com/sfdeloach/churchsite/templates/pages"
)

// HomeHandler handles the homepage.
type HomeHandler struct {
	events *services.EventService
}

// NewHomeHandler creates a new HomeHandler.
func NewHomeHandler(events *services.EventService) *HomeHandler {
	return &HomeHandler{
		events: events,
	}
}

// Index renders the homepage with service times and upcoming events.
func (h *HomeHandler) Index(w http.ResponseWriter, r *http.Request) {
	events, err := h.events.GetUpcoming(6)
	if err != nil {
		slog.Error("failed to load upcoming events", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	component := pages.Home(events)
	if err := component.Render(r.Context(), w); err != nil {
		slog.Error("failed to render homepage", "error", err)
	}
}
