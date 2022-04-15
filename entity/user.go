package entity

import "time"

type User struct {
	ID             int16     `json:"id" sql:"id"`
	Email          string    `json:"email" validate:"required" sql:"email"`
	HashedPassword string    `json:"password" validate:"required" sql:"password"`
	Username       string    `json:"username" sql:"username"`
	TokenHash      string    `json:"tokenhash" sql:"tokenhash"`
	CreatedAt      time.Time `json:"createdat" sql:"createdat"`
	UpdatedAt      time.Time `json:"updatedat" sql:"updatedat"`
}
