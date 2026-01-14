package repository

import (
	commonModel "casino/model"
	"errors"
	"fmt"

	"gorm.io/gorm"
)

type InventoryRepository struct {
	db *gorm.DB
}

func NewInventoryRepository(db *gorm.DB) *InventoryRepository {
	return &InventoryRepository{db: db}
}

func (r *InventoryRepository) SellAllUserInventory(userId int) (float32, error) {
	var userCars []commonModel.UsersCar
	var total float32

	if err := r.db.Where("user_id = ?", userId).Find(&userCars).Error; err != nil {
		return 0, err
	}

	if len(userCars) == 0 {
		return 0, nil
	}

	for _, uc := range userCars {
		var car commonModel.Car
		if err := r.db.Where("id = ?", uc.CarId).First(&car).Error; err != nil {
			return 0, err
		}
		total += car.Price
	}

	if err := r.db.Model(&commonModel.User{}).Where("id = ?", userId).
		Update("balance", gorm.Expr("balance + ?", total)).Error; err != nil {
		return 0, err
	}

	if err := r.db.Where("user_id = ?", userId).Delete(&commonModel.UsersCar{}).Error; err != nil {
		return 0, err
	}

	return total, nil
}

func (r *InventoryRepository) SellUserCar(userId, carId int) (float32, error) {

	var userCar commonModel.UsersCar
	if err := r.db.Where("user_id = ? AND id = ?", userId, carId).First(&userCar).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return 0, errors.New("user does not own this car")
		}
		return 0, err
	}

	var car commonModel.Car
	if err := r.db.Where("id = ?", userCar.CarId).First(&car).Error; err != nil {
		return 0, err
	}

	fmt.Println(car)
	if err := r.db.Model(&commonModel.User{}).
		Where("id = ?", userId).
		Update("balance", gorm.Expr("balance + ?", car.Price)).Error; err != nil {
		return 0, err
	}

	if err := r.db.Where("user_id = ? AND id = ?", userId, carId).Delete(&commonModel.UsersCar{}).Error; err != nil {
		return 0, err
	}

	return car.Price, nil
}
