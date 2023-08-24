package models

import "gorm.io/gorm"

type UserModel struct {
	gorm.Model
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password" gorm:"not null"`
}
