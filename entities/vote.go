package entities

import (
	"github.com/google/uuid"
	"github.com/lib/pq"

	"time"
)

type Vote struct {
	ID          int            `json:"-" gorm:"primary_key;auto_increment"`
	UUID        uuid.UUID      `json:"uuid" gorm:"type:uuid;unique_index"`
	Title       string         `json:"title" gorm:"type:varchar(255)"`
	Description string         `json:"desc" gorm:"type:text"`
	UuidVotes   pq.StringArray `json:"uuid_votes" gorm:"type:varchar(255)[]"`
	StartDate   *time.Time     `json:"start_date"`
	EndDate     *time.Time     `json:"end_date"`
	CreatedAt   time.Time      `json:"-"`
	UpdatedAt   time.Time      `json:"-"`
	DeletedAt   *time.Time     `json:"-"`
}

type VoteEdition struct {
	Title       string     `json:"title"`
	Description string     `json:"description"`
	StartDate   *time.Time `json:"start_date"`
	EndDate     *time.Time `json:"end_date"`
}
