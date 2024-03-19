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

// Create
// @Summary Create a pet
// @Description Add a new pet to the store
// @Accept json
// @Produce json
// @Param pet body Pet true "Pet object that needs to be added to the store"
// @Success 200 {object} string "Successful operation"
// @Failure 400 {string} string "Invalid pet data"
// @Failure 500 {string} string "Internal server error"
// @Router /pet [post]
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

// GetByID
// @Summary Get pet by ID
// @Description Get pet information by ID
// @Produce json
// @Param petId path string true "ID of the pet to get"
// @Success 200 {object} Pet "Successful operation"
// @Failure 404 {string} string "Pet not found"
// @Router /pet/{petId} [get]
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

// GetByStatus
// @Summary Find pets by status
// @Description Finds pets by status
// @Produce json
// @Param status query string true "Status value to search for"
// @Success 200 {object} []Pet "Successful operation"
// @Failure 404 {string} string "Pets not found"
// @Router /pet/findByStatus [get]
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

// FullUpdate
// @Summary Update a pet
// @Description Updates a pet in the store
// @Accept json
// @Produce json
// @Param pet body Pet true "Pet object that needs to be updated in the store"
// @Success 200 "OK"
// @Failure 400 {string} string "Failed to decode request body"
// @Failure 500 {string} string "Failed to update pet"
// @Router /pet [put]
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

// PartialUpdate
// @Summary Update a pet with form data
// @Description Updates a pet in the store with form data
// @Accept json
// @Produce json
// @Param petId path string true "ID of the pet to update"
// @Param name formData string false "Name of the pet"
// @Param status formData string false "Status of the pet"
// @Success 200 "OK"
// @Failure 400 {string} string "Failed to decode request body"
// @Failure 500 {string} string "Failed to update pet"
// @Router /pet/{petId} [post]
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

// UploadImages
// @Summary Upload an image for a pet
// @Description Uploads an image for a pet
// @Param petId path string true "ID of the pet to upload image for"
// @Accept multipart/form-data
// @Produce json
// @Param file formData file true "Image file to upload"
// @Success 200 {object} object "Successful operation"
// @Failure 400 {string} string "Failed to get file from request"
// @Failure 500 {string} string "Failed to upload image"
// @Router /pet/{petId}/uploadImage [post]
func (p *PetContr) UploadImages(w http.ResponseWriter, r *http.Request) {
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Failed to get file from request", http.StatusBadRequest)
		return
	}
	defer file.Close()

	imageURL, err := p.serv.UploadImages(file, header.Filename)
	if err != nil {
		http.Error(w, "Failed to upload image", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err = json.NewEncoder(w).Encode(map[string]string{"url": imageURL})
	if err != nil {
		return
	}

}
