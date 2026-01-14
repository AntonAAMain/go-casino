package repository

import (
	commonModel "casino/model"

	"gorm.io/gorm"
)

type CarsRepository struct {
	db *gorm.DB
}

func NewCarsRepository(db *gorm.DB) *CarsRepository {
	return &CarsRepository{db: db}
}

func (r *CarsRepository) CreateCar(name string, price float32) (*commonModel.Car, error) {

	car := &commonModel.Car{
		Name:  name,
		Price: price,
	}

	if err := r.db.Create(car).Error; err != nil {
		return nil, err
	}

	return car, nil
}

func (r *CarsRepository) GetCars() ([]*commonModel.Car, error) {

	cars := []*commonModel.Car{}

	if err := r.db.Find(&cars).Error; err != nil {
		return nil, err
	}

	return cars, nil
}
