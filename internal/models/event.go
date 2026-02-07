package models

import (
	"time"

	"gorm.io/gorm"
)

// Event represents a church event. Soft-delete model (embeds gorm.Model).
type Event struct {
	gorm.Model
	Title                string     `gorm:"column:title;type:varchar(255);not null" json:"title"`
	Description          string     `gorm:"column:description;type:text" json:"description"`
	EventDate            time.Time  `gorm:"column:event_date;not null" json:"event_date"`
	EndDate              *time.Time `gorm:"column:end_date" json:"end_date"`
	Location             string     `gorm:"column:location;type:varchar(255)" json:"location"`
	LocationDetails      string     `gorm:"column:location_details;type:text" json:"location_details"`
	IsRecurring          bool       `gorm:"column:is_recurring;default:false" json:"is_recurring"`
	RecurrenceRule       string     `gorm:"column:recurrence_rule;type:varchar(20);default:'none'" json:"recurrence_rule"`
	RecurrenceEnd        *time.Time `gorm:"column:recurrence_end;type:date" json:"recurrence_end"`
	RegistrationEnabled  bool       `gorm:"column:registration_enabled;default:false" json:"registration_enabled"`
	CapacityLimit        *int       `gorm:"column:capacity_limit" json:"capacity_limit"`
	RegistrationDeadline *time.Time `gorm:"column:registration_deadline" json:"registration_deadline"`
	VisibleFrom          *time.Time `gorm:"column:visible_from" json:"visible_from"`
	VisibleUntil         *time.Time `gorm:"column:visible_until" json:"visible_until"`
	IsPublic             bool       `gorm:"column:is_public;default:true" json:"is_public"`
	MinistryID           *uint      `gorm:"column:ministry_id" json:"ministry_id"`
	CreatedBy            *uint      `gorm:"column:created_by" json:"created_by"`
}

func (Event) TableName() string {
	return "events"
}
