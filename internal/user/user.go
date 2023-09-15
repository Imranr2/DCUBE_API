package user

import "time"

type User struct {
	ID        uint      `json:"-" gorm:"primaryKey"`
	Username  string    `json:"-" gorm:"unique;not null"`
	Password  string    `json:"-" gorm:"not null"`
	CreatedAt time.Time `json:"-" gorm:"type:timestamp;default:current_timestamp"`
	UpdatedAt time.Time `json:"-" gorm:"type:timestamp;default:current_timestamp ON update current_timestamp"`
}

type Request struct {
	Username string `json:"username" validate:"required,max=32"`
	Password string `json:"password" validate:"required,min=8"`
}

type Response struct {
	User User `json:"user"`
}
