package entities

import "time"

type Model struct {
	ID        string `gorm:"primary_key"`
	CreatedAt time.Time
	UpdatedAt time.Time
	DeletedAt *time.Time
}

const (
	STATUS_SUCCESS = "success"
	STATUS_FAILURE = "failure"
)
