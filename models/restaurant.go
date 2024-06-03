package models

import "github.com/google/uuid"

type MenuItem struct {
	ItemID      uuid.UUID `json:"item_id"`
	ItemName    string    `json:"item_name"`
	Description string    `json:"description"`
	IsAvailable bool      `json:"is_available"`
	Price       float64   `json:"price"`
}

type Menu struct {
	RestaurantID int        `json:"restaurant_id"`
	Items        []MenuItem `json:"items"`
}

type Restaurant struct {
	ID          int     `json:"id"`
	Name        string  `json:"name"`
	Description string  `json:"description"`
	IsAvailable bool    `json:"is_available"`
	Address     string  `json:"address"`
	Pincode     string  `json:"pincode"`
	Latitude    float64 `json:"latitude"`
	Longitude   float64 `json:"longitude"`
}

type OrderStatusUpdate struct {
	OrderID      int    `json:"order_id"`
	RestaurantID int    `json:"restaurant_id"`
	Status       string `json:"status"`
	ItemID       int    `json:"item_id"`
}

type OrderPlacedEvent struct {
	OrderID      int    `json:"order_id"`
	RestaurantID int    `json:"restaurant_id"`
	MenuID       int    `json:"menu_id"`
	Status       string `json:"status"`
}

type OrderStatusUpdatedEvent struct {
	OrderID int     `json:"order_id"`
	Status  string  `json:"status"`
	Lat     float64 `json:"lat"`
	Long    float64 `json:"long"`
}
