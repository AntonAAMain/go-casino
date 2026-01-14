package service

import (
	"casino/admin/internal/admin/repository"
	"casino/pkg/redis"
	"encoding/json"
	"fmt"

	commonModel "casino/model"
)

type CarsService struct {
	repo repository.CarsRepository
}

func NewCarsService(repo *repository.CarsRepository) *CarsService {
	return &CarsService{repo: *repo}
}

func (s *CarsService) CreateCar(name string, price float32) (*commonModel.Car, error) {

	car, err := s.repo.CreateCar(name, price)

	rdb := redis.NewRedisClient()

	rdb.Del(redis.Ctx, redis.AllCarsRedis)

	if err != nil {
		return nil, err
	}

	return car, nil

}

func (s *CarsService) GetAllCars() ([]*commonModel.Car, error) {
	rdb := redis.NewRedisClient()

	val, err := rdb.Get(redis.Ctx, redis.AllCarsRedis).Result()
	if err == nil && val != "" {
		var cars []*commonModel.Car
		if err := json.Unmarshal([]byte(val), &cars); err == nil {
			fmt.Println("we are from redis cache")
			return cars, nil
		}

		_ = rdb.Del(redis.Ctx, redis.AllCarsRedis).Err()
	}

	cars, err := s.repo.GetCars()
	if err != nil {
		return nil, err
	}

	if jsonCars, err := json.Marshal(cars); err == nil {
		_ = rdb.Set(
			redis.Ctx,
			redis.AllCarsRedis,
			jsonCars,
			redis.AllCarsExpirationTime,
		).Err()
	}

	return cars, nil
}
