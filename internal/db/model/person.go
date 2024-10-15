package model

type Person struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"name,size:255"`
	Email string `gorm:"email,unique,size:255"`
}

func (Person) TableName() string {
	return "person"
}

type PersonWithAssignment struct {
	ID         uint
	Name       string
	Email      string
	IsAssigned bool
	IsGlobal   bool
}

func (PersonWithAssignment) TableName() string {
	return "person"
}

type AssignedFeatureFlag struct {
	ID         uint   `json:"id"`
	Name       string `json:"name"`
	IsActive   bool   `json:"is_active"`
	IsGlobal   bool   `json:"is_global"`
	IsAssigned bool   `json:"is_assigned"`
}
