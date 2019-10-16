package entity

import "github.com/google/uuid"

type Point struct {
	Base
	AuthorID  uuid.UUID `json:"author_id" gorm:"type:char(36);not null"`
	Author    User      `json:"-" gorm:"foreignkey:AuthorID"`
	Latitude  float32   `json:"latitude" gorm:"not null"`
	Longitude float32   `json:"longitude" gorm:"not null"`
	PhotoUrl  string    `json:"photo_url"`
}
