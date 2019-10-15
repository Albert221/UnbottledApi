package entity

type User struct {
	Base
	Username  string    `json:"username" gorm:"type:varchar(191);unique_index;not null"`
	Email     string    `json:"email" gorm:"type:varchar(191);unique_index;not null"`
	Password  string    `json:"password" gorm:"not null"`
}
