package models

import (
	"time"
)

type (
	User struct {
		ID         uint       `gorm:"primary_key"`
		CreatedAt  time.Time  `json:"created_at"`
		UpdatedAt  time.Time  `json:"updated_at"`
		DeletedAt  *time.Time `sql:"index" json:"deleted_at"`
		UserName   string     `json:"userName"`
		Password   string     `json:"-"`
		LastName   string     `json:"lastName"`
		FirstName  string     `json:"firstName"`
		MiddleName string     `json:"middleName"`
		Avatar     string     `json:"avatar"`
		Phone      int64      `json:"phone"`
		Email      string     `json:"email"`
	}
)
