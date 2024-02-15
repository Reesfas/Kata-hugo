package repository

import (
	"context"
	"database/sql"
)

type Conditions struct {
}

type UserRepository interface {
	Create(ctx context.Context, user User) error
	GetByID(ctx context.Context, id string) (User, error)
	Update(ctx context.Context, user User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, c Conditions) ([]User, error)
	// Другие методы, необходимые для работы с пользователями
}

type User struct {
	ID   string
	Name string
}

type UserRep struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRep {
	return &UserRep{db: db}
}

func (u *UserRep) Create(ctx context.Context, user User) error {
	query := `INSERT INTO users (id, username,) VALUES ($1, $2)`
	_, err := u.db.ExecContext(ctx, query, user.ID, user.Name)
	return err
}

func (u *UserRep) GetByID(ctx context.Context, id string) (User, error) {
	user := User{}
	query := `SELECT * FROM users WHERE id = $1`
	err := u.db.QueryRowContext(ctx, string(query), id).Scan(&user.ID, &user.Name)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (u *UserRep) Update(ctx context.Context, user User) error {
	query := `UPDATE users SET username = $2 WHERE id = $1`
	_, err := u.db.ExecContext(ctx, query, user.ID, user.Name)
	if err != nil {
		return err
	}
	return nil
}

// Delete Must be rewrote to soft delete
func (u *UserRep) Delete(ctx context.Context, id string) error {
	query := `DELETE FROM users WHERE id = $1`
	_, err := u.db.ExecContext(ctx, query, id)
	return err
}

func (u *UserRep) List(ctx context.Context, c Conditions) ([]User, error) {
	return []User{}, nil
}
