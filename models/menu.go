package models

import (
	"time"
)

type Menu struct {
	ID        int        `gorm:"type:bigserial; primary_key"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deletedAt"`
	AuthorID  int        `gorm:"not null"`
	Name      string     `gorm:"not null"`
	Items     []MenuItem
}

type MenuItem struct {
	ID        int        `gorm:"type:bigserial; primary_key"`
	CreatedAt time.Time  `json:"created_at"`
	UpdatedAt time.Time  `json:"updated_at"`
	DeletedAt *time.Time `sql:"index" json:"deletedAt"`
	MenuID    int        `json:"menu_id"`
	AuthorID  int        `gorm:"not null" json:"author_id"`
	Title     string     `gorm:"type:varchar(255)"`
	Type      string     `gorm:"type:varchar(50)"`
	Object    string     `gorm:"type:varchar(50)"`
	ObjectID  int        `json:"object_id"`
	Url       string
	Parent    int
	Order     int
	Target    string `gorm:"type:varchar(50)"`
	Classes   string `gorm:"type:varchar(255)"`
}
