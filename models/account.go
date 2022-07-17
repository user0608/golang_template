package models

import "time"

type Account struct {
	Username        string    `gorm:"primaryKey" json:"username"`
	FirstName       string    `json:"first_name"`
	LastName        string    `json:"last_name"`
	Email           string    `json:"email"`
	Password        string    `json:"-"`
	PasswordAttempt int       `json:"-"`
	CreatedAt       time.Time `json:"created_at"`
}
