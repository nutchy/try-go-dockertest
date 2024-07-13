package main

import (
	"context"
	"time"
	"try-go-dockertest/pkg/todo"

	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func main() {

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().
		ApplyURI("mongodb://localhost:27018").
		SetAuth(options.Credential{Username: "root", Password: "example"}))
	if err != nil {
		panic(err)
	}
	todoRepo := todo.NewRepository(client)
	todoService := todo.NewService(todoRepo)

	err = todoService.CreateTodo(context.Background(), todo.Todo{
		ID:     primitive.NewObjectID(),
		Title:  "DockerTest",
		IsDone: false,
	})

	if err != nil {
		panic(err)
	}

	defer func() {
		client.Disconnect(ctx)
	}()
}
