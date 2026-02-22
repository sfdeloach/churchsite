package models

import "gorm.io/gorm"

// Ministry represents a church ministry. Soft-delete model (embeds gorm.Model).
type Ministry struct {
	gorm.Model
	Name         string  `gorm:"column:name;type:varchar(255);not null" json:"name"`
	Slug         string  `gorm:"column:slug;type:varchar(255);uniqueIndex;not null" json:"slug"`
	Description  string  `gorm:"column:description;type:text" json:"description"`
	LeaderID     *uint   `gorm:"column:leader_id" json:"leader_id"`
	ContactEmail string  `gorm:"column:contact_email;type:varchar(255)" json:"contact_email"`
	MeetingTime  string  `gorm:"column:meeting_time;type:varchar(255)" json:"meeting_time"`
	Location     string  `gorm:"column:location;type:varchar(255)" json:"location"`
	IsActive     bool    `gorm:"column:is_active;default:true" json:"is_active"`
	SortOrder    int     `gorm:"column:sort_order;default:0" json:"sort_order"`
	PageContent  string  `gorm:"column:page_content;type:text" json:"page_content"`
}

func (Ministry) TableName() string { return "ministries" }
