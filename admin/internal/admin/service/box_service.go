package service

import (
	"casino/admin/internal/admin/repository"

	commonModel "casino/model"
)

type BoxService struct {
	repo repository.BoxRepository
}

func NewBoxService(repo *repository.BoxRepository) *BoxService {
	return &BoxService{repo: *repo}
}

func (s *BoxService) CreateBox(name, mode string, price float32) (*commonModel.Box, error) {

	box, err := s.repo.CreateBox(name, mode, price)

	if err != nil {
		return nil, err
	}

	return box, nil

}
