package entity

type User struct {
	Base
	Username string `json:"username" gorm:"type:varchar(191);unique_index;not null"`
	Email    string `json:"email" gorm:"type:varchar(191);unique_index;not null"`
	Password string `json:"-" gorm:"not null"`
	Active   bool   `json:"active" gorm:"not null"`
}
