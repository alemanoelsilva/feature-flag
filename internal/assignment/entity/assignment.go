package entity

import "errors"

type Assignment struct {
	PersonID      uint `json:"personId"`
	FeatureFlagID uint `json:"featureFlagId"`
}

func (ff *Assignment) Validate() error {
	if ff.PersonID == 0 {
		return errors.New("person id is required")
	}

	if ff.FeatureFlagID == 0 {
		return errors.New("feature flag id is required")
	}

	return nil
}
