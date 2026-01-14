package model

import "time"

type User struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"` // первичный ключ
	Name      string    `json:"name" gorm:"size:100;not null"`      // имя пользователя
	Password  string    `json:"-" gorm:"not null"`                  // пароль хранится только в виде хэша
	Balance   int       `json:"balance" gorm:"default:0"`           // баланс пользователя
	Role      string    `json:"role" gorm:"size:50;default:'USER'"` // роль пользователя
	CreatedAt time.Time `json:"created_at"`                         // дата создания
	UpdatedAt time.Time `json:"updated_at"`                         // дата обновления
}

type UsersCar struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	UserId    int       `json:"user_id" gorm:"not null"`
	CarId     int       `json:"car_id" gorm:"not null"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Box struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"size:100;not null"`
	Price     float32   `json:"price" gorm:"not null'"`
	Mode      string    `json:"mode" gorm:"not null'"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type Car struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	Name      string    `json:"name" gorm:"size:100;not null"`
	Price     float32   `json:"price" gorm:"not null'"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type BoxCar struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"`
	BoxId     int       `json:"box_id"`
	CarId     int       `json:"car_id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}
