package services

import (
	"github.com/sfdeloach/churchsite/internal/models"
	"gorm.io/gorm"
)

// MinistryService handles ministry queries.
type MinistryService struct {
	db *gorm.DB
}

// NewMinistryService creates a new MinistryService.
func NewMinistryService(db *gorm.DB) *MinistryService {
	return &MinistryService{db: db}
}

// GetActive returns active, non-deleted ministries ordered by sort_order then name.
func (s *MinistryService) GetActive() ([]models.Ministry, error) {
	var ministries []models.Ministry

	err := s.db.
		Where("is_active = ?", true).
		Order("sort_order ASC, name ASC").
		Find(&ministries).Error

	return ministries, err
}

// GetBySlug returns a single active ministry by its slug.
// Returns gorm.ErrRecordNotFound if no active ministry with that slug exists.
func (s *MinistryService) GetBySlug(slug string) (*models.Ministry, error) {
	var ministry models.Ministry

	err := s.db.
		Where("slug = ? AND is_active = ?", slug, true).
		First(&ministry).Error

	if err != nil {
		return nil, err
	}

	return &ministry, nil
}
