package dto

type CarResponse struct {
	ID    uint    `json:"id"`
	Name  string  `json:"name"`
	Price float32 `json:"price"`
}
