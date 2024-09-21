package entity

type FeatureFlag struct {
	ID             uint   `json:"id"`
	Name           string `json:"name"`
	Description    string `json:"description"`
	IsActive       bool   `json:"isActive"`
	ExpirationDate string `json:"expirationDate"`
}
