package models

type UserModel struct {
	ID       uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Name     string `json:"firstName" gorn:"not null"`
	Password string `json:"password"  gorn:"not null"`
}
