package controller

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"task3.4.3/internal/repository"
	"task3.4.3/internal/service"
)

type OrderController interface {
	Create(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
	GetInventory(w http.ResponseWriter, r *http.Request)
}

type OrderContr struct {
	serv service.OrderService
}

func NewOrderRep(serv service.OrderService) OrderContr {
	return OrderContr{serv}
}

// Create
// @Summary Create an order
// @Description Create a new order in the store
// @Accept json
// @Produce json
// @Param order body repository.Order true "Order object that needs to be added to the store"
// @Success 200 {object} string "Successful operation"
// @Failure 400 {string} string "Invalid order data"
// @Failure 500 {string} string "Internal server error"
// @Router /store/order [post]
func (o *OrderContr) Create(w http.ResponseWriter, r *http.Request) {
	var order repository.Order
	if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	err := o.serv.Create(r.Context(), order)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
}

// GetByID
// @Summary Get order by ID
// @Description Get order information by ID
// @Produce json
// @Param orderId path string true "ID of the order to get"
// @Success 200 {object} repository.Order "Successful operation"
// @Failure 404 {string} string "Order not found"
// @Router /store/order/{orderId} [get]
func (o *OrderContr) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	user, err := o.serv.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		return
	}
}

// Delete
// @Summary Delete order by ID
// @Description Delete order by ID
// @Param orderId path string true "ID of the order to delete"
// @Success 204 "No Content"
// @Failure 500 {string} string "Internal server error"
// @Router /store/order/{orderId} [delete]
func (o *OrderContr) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := o.serv.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

// GetInventory
// @Summary Get inventory
// @Description Get inventory of the store
// @Produce json
// @Success 200 {object} object "Successful operation"
// @Failure 500 {string} string "Internal server error"
// @Router /store/inventory [get]
func (o *OrderContr) GetInventory(w http.ResponseWriter, r *http.Request) {

	inventory, err := o.serv.GetInventory(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(inventory)
}
