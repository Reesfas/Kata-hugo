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
	GetByUsername(ctx context.Context, username string) (User, error)
	Update(ctx context.Context, user User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, c Conditions) ([]User, error)
	// add login and logout method!
}

type User struct {
	id        int
	username  string
	firstName string
	lastName  string
	email     string
	password  string
	phone     string
	status    int
	deleted   bool
}

type UserRep struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRep {
	return &UserRep{db: db}
}

func (u *UserRep) Create(ctx context.Context, user User) error {
	query := `INSERT INTO users (username, firstName, lastName, email, password, phone, userStatus) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := u.db.ExecContext(ctx, query, user.username, user.firstName, user.lastName, user.email, user.password, user.phone, user.status)
	return err
}

func (u *UserRep) GetByUsername(ctx context.Context, username string) (User, error) {
	user := User{}
	query := `SELECT * FROM users WHERE username = $1 AND deleted = false`
	err := u.db.QueryRowContext(ctx, string(query), username).Scan(&user.username, &user.firstName, &user.lastName, &user.email, &user.password, &user.phone, &user.status)
	if err != nil {
		return User{}, err
	}
	return user, nil
}

func (u *UserRep) Update(ctx context.Context, user User) error {
	query := `
		UPDATE users 
		SET 
			username = $2,
			firstName = $3,
			lastName = $4,
			email = $5,
			password = $6,
			phone = $7,
			userStatus = $8
		WHERE 
			id = $1 AND 
			deleted = false
	`
	_, err := u.db.ExecContext(ctx, query, user.id, user.username, user.firstName, user.lastName, user.email, user.password, user.phone, user.status)
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
		if err = rows.Scan(&user.id, &user.username, &user.firstName, &user.lastName, &user.email, &user.password, &user.phone, &user.status); err != nil {
			return nil, err
		}
		if !user.deleted {
			users = append(users, user)
		}
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return users, nil
}
