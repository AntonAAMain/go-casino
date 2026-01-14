package model

import "time"

type Transaction struct {
	ID        uint      `json:"id" gorm:"primaryKey;autoIncrement"` // первичный ключ
	Amount    int       `json:"name" gorm:"not null"`               // имя пользователя
	UserId    int       `json:"user_id" gorm:"not null"`            // пароль хранится только в виде хэша
	CreatedAt time.Time `json:"created_at"`                         // дата создания
	UpdatedAt time.Time `json:"updated_at"`                         // дата обновления
}
