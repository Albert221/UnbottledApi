package domain

type User struct {
	Base
	Username  string    `json:"username" gorm:"unique_index"`
	Email     string    `json:"email" gorm:"unique_index"`
	Password  string    `json:"password"`
}
