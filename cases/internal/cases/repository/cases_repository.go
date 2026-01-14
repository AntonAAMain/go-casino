package repository

import (
	"casino/model"
	"errors"

	commonModel "casino/model"

	"gorm.io/gorm"
)

type CasesRepository struct {
	db *gorm.DB
}

func NewCasesRepository(db *gorm.DB) *CasesRepository {
	return &CasesRepository{db: db}
}

func (r *CasesRepository) GetBoxInfo(boxId int) (float32, string, error) {

	box := commonModel.Box{}

	if err := r.db.First(&box, boxId).Error; err != nil {
		return 0, "", err
	}

	return box.Price, box.Mode, nil
}

func (r *CasesRepository) GetAllCarsInCase(boxId int) ([]*model.Car, error) {
	var boxCars []model.BoxCar
	if err := r.db.Where("box_id = ?", boxId).Find(&boxCars).Error; err != nil {
		return nil, err
	}

	if len(boxCars) == 0 {
		return nil, errors.New("case has no cars")
	}

	carIDs := make([]int, 0, len(boxCars))
	for _, bc := range boxCars {
		carIDs = append(carIDs, bc.CarId)
	}

	var cars []*model.Car
	if err := r.db.Where("id IN ?", carIDs).Find(&cars).Error; err != nil {
		return nil, err
	}

	return cars, nil
}

func (r *CasesRepository) AddCarToUser(userID, carID int, price float32) error {
	return r.db.Transaction(func(tx *gorm.DB) error {

		result := tx.Model(&model.User{}).
			Where("id = ? AND balance >= ?", userID, price).
			UpdateColumn("balance", gorm.Expr("balance - ?", price))

		if result.Error != nil {
			return result.Error
		}

		if result.RowsAffected == 0 {
			return errors.New("user not found or insufficient balance")
		}

		userCar := &model.UsersCar{
			UserId: userID,
			CarId:  carID,
		}

		if err := tx.Create(userCar).Error; err != nil {
			return err
		}

		return nil
	})
}
