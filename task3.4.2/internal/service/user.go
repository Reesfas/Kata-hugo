package service

import (
	"context"
	"hugo/task3.4.2/internal/repository"
)

// Как я понимаю, тут должны быть какие-то бизнес логики типа валидации данных, но у нас в задании ничего нет поэтому просто он вызывает методы репозитория

type UserService interface {
	Create(ctx context.Context, user repository.User) error
	GetByID(ctx context.Context, id string) (repository.User, error)
	Update(ctx context.Context, user repository.User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, c repository.Conditions) ([]repository.User, error)
}

type UserSer struct {
	repo repository.UserRepository
}

func NewUserService(user repository.UserRepository) *UserSer {
	return &UserSer{repo: user}
}

func (u *UserSer) Create(ctx context.Context, user repository.User) error {
	return u.repo.Create(ctx, user)
}

func (u *UserSer) GetByID(ctx context.Context, id string) (repository.User, error) {
	return u.repo.GetByID(ctx, id)
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
