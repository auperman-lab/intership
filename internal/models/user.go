package models

import "go.mongodb.org/mongo-driver/bson/primitive"

type UserModel struct {
	ID       primitive.ObjectID `bson:"_id,omitempty"`
	Email    string             `json:"email" gorm:"unique"`
	Password string             `json:"password" gorm:"not null"`
}

type TokenDetails struct {
	Token     *string
	UserID    uint
	ExpiresIn *int64
}
