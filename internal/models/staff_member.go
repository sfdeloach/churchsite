package models

import (
	"sort"

	"gorm.io/gorm"
)

// StaffCategory is a typed string for staff member categories.
type StaffCategory string

const (
	CategoryPastor StaffCategory = "pastor"
	CategoryStaff  StaffCategory = "staff"
)

// StaffCategoryInfo holds the display label and sort order for a category.
type StaffCategoryInfo struct {
	Label        string
	DisplayOrder int
}

// StaffCategories maps each category to its display metadata.
var StaffCategories = map[StaffCategory]StaffCategoryInfo{
	CategoryPastor: {Label: "Teaching Elders", DisplayOrder: 1},
	CategoryStaff:  {Label: "Staff", DisplayOrder: 2},
}

// OrderedStaffCategories returns categories sorted by DisplayOrder.
func OrderedStaffCategories() []StaffCategory {
	cats := make([]StaffCategory, 0, len(StaffCategories))
	for c := range StaffCategories {
		cats = append(cats, c)
	}
	sort.Slice(cats, func(i, j int) bool {
		return StaffCategories[cats[i]].DisplayOrder < StaffCategories[cats[j]].DisplayOrder
	})
	return cats
}

// StaffMember represents a church staff member. Soft-delete model (embeds gorm.Model).
type StaffMember struct {
	gorm.Model
	UserID       *uint         `gorm:"column:user_id" json:"user_id"`
	Name         string        `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Title        string        `gorm:"column:title;type:varchar(255);not null" json:"title"`
	Bio          string        `gorm:"column:bio;type:text" json:"bio"`
	Email        string        `gorm:"column:email;type:varchar(255)" json:"email"`
	Phone        string        `gorm:"column:phone;type:varchar(20)" json:"phone"`
	PhotoURL     string        `gorm:"column:photo_url;type:varchar(500)" json:"photo_url"`
	DisplayOrder int           `gorm:"column:display_order;default:0" json:"display_order"`
	IsActive     bool          `gorm:"column:is_active;default:true" json:"is_active"`
	Category     StaffCategory `gorm:"column:category;type:varchar(50);not null;default:'staff'" json:"category"`
}

func (StaffMember) TableName() string {
	return "staff_members"
}
