package service

import (
	"context"
	"task3.4.3/internal/repository"
)

type UserService interface {
	Create(ctx context.Context, user repository.User) error
	GetByUsername(ctx context.Context, username string) (repository.User, error)
	Update(ctx context.Context, user repository.User) error
	Delete(ctx context.Context, id string) error
	List(ctx context.Context, c repository.Conditions) ([]repository.User, error)
	// add login and logout method!
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
