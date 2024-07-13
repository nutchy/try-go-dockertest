package todo

import (
	"context"
)

type Service interface {
	CreateTodo(ctx context.Context, todo Todo) error
}

type serviceImpl struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &serviceImpl{repo}
}

func (s *serviceImpl) CreateTodo(ctx context.Context, todo Todo) error {
	return s.repo.CreateTodo(ctx, todo)
}
