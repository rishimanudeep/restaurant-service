package http

import (
	"restaurant-service/models"
)

type Restaurant interface {
	Create(payload *models.Restaurant) error
	Read() ([]*models.Restaurant, error)
	Update(restaurant *models.Restaurant) error
	OrderStatusUpdate(orderStatus *models.OrderStatusUpdate) error
	Delete(id int) error
	ProcessOrderPlacedEvent(message []byte) error
}

type Menu interface {
	Create(payload *models.MenuItem) error
	GetMenuByRestaurantID(restaurantID int) (*models.Menu, error)
	Update(menu *models.Menu) error
	Delete(restaurantID int) error
}
