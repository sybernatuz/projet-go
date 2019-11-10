package entities

import (
	"github.com/google/uuid"
	"time"
)

type User struct {
	ID          int        `json:"-" gorm:"primary_key, AUTO_INCREMENT"`
	UUID        uuid.UUID  `json:"uuid" gorm:"type:uuid"`
	AccessLevel int        `json:"-"`
	Password    string     `json:"-"`
	FirstName   string     `json:"first_name"`
	LastName    string     `json:"last_name"`
	Email       string     `json:"email"`
	DateOfBirth time.Time  `json:"birth_date"`
	CreatedAt   time.Time  `json:"-"`
	UpdatedAt   time.Time  `json:"-"`
	DeletedAt   *time.Time `json:"-"`
}
