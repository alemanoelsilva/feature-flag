package model

import "time"

type FeatureFlag struct {
	ID             uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name           string    `gorm:"unique;not null" json:"name"`
	Description    string    `gorm:"not null" json:"description"`
	IsActive       bool      `gorm:"not null;default:false" json:"is_active"`
	ExpirationDate string    `gorm:"null" json:"expiration_date"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time `gorm:"autoUpdateTime" json:"updated_at"`
	Person         *Person   `gorm:"foreignKey:PersonID"`
	PersonID       uint      `gorm:"column:person_id" json:"person_id"`
}

type FeatureFlagFilters struct {
	ID       uint
	Name     string
	IsActive *bool
	PersonID uint
}
