package models

import (
	"time"
)

type PostConfig struct {
	Slug       string `json:"slug"`
	PluralSlug string `json:"pluralSlug"`
}

type Post struct {
	ID        int        `gorm:"type:bigserial; primary_key"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `sql:"index" json:"deletedAt"`
	AuthorID  int        `gorm:"not null"`
	Title     string     `gorm:"not null"`
	Content   string
	Excerpt   string
	Status    int    `gorm:"not null"`
	Slug      string `gorm:"not null; type:varchar(255); unique;"`
	Type      string `gorm:"type:varchar(50)"`
	Meta      []PostMeta
}

type PostMeta struct {
	ID        int        `gorm:"type:bigserial; primary_key"`
	CreatedAt time.Time  `json:"createdAt"`
	UpdatedAt time.Time  `json:"updatedAt"`
	DeletedAt *time.Time `sql:"index" json:"deletedAt"`
	PostID    int
	Key       string `gorm:"type:varchar(255)"`
	Value     string
}

type MetaQuery struct {
	Key     string
	Value   []interface{}
	Compare string
	Type    string
}
