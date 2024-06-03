package service

import (
	"restaurant-service/models"
)

type Restaurant interface {
	Create(payload *models.Restaurant) error
	Read() ([]*models.Restaurant, error)
	Update(restaurant *models.Restaurant) error
	Delete(id int) error
	UpdateOrderStatus(orderID int, status string) error
	InsertOrder(order *models.OrderPlacedEvent) error
	GetByID(id int) (*models.Restaurant, error)
}

type Menu interface {
	Create(payload *models.Menu) error
	GetMenuByRestaurantID(restaurantID int) (*models.Menu, error)
	Update(menu *models.Menu) error
	Delete(restaurantID int) error
}
