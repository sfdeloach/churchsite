package handlers

import (
	"log/slog"
	"net/http"

	"github.com/sfdeloach/churchsite/internal/services"
	"github.com/sfdeloach/churchsite/templates/pages"
)

// AboutHandler handles about section pages.
type AboutHandler struct {
	siteSettings *services.SiteSettingsService
	staffMembers *services.StaffMemberService
}

// NewAboutHandler creates a new AboutHandler.
func NewAboutHandler(siteSettings *services.SiteSettingsService, staffMembers *services.StaffMemberService) *AboutHandler {
	return &AboutHandler{
		siteSettings: siteSettings,
		staffMembers: staffMembers,
	}
}

// Index redirects /about to /about/history.
func (h *AboutHandler) Index(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "/about/history", http.StatusMovedPermanently)
}

// History renders the church history page.
func (h *AboutHandler) History(w http.ResponseWriter, r *http.Request) {
	settings, err := h.siteSettings.GetAll()
	if err != nil {
		slog.Error("failed to load site settings", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	component := pages.AboutHistory(settings)
	if err := component.Render(r.Context(), w); err != nil {
		slog.Error("failed to render about history page", "error", err)
	}
}

// Beliefs renders the doctrine and beliefs page.
func (h *AboutHandler) Beliefs(w http.ResponseWriter, r *http.Request) {
	settings, err := h.siteSettings.GetAll()
	if err != nil {
		slog.Error("failed to load site settings", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	component := pages.AboutBeliefs(settings)
	if err := component.Render(r.Context(), w); err != nil {
		slog.Error("failed to render about beliefs page", "error", err)
	}
}

// Worship renders the theology of worship page.
func (h *AboutHandler) Worship(w http.ResponseWriter, r *http.Request) {
	settings, err := h.siteSettings.GetAll()
	if err != nil {
		slog.Error("failed to load site settings", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	component := pages.AboutWorship(settings)
	if err := component.Render(r.Context(), w); err != nil {
		slog.Error("failed to render about worship page", "error", err)
	}
}

// Gospel renders the gospel explanation page.
func (h *AboutHandler) Gospel(w http.ResponseWriter, r *http.Request) {
	settings, err := h.siteSettings.GetAll()
	if err != nil {
		slog.Error("failed to load site settings", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	component := pages.AboutGospel(settings)
	if err := component.Render(r.Context(), w); err != nil {
		slog.Error("failed to render about gospel page", "error", err)
	}
}

// Staff renders the pastors and staff page.
func (h *AboutHandler) Staff(w http.ResponseWriter, r *http.Request) {
	settings, err := h.siteSettings.GetAll()
	if err != nil {
		slog.Error("failed to load site settings", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	members, err := h.staffMembers.GetActive()
	if err != nil {
		slog.Error("failed to load staff members", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	component := pages.AboutStaff(settings, members)
	if err := component.Render(r.Context(), w); err != nil {
		slog.Error("failed to render about staff page", "error", err)
	}
}

// Building renders the building/place of worship page.
func (h *AboutHandler) Building(w http.ResponseWriter, r *http.Request) {
	settings, err := h.siteSettings.GetAll()
	if err != nil {
		slog.Error("failed to load site settings", "error", err)
		http.Error(w, "Internal Server Error", http.StatusInternalServerError)
		return
	}

	component := pages.AboutBuilding(settings)
	if err := component.Render(r.Context(), w); err != nil {
		slog.Error("failed to render about building page", "error", err)
	}
}
