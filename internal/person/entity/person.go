package entity

type PersonResponse struct {
	ID    uint   `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type PersonWithAssignmentResponse struct {
	ID         string `json:"id"`
	Name       string `json:"name"`
	Email      string `json:"email"`
	IsAssigned bool   `json:"isAssigned"`
}

type AssignedFeatureFlagResponse struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	IsActive   bool   `json:"isActive"`
	IsAssigned bool   `json:"isAssigned"`
}

type PersonFilters struct {
	FeatureFlagID uint   `json:"id"`
	Name          string `json:"name"`
	IsAssigned    *bool  `json:"isAssigned"`
}
