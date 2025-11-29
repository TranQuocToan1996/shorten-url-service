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
	GetByCode(code string) (*model.ShortenURL, error)
	GetByLongURL(longURL string) (*model.ShortenURL, error)
}

func NewURLRepository(db *gorm.DB) URLRepository {
	return &urlRepo{
		db: db,
	}
}

func (r *urlRepo) Save(shortenURL *model.ShortenURL) error {
	return r.db.Save(shortenURL).Error
}

func (r *urlRepo) GetByCode(code string) (*model.ShortenURL, error) {
	var res model.ShortenURL
	err := r.db.Where("code = ?", code).First(&res).Error
	if err != nil {
		return nil, err
	}
	return &res, nil
}

func (r *urlRepo) GetByLongURL(longURL string) (*model.ShortenURL, error) {
	var res model.ShortenURL
	err := r.db.Where("long_url = ?", longURL).First(&res).Error
	if err != nil {
		return nil, err
	}
	return &res, nil
}
