package services

import (
	"github.com/sfdeloach/churchsite/internal/models"
	"gorm.io/gorm"
)

// SiteSettingsService handles site_settings queries.
type SiteSettingsService struct {
	db *gorm.DB
}

// NewSiteSettingsService creates a new SiteSettingsService.
func NewSiteSettingsService(db *gorm.DB) *SiteSettingsService {
	return &SiteSettingsService{db: db}
}

// GetAll returns all site settings as a key-value map.
func (s *SiteSettingsService) GetAll() (map[string]string, error) {
	var settings []models.SiteSetting
	if err := s.db.Find(&settings).Error; err != nil {
		return nil, err
	}

	result := make(map[string]string, len(settings))
	for _, setting := range settings {
		result[setting.Key] = setting.Value
	}
	return result, nil
}
