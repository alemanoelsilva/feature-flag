package model

import "time"

type FeatureFlag struct {
	ExpirationDate string    `gorm:"null" json:"expiration_date"`
	CreatedAt      time.Time `gorm:"autoCreateTime" json:"created_at"`
	PersonID       uint      `gorm:"column:person_id" json:"person_id"`
}
}
