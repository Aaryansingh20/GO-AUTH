package models

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
    ID            uint       `gorm:"primaryKey" json:"id"`
    First_name    *string    `json:"first_name" validate:"required,min=2,max=100" gorm:"size:100;not null"`
    Last_name     *string    `json:"last_name" validate:"required,min=2,max=100" gorm:"size:100;not null"`
    Password      *string    `json:"password" validate:"required,min=6" gorm:"size:255;not null"`
    Email         *string    `json:"email" validate:"email,required" gorm:"size:100;uniqueIndex;not null"` //validate email means it should have an @
    Phone         *string    `json:"phone" validate:"required" gorm:"size:20;not null"`
    Token         *string    `json:"token" gorm:"size:500"`
    User_type     *string    `json:"user_type" validate:"required,eq=ADMIN|eq=USER" gorm:"size:20;not null"`
    Refresh_token *string    `json:"refresh_token" gorm:"size:500"`
    Created_at    time.Time  `json:"created_at"`
    Updated_at    time.Time  `json:"updated_at"`
    DeletedAt     gorm.DeletedAt `gorm:"index" json:"-"`
    User_id       string     `json:"user_id" gorm:"size:100;uniqueIndex;not null"`
}

func (User) TableName() string {
    return "users"
}
