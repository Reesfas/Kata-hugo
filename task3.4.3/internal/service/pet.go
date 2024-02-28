package service

import (
	"context"
	"task3.4.3/internal/repository"
)

type PetService interface {
	Create(ctx context.Context, pet repository.Pet) error
	GetByID(ctx context.Context, id string) (repository.Pet, error)
	GetByStatus(ctx context.Context, status string) (repository.Pet, error)
	UploadImages()
	// Update Тут два апдейта но я не понял разницы do it later
	Update(ctx context.Context, pet repository.Pet) error
	Delete(ctx context.Context, id string) error
}

type PetServ struct {
	repo repository.PetRepository
}

func NewPetRep(pet repository.PetRepository) *PetServ {
	return &PetServ{repo: pet}
}
func (s *PetServ) Create(ctx context.Context, pet repository.Pet) error {
	return s.repo.Create(ctx, pet)
}

func (s *PetServ) GetByID(ctx context.Context, id string) (repository.Pet, error) {
	return s.repo.GetByID(ctx, id)
}

func (s *PetServ) GetByStatus(ctx context.Context, status string) (repository.Pet, error) {
	return s.repo.GetByStatus(ctx, status)
}

func (s *PetServ) UploadImages() {

}

func (s *PetServ) Update(ctx context.Context, pet repository.Pet) error {
	return s.repo.Update(ctx, pet)
}

func (s *PetServ) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
