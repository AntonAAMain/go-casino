package repository

import (
	commonModel "casino/model"

	"gorm.io/gorm"
)

type BoxCarsRepository struct {
	db *gorm.DB
}

func NewBoxCarsRepository(db *gorm.DB) *BoxCarsRepository {
	return &BoxCarsRepository{db: db}
}

func (r *BoxCarsRepository) CreateBoxCars(carsId []int, boxId int) ([]*commonModel.BoxCar, error) {

	var boxCars []*commonModel.BoxCar

	for _, carId := range carsId {
		boxCar := &commonModel.BoxCar{
			CarId: carId,
			BoxId: boxId,
		}

		if err := r.db.Create(boxCar).Error; err != nil {
			return nil, err
		}

		boxCars = append(boxCars, boxCar)
	}

	return boxCars, nil
}
