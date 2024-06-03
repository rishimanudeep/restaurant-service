package menu

import (
	"database/sql"
	"encoding/json"

	"restaurant-service/errors"
	"restaurant-service/models"
)

type menuStore struct {
	db *sql.DB
}

func New(db *sql.DB) menuStore {
	return menuStore{db: db}
}

func (s *menuStore) Create(paylaod *models.Menu) error {
	itemsJSON, err := json.Marshal(paylaod.Items)
	if err != nil {
		return err
	}

	_, err = s.db.Exec("INSERT INTO menu_items (restaurant_id, items) VALUES ($1, $2)",
		paylaod.RestaurantID, itemsJSON)
	if err != nil {
		return &errors.InternalServerError{Message: "Query Execution Failed"}
	}
	return nil
}

func (s *menuStore) GetMenuByRestaurantID(restaurantID int) (*models.Menu, error) {
	row := s.db.QueryRow("SELECT items FROM menu_items WHERE restaurant_id = $1", restaurantID)

	var itemsJSON []byte
	if err := row.Scan(&itemsJSON); err != nil {
		if err == sql.ErrNoRows {
			return nil, &errors.EntityNotFound{Message: "No rows found in db"}
		}
		return nil, &errors.InternalServerError{Message: "Query Execution Failed"}
	}

	var items []models.MenuItem
	if err := json.Unmarshal(itemsJSON, &items); err != nil {
		return nil, err
	}

	return &models.Menu{
		RestaurantID: restaurantID,
		Items:        items,
	}, nil
}

func (s *menuStore) Update(menu *models.Menu) error {
	itemsJSON, err := json.Marshal(menu.Items)
	if err != nil {
		return &errors.InternalServerError{Message: "marshal error"}
	}

	_, err = s.db.Exec("UPDATE menu_items SET items = $1 WHERE restaurant_id = $2",
		itemsJSON, menu.RestaurantID)
	if err != nil {
		return &errors.InternalServerError{Message: "Query Execution Failed"}
	}

	return nil
}

func (s *menuStore) Delete(restaurantID int) error {
	_, err := s.db.Exec("DELETE FROM menu_items WHERE restaurant_id = $1", restaurantID)
	if err != nil {
		return &errors.InternalServerError{Message: "Query Execution Failed"}
	}
	return nil
}
