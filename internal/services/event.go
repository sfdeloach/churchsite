package services

import (
	"time"

	"github.com/sfdeloach/churchsite/internal/models"
	"gorm.io/gorm"
)

// EventService handles event queries.
type EventService struct {
	db *gorm.DB
}

// NewEventService creates a new EventService.
func NewEventService(db *gorm.DB) *EventService {
	return &EventService{db: db}
}

// GetUpcoming returns upcoming public events ordered by date, limited to `limit` results.
// Only returns events that are public, not soft-deleted, and within their visibility window.
func (s *EventService) GetUpcoming(limit int) ([]models.Event, error) {
	var events []models.Event
	now := time.Now()

	err := s.db.
		Where("is_public = ? AND event_date >= ?", true, now).
		Where("(visible_from IS NULL OR visible_from <= ?)", now).
		Where("(visible_until IS NULL OR visible_until >= ?)", now).
		Order("event_date ASC").
		Limit(limit).
		Find(&events).Error

	return events, err
}
