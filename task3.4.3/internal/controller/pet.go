package controller

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"net/http"
	"task3.4.3/internal/repository"
	"task3.4.3/internal/service"
)

type PetController interface {
	Create(w http.ResponseWriter, r *http.Request)
	GetByID(w http.ResponseWriter, r *http.Request)
	GetByStatus(w http.ResponseWriter, r *http.Request)
	UploadImages(w http.ResponseWriter, r *http.Request)
	FullUpdate(w http.ResponseWriter, r *http.Request)
	PartialUpdate(w http.ResponseWriter, r *http.Request)
	Delete(w http.ResponseWriter, r *http.Request)
}

type PetContr struct {
	serv service.PetService
}

func NewPetRep(pet service.PetService) *PetContr {
	return &PetContr{serv: pet}
}

func (p *PetContr) Create(w http.ResponseWriter, r *http.Request) {
	var pet repository.Pet
	if err := json.NewDecoder(r.Body).Decode(&pet); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	err := p.serv.Create(r.Context(), pet)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

}

func (p *PetContr) GetByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	user, err := p.serv.GetByID(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(user)
	if err != nil {
		return
	}
}

func (p *PetContr) GetByStatus(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	status := vars["status"]

	pet, err := p.serv.GetByStatus(r.Context(), status)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	err = json.NewEncoder(w).Encode(pet)
	if err != nil {
		return
	}
}

func (p *PetContr) FullUpdate(w http.ResponseWriter, r *http.Request) {
	var pet repository.Pet
	err := json.NewDecoder(r.Body).Decode(&pet)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	err = p.serv.FullUpdate(r.Context(), pet)
	if err != nil {
		http.Error(w, "Failed to update pet", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (p *PetContr) PartialUpdate(w http.ResponseWriter, r *http.Request) {
	var pet repository.Pet
	err := json.NewDecoder(r.Body).Decode(&pet)
	if err != nil {
		http.Error(w, "Failed to decode request body", http.StatusBadRequest)
		return
	}

	err = p.serv.PartialUpdate(r.Context(), pet)
	if err != nil {
		http.Error(w, "Failed to update pet", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func (p *PetContr) Delete(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	err := p.serv.Delete(r.Context(), id)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (p *PetContr) UploadImages(w http.ResponseWriter, r *http.Request) {
	// Обработка загруженного файла
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file from request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	// Процесс сохранения изображения
	imageURL, err := p.serv.UploadImages(file, header.Filename)
	if err != nil {
		http.Error(w, "Failed to upload image", http.StatusInternalServerError)
		return
	}

	// Отправка URL загруженного изображения в ответе
	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]string{"url": imageURL})
	if err != nil {
		return
	}

}
