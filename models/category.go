package models

type Category struct {
	ID   uint   `gorm:"primaryKey"`
	Code string `gorm:"unique;not null"`
	Name string `gorm:"not null"`
}
