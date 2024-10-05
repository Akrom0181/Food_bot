package models

type Order struct {
	ID         int     `json:"id"`
	UserID     int     `json:"user_id"`
	TotalPrice float64 `json:"total_price"`
	CreatedAt  string  `json:"created_at"`
	Status     string  `json:"status"`
}
