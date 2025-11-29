package model

import (
	"database/sql"
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
	ID        int64        `gorm:"column:id" json:"id"`
	CreatedAt time.Time    `gorm:"column:created_at" json:"created_at"`
	UpdatedAt time.Time    `gorm:"column:updated_at" json:"updated_at"`
	DeletedAt sql.NullTime `gorm:"column:deleted_at" json:"deleted_at"`
	Status    string       `gorm:"column:status" json:"status"`
	Code      string       `gorm:"column:code" json:"code"`
	Algo      string       `gorm:"column:algo" json:"algo"`
	LongURL   string       `gorm:"column:long_url" json:"long_url"`
}

func (ShortenURL) TableName() string {
	return "shorten_urls"
}
