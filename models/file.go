package models

import (
	"time"
)

type File struct {
	ID        int        `gorm:"type:bigserial; primary_key"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deleted_at"`
	AuthorID  int        `gorm:"not null" json:"author_id"`
	Title     string     `gorm:"not null" json:"title"`
	MimeType  string     `gorm:"type:varchar(100)" json:"mime_type"`
	File      string     `gorm:"not null" json:"file"`
}
