package models

import "gorm.io/gorm"

// StaffMember represents a church staff member. Soft-delete model (embeds gorm.Model).
type StaffMember struct {
	gorm.Model
	UserID       *uint  `gorm:"column:user_id" json:"user_id"`
	Name         string `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Title        string `gorm:"column:title;type:varchar(255);not null" json:"title"`
	Bio          string `gorm:"column:bio;type:text" json:"bio"`
	Email        string `gorm:"column:email;type:varchar(255)" json:"email"`
	Phone        string `gorm:"column:phone;type:varchar(20)" json:"phone"`
	PhotoURL     string `gorm:"column:photo_url;type:varchar(500)" json:"photo_url"`
	DisplayOrder int    `gorm:"column:display_order;default:0" json:"display_order"`
	IsActive     bool   `gorm:"column:is_active;default:true" json:"is_active"`
}

func (StaffMember) TableName() string {
	return "staff_members"
}
