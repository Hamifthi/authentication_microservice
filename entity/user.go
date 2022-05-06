package entity

import (
	"encoding/json"
	"github.com/go-playground/validator"
	"io"
	"time"
)

type User struct {
	ID             int16     `gorm:"primaryKey;autoIncrement" json:"id"`
	Email          string    `gorm:"not null" json:"email" validate:"required"`
	Password       string    `sql:"-" json:"password" validate:"required"`
	HashedPassword string    `json:"-"`
	TokenHash      string    `json:"-"`
	CreatedAt      time.Time `gorm:"autoCreateTime:milli" json:"-"`
	UpdatedAt      time.Time `gorm:"autoCreateTime:milli" json:"-"`
}

func (u *User) Validate() error {
	validate := validator.New()
	return validate.Struct(u)
}

func (u *User) FromJson(r io.Reader) error {
	decoder := json.NewDecoder(r)
	return decoder.Decode(u)
}
