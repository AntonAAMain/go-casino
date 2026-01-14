package service

import (
	"casino/payment/internal/payment/model"
	"casino/payment/internal/payment/repository"
	"errors"
)

type PaymentService struct {
	repo repository.TransactionRepository
}

func NewPaymentService(repo *repository.TransactionRepository) *PaymentService {
	return &PaymentService{repo: *repo}
}

func (s *PaymentService) CreateTransaction(amount int, userId int) (*model.Transaction, error) {

	if amount == 0 {
		return nil, errors.New("min deposit")
	}

	transaction, err := s.repo.CreateTransaction(amount, userId)

	if err != nil {
		return nil, err
	}

	return transaction, nil

}
