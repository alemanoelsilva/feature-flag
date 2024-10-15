package model

type AssignmentGroup struct {
	ID       uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Name     string    `gorm:"not null;unique" json:"name"`
	Person   *[]Person `gorm:"foreignKey:PersonID"`
	PersonID []uint    `gorm:"column:person_id" json:"person_id"`
}
