package models

import (
	"gorm.io/gorm"
)

type UserModel struct {
	gorm.Model
	Email    string `json:"email" gorm:"unique"`
	Password string `json:"password" gorm:"not null"`
}

type TokenDetails struct {
	Token     *string
	TokenUuid string
	UserID    uint
	ExpiresIn *int64
}
