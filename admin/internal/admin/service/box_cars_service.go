package service

import (
	"casino/admin/internal/admin/repository"

	commonModel "casino/model"
)

type BoxCarService struct {
	repo repository.BoxCarsRepository
}

func NewBoxCarService(repo *repository.BoxCarsRepository) *BoxCarService {
	return &BoxCarService{repo: *repo}
}

func (s *BoxCarService) CreateBoxCar(carId []int, boxId int) ([]*commonModel.BoxCar, error) {

	boxCar, err := s.repo.CreateBoxCars(carId, boxId)

	if err != nil {
		return nil, err
	}

	return boxCar, nil

}
