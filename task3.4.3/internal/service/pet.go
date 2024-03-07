package service

import (
	"context"
	"github.com/google/uuid"
	"io"
	"os"
	"path/filepath"
	"task3.4.3/internal/repository"
)

type PetService interface {
	Create(ctx context.Context, pet repository.Pet) error
	GetByID(ctx context.Context, id string) (repository.Pet, error)
	GetByStatus(ctx context.Context, status string) (repository.Pet, error)
	UploadImages(file io.Reader, filename string) (string, error)
	FullUpdate(ctx context.Context, pet repository.Pet) error
	PartialUpdate(ctx context.Context, pet repository.Pet) error
	Delete(ctx context.Context, id string) error
}

type PetServ struct {
	repo repository.PetRepository
}

func NewPetServ(pet repository.PetRepository) *PetServ {
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

func (s *PetServ) UploadImages(file io.Reader, filename string) (string, error) {
	ext := filepath.Ext(filename)
	uniqueFilename := uuid.New().String() + ext

	imagePath := filepath.Join("uploads", uniqueFilename)
	f, err := os.Create(imagePath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = io.Copy(f, file)
	if err != nil {
		return "", err
	}

	return "/uploads/" + uniqueFilename, nil
}

func (s *PetServ) FullUpdate(ctx context.Context, pet repository.Pet) error {
	return s.repo.FullUpdate(ctx, pet)
}

func (s *PetServ) PartialUpdate(ctx context.Context, pet repository.Pet) error {
	return s.repo.PartialUpdate(ctx, pet)
}

func (s *PetServ) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}
