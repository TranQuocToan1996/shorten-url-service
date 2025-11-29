package dto

import "time"

type SubmitShortenURLRequest struct {
	LongURL string `json:"number" binding:"required,url"`
}

type GetDecodeURLRequest struct {
	ShortenURL string `json:"shorten_url" binding:"required,url"`
}

type GetDecodeURLResponse struct {
	ID        int64     `gorm:"column:id" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
	Status    string    `gorm:"column:code" json:"code"`
	Code      string    `gorm:"column:status" json:"status"`
	CleanURL  string    `gorm:"column:clean_url" json:"clean_url"`
}
