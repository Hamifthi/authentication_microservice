package entity

import "time"

type User struct {
	ID             int16     `json:"id"`
	Email          string    `json:"email"`
	HashedPassword string    `json:"password"`
	TokenHash      string    `json:"tokenhash"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}