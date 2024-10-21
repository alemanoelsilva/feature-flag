package entity

import (
	"errors"
	personEntity "ff/internal/person/entity"
	"regexp"
	"time"
)

type FeatureFlag struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	IsActive       bool   `json:"isActive"`
	IsGlobal       bool   `json:"isGlobal"`
	ExpirationDate string `json:"expirationDate"`
}

func (ff *FeatureFlag) Validate() error {
	if ff.Name == "" {
		return errors.New("Name|Name is required")
	}

	// Validate Name format
	nameRegex := regexp.MustCompile(`^[A-Z0-9_]+$`)
	if !nameRegex.MatchString(ff.Name) {
		return errors.New("Name|Name must be uppercase and contain only letters, numbers, underscores")
	}

	if ff.Description == "" {
		return errors.New("Description|Description is required")
	}

	if ff.ExpirationDate != "" {
		if _, err := time.Parse(time.DateOnly, ff.ExpirationDate); err != nil {
			return errors.New("ExpirationDate|Expiration date must be in YYYY-MM-DD format")
		}
	}

	return nil
}

type UpdateFeatureFlag struct {
	Description    string `json:"description"`
	IsActive       bool   `json:"isActive"`
	IsGlobal       bool   `json:"isGlobal"`
	ExpirationDate string `json:"expirationDate"`
}

func (ff *UpdateFeatureFlag) Validate() error {
	if ff.Description == "" {
		return errors.New("Description|Description is required")
	}

	if ff.ExpirationDate != "" {
		if _, err := time.Parse(time.DateOnly, ff.ExpirationDate); err != nil {
			return errors.New("ExpirationDate|Expiration date must be in YYYY-MM-DD format")
		}
	}

	return nil
}

type FeatureFlagResponse struct {
	ID             string                      `json:"id"`
	Name           string                      `json:"name"`
	Description    string                      `json:"description"`
	IsActive       bool                        `json:"isActive"`
	IsGlobal       bool                        `json:"isGlobal"`
	ExpirationDate string                      `json:"expirationDate"`
	CreatedAt      string                      `json:"createdAt"`
	UpdatedAt      string                      `json:"updatedAt"`
	Person         personEntity.PersonResponse `json:"person"`
}

// type AssignedFeatureFlagResponse struct {
// 	ID         uint   `json:"id"`
// 	Name       string `json:"name"`
// 	IsAssigned bool   `json:"isAssigned"`
// }

type FeatureFlagFilters struct {
	ID       uint   `json:"id"`
	Name     string `json:"name"`
	PersonID uint   `json:"personId"`
	IsActive *bool  `json:"isActive"`
	IsGlobal *bool  `json:"isGlobal"`
}
