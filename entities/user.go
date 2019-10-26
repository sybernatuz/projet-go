package entities

import (
	"github.com/google/uuid"
)

type User struct {
	Uuid        uuid.UUID `gorm:"type:uuid;primary_key;"`
	AccessLevel int
	Username    string
	Password    string
	FirstName   string
	LastName    string
	Email       string
	BirthDate   string
}
