package http

import (
	"encoding/json"
	"net/http"
	"restaurant-service/errors"
	"strconv"

	"github.com/gorilla/mux"
	handler "restaurant-service/handler"
	"restaurant-service/models"
)

type registerRestaurant struct {
	restaurantService handler.Restaurant
}

func New(s handler.Restaurant) registerRestaurant {
	return registerRestaurant{restaurantService: s}
}

func (reg *registerRestaurant) Create(w http.ResponseWriter, r *http.Request) {
	var payload models.Restaurant

	err := json.NewDecoder(r.Body).Decode(&payload)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}

	err = reg.restaurantService.Create(&payload)
	if err != nil {
		reg.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
}

func (reg *registerRestaurant) Read(w http.ResponseWriter, r *http.Request) {
	restaurants, err := reg.restaurantService.Read()
	if err != nil {
		reg.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(restaurants)
}

func (reg *registerRestaurant) Update(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	var restaurant models.Restaurant
	if err := json.NewDecoder(r.Body).Decode(&restaurant); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	restaurant.ID = id

	err := reg.restaurantService.Update(&restaurant)
	if err != nil {
		reg.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(restaurant)
}

func (reg *registerRestaurant) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])

	err := reg.restaurantService.Delete(id)
	if err != nil {
		reg.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (reg *registerRestaurant) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	restaurant_id, _ := strconv.Atoi(vars["restaurant_id"])
	order_id, _ := strconv.Atoi(vars["order_id"])

	var orderStatusUpdate models.OrderStatusUpdate
	if err := json.NewDecoder(r.Body).Decode(&orderStatusUpdate); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	orderStatusUpdate.RestaurantID = restaurant_id
	orderStatusUpdate.OrderID = order_id

	err := reg.restaurantService.OrderStatusUpdate(&orderStatusUpdate)
	if err != nil {
		reg.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(orderStatusUpdate)
}

func (h *registerRestaurant) handleError(w http.ResponseWriter, err error) {
	switch e := err.(type) {
	case *errors.EntityNotFound:
		http.Error(w, e.Error(), http.StatusNotFound)
	case *errors.NoResponse:
		http.Error(w, e.Error(), http.StatusNotFound)
	case *errors.MissingParam:
		http.Error(w, e.Error(), http.StatusBadRequest)
	case *errors.ValidationError:
		http.Error(w, e.Error(), http.StatusBadRequest)
	case *errors.InternalServerError:
		http.Error(w, e.Error(), http.StatusInternalServerError)
	default:
		http.Error(w, "unknown error", http.StatusInternalServerError)
	}
}
