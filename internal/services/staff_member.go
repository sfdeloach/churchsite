package services

import (
	"github.com/sfdeloach/churchsite/internal/models"
	"gorm.io/gorm"
)

// StaffMemberService handles staff member queries.
type StaffMemberService struct {
	db *gorm.DB
}

// NewStaffMemberService creates a new StaffMemberService.
func NewStaffMemberService(db *gorm.DB) *StaffMemberService {
	return &StaffMemberService{db: db}
}

// GetActive returns active, non-deleted staff members ordered by display_order then name.
func (s *StaffMemberService) GetActive() ([]models.StaffMember, error) {
	var members []models.StaffMember

	err := s.db.
		Where("is_active = ?", true).
		Order("display_order ASC, name ASC").
		Find(&members).Error

	return members, err
}

// GroupByCategory groups staff members by their category.
func GroupByCategory(members []models.StaffMember) map[models.StaffCategory][]models.StaffMember {
	grouped := make(map[models.StaffCategory][]models.StaffMember)
	for _, m := range members {
		grouped[m.Category] = append(grouped[m.Category], m)
	}
	return grouped
}
