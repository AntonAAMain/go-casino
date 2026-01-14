package service

import (
	"casino/inventory/internal/inventory/repository"
)

type InventoryService struct {
	repo repository.InventoryRepository
}

func NewInventoryService(repo *repository.InventoryRepository) *InventoryService {
	return &InventoryService{repo: *repo}
}

func (s *InventoryService) SellUserInventory(userId int) (float32, error) {
	total, err := s.repo.SellAllUserInventory(userId)

	if err != nil {
		return 0, err
	}

	return total, nil
}

func (s *InventoryService) SellUserCar(userId, carId int) (float32, error) {
	total, err := s.repo.SellUserCar(userId, carId)

	if err != nil {
		return 0, err
	}

	return total, nil
}
