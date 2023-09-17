package urlshortener

import (
	"time"

	"github.com/Imranr2/DCUBE_API/internal/user"
)

type ShortenedURL struct {
	ID        uint      `json:"id" gorm:"primaryKey"`
	Original  string    `json:"original" gorm:"not null"`
	Shortened string    `json:"shortened" gorm:"index;unique;not null"`
	UserID    uint      `json:"-" gorm:"not null"`
	User      user.User `json:"-" gorm:"foreignKey:UserID;not null"`
	CreatedAt time.Time `json:"createdAt" gorm:"type:timestamp;default:current_timestamp"`
}

type GetRequest struct {
	UserID uint
}

type CreateRequest struct {
	UserID      uint
	OriginalURL string `json:"original_url" validate:"required"`
}

type DeleteRequest struct {
	UserID uint
	ID     uint
}

type GetResponse struct {
	ShortenedURLs []ShortenedURL `json:"shortened_urls"`
}

type CreateResponse struct {
	ShortenedURL ShortenedURL `json:"shortened_url"`
}

type DeleteResponse struct {
	ShortenedURL ShortenedURL `json:"shortened_url"`
}
