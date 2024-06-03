package store

import (
	"database/sql"
	"log"
	"restaurant-service/errors"

	"restaurant-service/models"
)

type restaurantStore struct {
	db *sql.DB
}

func New(db *sql.DB) restaurantStore {
	return restaurantStore{db: db}
}

func (s *restaurantStore) Create(paylaod *models.Restaurant) error {
	query := `INSERT INTO restaurants (id,name,description,is_available,address,pincode,latitude,longitude) VALUES ($1, $2, $3,$4,$5,$6,$7,$8)`
	args := []interface{}{
		paylaod.ID, paylaod.Name, paylaod.Description, paylaod.IsAvailable, paylaod.Address, paylaod.Pincode, paylaod.Latitude, paylaod.Longitude,
	}
	_, err := s.db.Exec(query, args...)
	if err != nil {
		return &errors.InternalServerError{Message: "Query Execution Failed"}
	}
	return nil
}

func (s *restaurantStore) Read() ([]*models.Restaurant, error) {
	rows, err := s.db.Query("SELECT id, name, description, address, pincode, latitude,longitude FROM restaurants where is_available = $1", true)
	if err != nil {
		return nil, &errors.InternalServerError{Message: "Query Execution Failed"}
	}
	defer rows.Close()

	var restaurants []*models.Restaurant
	for rows.Next() {
		var restaurant models.Restaurant
		if err := rows.Scan(&restaurant.ID, &restaurant.Name, &restaurant.Description, &restaurant.Address, &restaurant.Pincode, &restaurant.Latitude, &restaurant.Longitude); err != nil {
			return nil, err
		}
		restaurants = append(restaurants, &restaurant)
	}

	if err = rows.Err(); err != nil {
		return nil, &errors.InternalServerError{Message: "Row Error"}
	}

	return restaurants, nil
}

func (s *restaurantStore) Update(restaurant *models.Restaurant) error {
	_, err := s.db.Exec("UPDATE restaurants SET name = $1, description = $2, address = $3, pincode = $4, latitude = $5,longitude = $6,is_available = $7 WHERE id = $8",
		restaurant.Name, restaurant.Description, restaurant.Address, restaurant.Pincode, restaurant.Latitude, restaurant.Longitude, restaurant.IsAvailable, restaurant.ID)
	if err != nil {
		return &errors.InternalServerError{Message: "Query Execution Failed"}
	}
	return nil
}

func (s *restaurantStore) Delete(id int) error {
	_, err := s.db.Exec("DELETE FROM restaurants WHERE id = $1", id)
	if err != nil {
		return &errors.InternalServerError{Message: "Query Execution Failed"}
	}

	return nil
}

func (s *restaurantStore) UpdateOrderStatus(orderID int, status string) error {
	query := "UPDATE restaurant_orders SET status = $1 WHERE order_id = $2"
	_, err := s.db.Exec(query, status, orderID)
	if err != nil {
		return &errors.InternalServerError{Message: "Query Execution Failed"}
	}

	return nil
}

func (s *restaurantStore) InsertOrder(order *models.OrderPlacedEvent) error {
	query := "INSERT INTO restaurant_orders (restaurant_id, order_id, menu_id, status) VALUES ($1, $2, $3, $4)"
	_, err := s.db.Exec(query, order.RestaurantID, order.OrderID, order.MenuID, order.Status)
	if err != nil {
		return &errors.InternalServerError{Message: "Query Execution Failed"}
	}

	return nil
}

func (s *restaurantStore) GetByID(id int) (*models.Restaurant, error) {
	query := `SELECT id, name, description, is_available, address, pincode, latitude, longitude FROM restaurants WHERE id = $1`
	row := s.db.QueryRow(query, id)

	var restaurant models.Restaurant
	err := row.Scan(&restaurant.ID, &restaurant.Name, &restaurant.Description, &restaurant.IsAvailable, &restaurant.Address,
		&restaurant.Pincode, &restaurant.Latitude, &restaurant.Longitude)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		log.Printf("Failed to get restaurant by ID: %v", err)
		return nil, &errors.InternalServerError{Message: "Scan error"}
	}

	return &restaurant, nil
}
