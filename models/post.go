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
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deletedAt"`
	AuthorID  int        `gorm:"not null"`
	Title     string     `gorm:"not null"`
	Content   string
	Excerpt   string
	Status    string `gorm:"not null; type:varchar(50)"`
	Slug      string `gorm:"not null; type:varchar(255); unique;"`
	Type      string `gorm:"type:varchar(50)"`
	Meta      []PostMeta
	Terms     []Term
}

type PostMeta struct {
	ID        int        `gorm:"type:bigserial; primary_key"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deletedAt"`
	PostID    int
	Key       string `gorm:"type:varchar(255)"`
	Value     string
}

type MetaQuery struct {
	Key     string
	Value   []interface{}
	Compare string
}

type DateQuery struct {
	Column    string
	Before    string
	After     string
	Compare   string
	Inclusive bool
	Year      []interface{}
	DayOfYear []interface{}
	Month     []interface{}
	Week      []interface{}
	Day       []interface{}
	DayOfWeek []interface{}
	Hour      []interface{}
	Minute    []interface{}
	Second    []interface{}
}

type DateQueries struct {
	DataPart string
	Values   []interface{}
}
