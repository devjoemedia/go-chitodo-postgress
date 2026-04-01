package models

import (
	"time"

	"gorm.io/gorm"
)

type CreateTicketRequest struct {
	Title       string `json:"title" validate:"required,min=1,max=255"`
	Description string `json:"description" validate:"omitempty,max=1000"`
	Status      string `json:"status" validate:"oneof=open in_progress closed"`
	Priority    string `json:"priority" validate:"oneof=low medium high"`
	AssigneeID  *uint  `json:"assignee_id" validate:"omitempty,gt=0"`
}

type UpdateTicketRequest struct {
	Title       *string `json:"title" validate:"omitempty,min=1,max=255"`
	Description *string `json:"description" validate:"omitempty,max=1000"`
	Status      *string `json:"status" validate:"omitempty,oneof=open in_progress closed"`
	Priority    *string `json:"priority" validate:"omitempty,oneof=low medium high"`
	AssigneeID  *uint   `json:"assignee_id" validate:"omitempty,gt=0"`
}

type Ticket struct {
	ID uint `gorm:"primaryKey" json:"id"`

	Title       string `json:"title" validate:"required,min=5,max=255"`
	Description string `json:"description" validate:"omitempty,max=1000"`
	Status      string `json:"status" validate:"oneof=open in_progress closed"`
	Priority    string `json:"priority" validate:"oneof=low medium high"`

	ReporterID uint  `json:"reporter_id"`
	Reporter   *User `gorm:"foreignKey:ReporterID" json:"reporter"`

	AssigneeID *uint `json:"assignee_id"`
	Assignee   *User `gorm:"foreignKey:AssigneeID" json:"assignee"`

	CreatedAt time.Time      `json:"created_at"`
	UpdatedAt time.Time      `json:"updated_at"`
	DeletedAt gorm.DeletedAt `gorm:"index" json:"-"`
}
