package entity

import (
	"errors"
	"regexp"
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

	// Validate Name format
	nameRegex := regexp.MustCompile(`^[A-Z0-9_]+$`)
	if !nameRegex.MatchString(ff.Name) {
		return errors.New("name must be uppercase and contain only letters, numbers, underscores")
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
