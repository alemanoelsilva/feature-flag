package model

type Person struct {
	ID    uint   `gorm:"primaryKey"`
	Name  string `gorm:"name,size:255"`
	Email string `gorm:"email,unique,size:255"`
}

func (Person) TableName() string {
	return "person"
}
