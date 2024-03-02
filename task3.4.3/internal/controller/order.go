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

func (o *OrderContr) GetInventory(w http.ResponseWriter, r *http.Request) {

	inventory, err := o.serv.GetInventory(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	err = json.NewEncoder(w).Encode(inventory)
}
