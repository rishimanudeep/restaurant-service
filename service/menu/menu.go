package menu

import (
	"restaurant-service/models"
	"restaurant-service/service"
)

type menuSvc struct {
	menuStore service.Menu
}

func New(menuStore service.Menu) menuSvc {
	return menuSvc{menuStore: menuStore}
}

func (s *menuSvc) Create(payload *models.Menu) error {
	err := s.menuStore.Create(payload)
	if err != nil {
		return err
	}

	return nil
}

func (s *menuSvc) GetMenuByRestaurantID(restaurantID int) (*models.Menu, error) {
	return s.menuStore.GetMenuByRestaurantID(restaurantID)
}

func (s *menuSvc) Update(menu *models.Menu) error {
	return s.menuStore.Update(menu)
}

func (s *menuSvc) Delete(restaurantID int) error {
	return s.menuStore.Delete(restaurantID)
}
