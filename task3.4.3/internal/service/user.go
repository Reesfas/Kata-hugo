package service

import (
	"context"
	"github.com/dgrijalva/jwt-go"
	"golang.org/x/crypto/bcrypt"
	"task3.4.3/internal/repository"
	"time"
)

type UserService interface {
	Create(ctx context.Context, user repository.User) error
	GetByUsername(ctx context.Context, username string) (repository.User, error)
	Update(ctx context.Context, user repository.User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, c repository.Conditions) ([]repository.User, error)
	Login(ctx context.Context, username, password string) (string, error)
	CreateUserWithArray(users []repository.User) error
	CreateUserWithList(users []repository.User) error
}

type UserSer struct {
	repo   repository.UserRepository
	jwtKey []byte
}

func NewUserService(user repository.UserRepository) *UserSer {
	return &UserSer{repo: user}
}

func (u *UserSer) Create(ctx context.Context, user repository.User) error {
	return u.repo.Create(ctx, user)
}

func (u *UserSer) GetByUsername(ctx context.Context, id string) (repository.User, error) {
	return u.repo.GetByUsername(ctx, id)
}

func (u *UserSer) Update(ctx context.Context, user repository.User) error {
	return u.repo.Update(ctx, user)
}

func (u *UserSer) Delete(ctx context.Context, id string) error {
	return u.repo.Delete(ctx, id)
}

func (u *UserSer) List(ctx context.Context, c repository.Conditions) ([]repository.User, error) {
	return u.repo.List(ctx, c)
}

func (u *UserSer) Login(ctx context.Context, username, password string) (string, error) {
	user, err := u.repo.Login(ctx, username, password)
	if err != nil {
		return "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		return "", err
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(),
	})

	tokenString, err := token.SignedString([]byte("secret"))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

func (u *UserSer) CreateUserWithArray(users []repository.User) error {
	return u.repo.CreateUserWithArray(users)
}

func (u *UserSer) CreateUserWithList(users []repository.User) error {
	return u.repo.CreateUserWithList(users)
}
