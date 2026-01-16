package service

import (
	"casino/cases/internal/cases/repository"
	"casino/cases/internal/dto"
	"errors"
	"math/rand"
	"sort"
	"time"
)

type CasesService struct {
	repo repository.CasesRepository
}

func NewBoxService(repo *repository.CasesRepository) *CasesService {
	return &CasesService{repo: *repo}
}

func min(a, b float64) float64 {
	if a < b {
		return a
	}
	return b
}

func max(a, b float64) float64 {
	if a > b {
		return a
	}
	return b
}

func 	GroupCarsByBoxPrice(cars []*dto.CarResponse, boxPrice float32) (cheaper, equal, moreExpensive []dto.CarResponse) {
	prices := make([]float64, len(cars))
	for i, car := range cars {
		prices[i] = float64(car.Price)
	}

	sort.Float64s(prices)
	median := prices[len(prices)/2]
	lowerBound := min(float64(boxPrice), median)
	upperBound := max(float64(boxPrice), median)

	for _, car := range cars {
		dtoCar := dto.CarResponse{ID: car.ID, Name: car.Name, Price: car.Price}
		switch {
		case float64(car.Price) < lowerBound:
			cheaper = append(cheaper, dtoCar)
		case float64(car.Price) > upperBound:
			moreExpensive = append(moreExpensive, dtoCar)
		default:
			equal = append(equal, dtoCar)
		}
	}

	return
}

func GetRandomCar(cars []dto.CarResponse) (*dto.CarResponse, error) {
	if len(cars) == 0 {
		return &dto.CarResponse{}, errors.New("no cars in the group")
	}

	r := rand.New(rand.NewSource(time.Now().UnixNano()))

	index := r.Intn(len(cars))
	return &cars[index], nil
}

func (s *CasesService) OpenCase(userID int, boxId int) (*dto.CarResponse, error) {

	cars, err := s.repo.GetAllCarsInCase(boxId)

	if err != nil {
		return nil, err
	}

	result := make([]*dto.CarResponse, 0, len(cars))

	for _, car := range cars {
		result = append(result, &dto.CarResponse{
			ID:    car.ID,
			Name:  car.Name,
			Price: car.Price,
		})
	}

	// return nil, nil

	boxPrice, boxMode, err := s.repo.GetBoxInfo(boxId)

	cheaper, equal, moreExpensive := GroupCarsByBoxPrice(result, boxPrice)

	switch boxMode {
	case "easy":
		randomCar, err := GetRandomCar(cheaper)

		if err != nil {
			return nil, err
		}

		err = s.repo.AddCarToUser(userID, int(randomCar.ID), randomCar.Price)

		if err != nil {
			return nil, err
		}

		return randomCar, nil

	case "medium":
		randomCar, err := GetRandomCar(equal)

		if err != nil {
			return nil, err
		}

		err = s.repo.AddCarToUser(userID, int(randomCar.ID), randomCar.Price)

		if err != nil {
			return nil, err
		}

		return randomCar, nil

	case "hard":
		randomCar, err := GetRandomCar(moreExpensive)

		if err != nil {
			return nil, err
		}

		err = s.repo.AddCarToUser(userID, int(randomCar.ID), randomCar.Price)

		if err != nil {
			return nil, err
		}

		return randomCar, nil

	default:
		return nil, nil
	}

}
