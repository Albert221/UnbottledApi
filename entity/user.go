package entity

type User struct {
	Base
	Username  string    `json:"username" gorm:"unique_index;not null"`
	Email     string    `json:"email" gorm:"unique_index;not null"`
	Password  string    `json:"password" gorm:"not null"`
}
