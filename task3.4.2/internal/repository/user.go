package repository

import (
	"context"
	"database/sql"
)

type Conditions struct {
	Limit  int `json:"limit"`
	Offset int `json:"offset"`
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
	ID      int    `json:"id"`
	Name    string `json:"name"`
	Deleted bool   `json:"deleted,omitempty"`
}

type UserRep struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRep {
	return &UserRep{db: db}
}

func (u *UserRep) Create(ctx context.Context, user User) error {
	query := `INSERT INTO users (name) VALUES ($1)`
	_, err := u.db.ExecContext(ctx, query, user.Name)
	return err
}

func (u *UserRep) GetByID(ctx context.Context, id string) (User, error) {
	user := User{}
	query := `SELECT * FROM users WHERE id = $1 AND deleted = false`
	err := u.db.QueryRowContext(ctx, string(query), id).Scan(&user.ID, &user.Name, &user.Deleted)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (u *UserRep) Update(ctx context.Context, user User) error {
	query := `UPDATE users SET name = $2 WHERE id = $1 AND deleted = false`
	_, err := u.db.ExecContext(ctx, query, user.ID, user.Name)
	if err != nil {
		return err
	}
	return nil
}

// Delete Must be rewrote to soft delete
func (u *UserRep) Delete(ctx context.Context, id string) error {
	query := `UPDATE users SET deleted = true WHERE id = $1 AND deleted = false`
	_, err := u.db.ExecContext(ctx, query, id)
	return err
}

func (u *UserRep) List(ctx context.Context, c Conditions) ([]User, error) {
	var users []User
	query := `SELECT * FROM users LIMIT $1 OFFSET $2`
	rows, err := u.db.QueryContext(ctx, query, c.Limit, c.Offset)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		var user User
		if err = rows.Scan(&user.ID, &user.Name, &user.Deleted); err != nil {
			return nil, err
		}
		if !user.Deleted {
			users = append(users, user)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
