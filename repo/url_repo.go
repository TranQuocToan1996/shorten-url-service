package repo

import (
	"shorten/model"

	"gorm.io/gorm"
)

type urlRepo struct {
	db *gorm.DB
}

type URLRepository interface {
	Save(shortenURL *model.ShortenURL) error
}

func NewURLRepository(db *gorm.DB) URLRepository {
	return &urlRepo{
		db: db,
	}
}

func (r *urlRepo) Save(shortenURL *model.ShortenURL) error {
	return r.db.Save(shortenURL).Error
}
