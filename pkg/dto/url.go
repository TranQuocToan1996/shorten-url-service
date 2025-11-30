package dto

import "time"

type SubmitShortenURLRequest struct {
	LongURL     string `json:"long_url" binding:"required,url"`
	CallbackURL string `json:"callback_url,omitempty"`
}

type GetDecodeURLRequest struct {
	ShortenURL string `form:"shorten_url" binding:"required"`
}

type GetDecodeURLResponse struct {
	ID        int64     `gorm:"column:id" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
	Status    string    `gorm:"column:status" json:"status"`
	Code      string    `gorm:"column:code" json:"code"`
	LongURL   string    `gorm:"column:long_url" json:"long_url"`
}
