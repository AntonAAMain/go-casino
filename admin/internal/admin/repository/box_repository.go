package repository

import (
	commonModel "casino/model"

	"gorm.io/gorm"
)

type BoxRepository struct {
	db *gorm.DB
}

func NewBoxRepository(db *gorm.DB) *BoxRepository {
	return &BoxRepository{db: db}
}

func (r *BoxRepository) CreateBox(name, mode string, price float32) (*commonModel.Box, error) {

	box := &commonModel.Box{
		Name:  name,
		Mode:  mode,
		Price: price,
	}

	if err := r.db.Create(box).Error; err != nil {
		return nil, err
	}

	return box, nil
}
