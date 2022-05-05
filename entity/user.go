package entity

import (
	"encoding/json"
	"github.com/go-playground/validator"
	"io"
	"time"
)

type User struct {
	ID             int16     `json:"id"`
	Email          string    `json:"email" validate:"required"`
	Password       string    `json:"-" validate:"required"`
	HashedPassword string    `json:"password"`
	TokenHash      string    `json:"tokenhash"`
	CreatedAt      time.Time `json:"-"`
	UpdatedAt      time.Time `json:"-"`
}

func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

func (u *User) FromJson(r io.Reader) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(u)
}
