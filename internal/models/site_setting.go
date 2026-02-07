package models

import "time"

// SiteSetting stores key-value configuration. Hard-delete model (no soft delete).
type SiteSetting struct {
	Key         string    `gorm:"column:key;primaryKey;type:varchar(100)" json:"key"`
	Value       string    `gorm:"column:value;type:text" json:"value"`
	Description string    `gorm:"column:description;type:text" json:"description"`
	UpdatedBy   *uint     `gorm:"column:updated_by" json:"updated_by"`
	UpdatedAt   time.Time `gorm:"column:updated_at;autoUpdateTime" json:"updated_at"`
}

func (SiteSetting) TableName() string {
	return "site_settings"
}
