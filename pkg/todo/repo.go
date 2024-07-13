package todo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type Repository interface {
	CreateTodo(ctx context.Context, todo Todo) error
}

type repositoryImpl struct {
	client *mongo.Client
}

func NewRepository(client *mongo.Client) Repository {
	return &repositoryImpl{client}
}

func (r *repositoryImpl) CreateTodo(ctx context.Context, todo Todo) error {
	_, err := r.client.Database("todoDB").Collection("todos").InsertOne(ctx, todo)
	return err
}
