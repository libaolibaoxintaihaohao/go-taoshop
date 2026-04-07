package models

import "time"

type User struct {
	ID           int64     `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
	CreatedAt    time.Time `json:"created_at"`
}

type Product struct {
	ID          int64     `json:"id"`
	Name        string    `json:"name"`
	Price       float64   `json:"price"`
	Stock       int       `json:"stock"`
	CoverURL    string    `json:"cover_url"`
	Description string    `json:"description"`
	CreatedAt   time.Time `json:"created_at"`
}

type Order struct {
	ID         int64     `json:"id"`
	UserID     int64     `json:"user_id"`
	ProductID  int64     `json:"product_id"`
	Product    Product   `json:"product"`
	Quantity   int       `json:"quantity"`
	TotalPrice float64   `json:"total_price"`
	Status     string    `json:"status"`
	CreatedAt  time.Time `json:"created_at"`
}
