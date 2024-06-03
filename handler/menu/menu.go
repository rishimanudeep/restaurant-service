package menu

import (
	"encoding/json"
	"net/http"
	"restaurant-service/errors"
	"strconv"

	"github.com/gorilla/mux"
	"restaurant-service/models"
	srv "restaurant-service/service"
)

type menuRegister struct {
	menuService srv.Menu
}

func New(s srv.Menu) menuRegister {
	return menuRegister{menuService: s}
}

func (m *menuRegister) Create(w http.ResponseWriter, r *http.Request) {
	var payload models.Menu

	if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := m.menuService.Create(&payload)
	if err != nil {
		m.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(payload)
}

func (m *menuRegister) GetMenu(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	restaurantID, err := strconv.Atoi(vars["restaurant_id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	menu, err := m.menuService.GetMenuByRestaurantID(restaurantID)
	if err != nil {
		m.handleError(w, err)
		return
	}

	if menu == nil {
		http.Error(w, "Menu not found", http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(menu)
}

func (m *menuRegister) Update(w http.ResponseWriter, r *http.Request) {
	var menu models.Menu
	if err := json.NewDecoder(r.Body).Decode(&menu); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := m.menuService.Update(&menu)
	if err != nil {
		m.handleError(w, err)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(menu)
}

func (m *menuRegister) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	restaurantID, err := strconv.Atoi(vars["restaurant_id"])
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err = m.menuService.Delete(restaurantID)
	if err != nil {
		m.handleError(w, err)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (m *menuRegister) handleError(w http.ResponseWriter, err error) {
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
