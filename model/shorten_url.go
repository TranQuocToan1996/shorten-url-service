package model

import (
	"time"
)

const (
	StatusSubmit  = "submitted"
	StatusEncoded = "encoded"
	StatusFailed  = "failed"
)

const (
	AlgoBase62 = "base62"
)

type ShortenURL struct {
	ID        int64     `gorm:"column:id" json:"id"`
	CreatedAt time.Time `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time `gorm:"column:updated_at" json:"updated_at"`
	Status    string    `gorm:"column:code" json:"code"`
	Code      string    `gorm:"column:status" json:"status"`
	Algo      string    `gorm:"column:algo" json:"algo"`
	CleanURL  string    `gorm:"column:clean_url" json:"clean_url"`
}

func (ShortenURL) TableName() string {
	return "shorten_urls"
}
