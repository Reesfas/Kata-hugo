package repository

import (
	"context"
	"database/sql"
	"errors"
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
	Login(ctx context.Context, username, password string) (*User, error)
	CreateUserWithArray(users []User) error
	CreateUserWithList(users []User) error
}

type User struct {
	Id        int
	Username  string
	FirstName string
	LastName  string
	Email     string
	Password  string
	Phone     string
	Status    int
	Deleted   bool
}

type UserRep struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRep {
	return &UserRep{db: db}
}

func (u *UserRep) Create(ctx context.Context, user User) error {
	query := `INSERT INTO users (username, firstName, lastName, email, password, phone, userStatus) VALUES ($1, $2, $3, $4, $5, $6, $7)`
	_, err := u.db.ExecContext(ctx, query, user.Username, user.FirstName, user.LastName, user.Email, user.Password, user.Phone, user.Status)
	return err
}

func (u *UserRep) GetByUsername(ctx context.Context, username string) (User, error) {
	user := User{}
	query := `SELECT * FROM users WHERE username = $1 AND deleted = false`
	err := u.db.QueryRowContext(ctx, string(query), username).Scan(&user.Username, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Phone, &user.Status)
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
	_, err := u.db.ExecContext(ctx, query, user.Id, user.Username, user.FirstName, user.LastName, user.Email, user.Password, user.Phone, user.Status)
	if err != nil {
		return err
	}
	return nil
}

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
		if err = rows.Scan(&user.Id, &user.Username, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Phone, &user.Status); err != nil {
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

func (u *UserRep) Login(ctx context.Context, username, password string) (*User, error) {
	user := &User{}

	query := "SELECT * FROM users WHERE username=$1 AND password=$2 AND deleted = false"
	err := u.db.QueryRowContext(ctx, string(query), username, password).Scan(&user.Username, &user.FirstName, &user.LastName, &user.Email, &user.Password, &user.Phone, &user.Status)
	if err != nil {
		return nil, errors.New("invalid username or password")
	}

	return user, nil
}

func (u *UserRep) CreateUserWithArray(users []User) error {
	tx, err := u.db.Begin()
	if err != nil {
		return err
	}

	defer func(tx *sql.Tx) {
		err = tx.Rollback()
		if err != nil {

		}
	}(tx)

	stmt, err := tx.Prepare(`INSERT INTO users (username, firstName, lastName, email, password, phone, userStatus) VALUES ($1, $2, $3, $4, $5, $6, $7)`)
	if err != nil {
		return err
	}
	defer func(stmt *sql.Stmt) {
		err = stmt.Close()
		if err != nil {

		}
	}(stmt)

	for _, user := range users {
		_, err = stmt.Exec(user.Username, user.FirstName, user.LastName, user.Email, user.Password, user.Phone, user.Status)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}

func (u *UserRep) CreateUserWithList(users []User) error {
	tx, err := u.db.Begin()
	if err != nil {
		return err
	}

	defer func(tx *sql.Tx) {
		err = tx.Rollback()
		if err != nil {

		}
	}(tx)

	for _, user := range users {
		_, err = tx.Exec("INSERT INTO users (username, firstName, lastName, email, password, phone, userStatus) VALUES ($1, $2, $3, $4, $5, $6, $7)",
			user.Username, user.FirstName, user.LastName, user.Email, user.Password, user.Phone, user.Status)
		if err != nil {
			return err
		}
	}

	err = tx.Commit()
	if err != nil {
		return err
	}

	return nil
}
