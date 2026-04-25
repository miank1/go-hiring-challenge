package models

type Category struct {
	ID   uint   `gorm:"primaryKey"`
	Code string `gorm:"unique;not null" json:"code"`
	Name string `gorm:"not null" json:"name"`
}
