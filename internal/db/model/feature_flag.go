package model

import "time"

type FeatureFlag struct {
	ID             uint       `gorm:"primaryKey;autoIncrement" json:"id"`
	Name           string     `gorm:"unique;not null" json:"name"`
	Description    string     `gorm:"not null" json:"description"`
	IsActive       bool       `gorm:"not null;default:false" json:"is_active"`
	ExpirationDate *time.Time `gorm:"type:date;null" json:"expiration_date"`
	CreatedAt      time.Time  `gorm:"autoCreateTime" json:"created_at"`
	UpdatedAt      time.Time  `gorm:"autoUpdateTime" json:"updated_at"`
	Person         Person     `gorm:"foreignkey:PersonId"`
	PersonId       int        `gorm:"person_id"`
}
