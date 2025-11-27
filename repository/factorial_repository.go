package repository

import (
	"shorten/model"

	"gorm.io/gorm"
)

// factorialRepository implements FactorialRepository interface
type factorialRepository struct {
	db *gorm.DB
}

// FactorialRepository defines the interface for factorial data operations
type FactorialRepository interface {
	Create(calc *model.ShortenURL) error
}

// NewFactorialRepository creates a new factorial repository
func NewFactorialRepository(db *gorm.DB) FactorialRepository {
	return &factorialRepository{
		db: db,
	}
}

// Create inserts a new factorial calculation record
func (r *factorialRepository) Create(calc *model.ShortenURL) error {
	return r.db.Create(calc).Error
}
