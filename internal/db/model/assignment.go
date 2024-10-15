package model

type Assignment struct {
	ID            uint         `gorm:"primaryKey;autoIncrement" json:"id"`
	Person        *Person      `gorm:"foreignKey:PersonID"`
	PersonID      uint         `gorm:"column:person_id" json:"person_id"`
	FeatureFlag   *FeatureFlag `gorm:"foreignKey:FeatureFlagID"`
	FeatureFlagID uint         `gorm:"column:feature_flag_id" json:"feature_flag_id"`
}

func (Assignment) TableName() string {
	return "feature_flag_assignments"
}
