package models

import (
	"time"
)

type (
	User struct {
		ID         uint       `gorm:"primary_key"`
		CreatedAt  time.Time  `json:"createdAt"`
		UpdatedAt  time.Time  `json:"updatedAt"`
		DeletedAt  *time.Time `sql:"index" json:"deletedAt"`
		UserName   string     `json:"userName" gorm:"unique;not null"`
		Password   string     `json:"-"`
		LastName   string     `json:"lastName"`
		FirstName  string     `json:"firstName"`
		MiddleName string     `json:"middleName"`
		Avatar     string     `json:"avatar"`
		Phone      int64      `json:"phone"`
		Email      string     `json:"email" gorm:"type:varchar(100);unique_index"`
	}
)
