package entity

import (
	"errors"
	"time"
)

type FeatureFlag struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	IsActive       bool   `json:"isActive"`
	ExpirationDate string `json:"expirationDate"`
}

func (ff *FeatureFlag) Validate() error {
	if ff.Name == "" {
		return errors.New("name is required")
	}
	if ff.Description == "" {
		return errors.New("description is required")
	}

	if ff.ExpirationDate != "" {
		if _, err := time.Parse(time.DateOnly, ff.ExpirationDate); err != nil {
			return errors.New("expirationDate must be in YYYY-MM-DD format")
		}
	}
	return nil
}
