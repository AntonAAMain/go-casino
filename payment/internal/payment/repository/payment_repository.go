package repository

import (
	commonModel "casino/model"
	"casino/payment/internal/payment/model"
	"errors"

	"gorm.io/gorm"
)

type TransactionRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) *TransactionRepository {
	return &TransactionRepository{db: db}
}

func (r *TransactionRepository) CreateTransaction(amount, userId int) (*model.Transaction, error) {

	transaction := &model.Transaction{
		Amount: amount,
		UserId: userId,
	}

	err := r.db.Transaction(func(tx *gorm.DB) error {

		transaction = &model.Transaction{
			UserId: userId,
			Amount: amount,
		}
		if err := tx.Create(transaction).Error; err != nil {
			return err
		}

		result := tx.Model(&commonModel.User{}).Where("id = ?", userId).
			UpdateColumn("balance", gorm.Expr("balance + ?", amount))
		if result.Error != nil {
			return result.Error
		}
		if result.RowsAffected == 0 {
			return errors.New("user not found")
		}

		return nil
	})

	if err != nil {
		return nil, err
	}

	return transaction, nil
}
