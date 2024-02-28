package repository

import (
	"context"
	"database/sql"
)

type Pet struct {
	id       int
	category Category
	name     string
	photoUrl string //change to xml later
	tags     string
	status   string
	deleted  bool
}

type Category struct {
	id   int
	name string
}

type Tag struct {
	id   int
	name string
}

type PetRepository interface {
	Create(ctx context.Context, pet Pet) error
	GetByID(ctx context.Context, id string) (Pet, error)
	GetByStatus(ctx context.Context, status string) (Pet, error)
	UploadImages()
	// Update Тут два апдейта но я не понял разницы do it later
	Update(ctx context.Context, pet Pet) error
	Delete(ctx context.Context, id string) error
}

type PetRep struct {
	db *sql.DB
}

func NewPetRep(db *sql.DB) PetRep {
	return PetRep{db: db}
}

func (p *PetRep) Create(ctx context.Context, pet Pet) error {
	query := `INSERT INTO pet (category, name, photo, tags, status) VALUES ($1, $2, $3, $4, $5)`
	_, err := p.db.ExecContext(ctx, query, pet.category, pet.name, pet.photoUrl, pet.tags, pet.status)
	return err
}

func (p *PetRep) GetByID(ctx context.Context, id string) (Pet, error) {
	pet := Pet{}
	query := `SELECT * FROM pet WHERE id = $1 AND deleted = false`
	err := p.db.QueryRowContext(ctx, string(query), id).Scan(&pet.category, &pet.name, &pet.photoUrl, &pet.tags, &pet.status, &pet.deleted)
	if err != nil {
		return Pet{}, err
	}
	return pet, nil
}

func (p *PetRep) GetByStatus(ctx context.Context, status string) (Pet, error) {
	pet := Pet{}
	query := `SELECT * FROM pet WHERE status = $1 AND deleted = false`
	err := p.db.QueryRowContext(ctx, string(query), status).Scan(&pet.category, &pet.name, &pet.photoUrl, &pet.tags, &pet.status, &pet.deleted)
	if err != nil {
		return Pet{}, err
	}
	return pet, nil
}

func (p *PetRep) UploadImages() {

}

func (p *PetRep) Delete(ctx context.Context, id string) error {
	query := `UPDATE pet SET deleted = true WHERE id = $1 AND deleted = false`
	_, err := p.db.ExecContext(ctx, query, id)
	return err
}
