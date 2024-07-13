package todo_test

import (
	"context"
	"fmt"
	"log"
	"testing"
	"try-go-dockertest/pkg/todo"

	"github.com/ory/dockertest"
	"github.com/ory/dockertest/docker"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const (
	Username = "root"
	Password = "password"
)

var dbClient *mongo.Client

func TestMain(m *testing.M) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		log.Fatalf("Could not construct pool: %s", err)
	}

	err = pool.Client.Ping()
	if err != nil {
		log.Fatalf("Could not connect to Docker: %s", err)
	}

	// pull mongodb docker image for version 5.0
	resource, err := pool.RunWithOptions(&dockertest.RunOptions{
		Repository: "mongo",
		Tag:        "6.0",
		Env: []string{
			// username and password for mongodb superuser
			fmt.Sprintf("MONGO_INITDB_ROOT_USERNAME=%s", Username),
			fmt.Sprintf("MONGO_INITDB_ROOT_PASSWORD=%s", Password),
		},
	}, func(config *docker.HostConfig) {
		// set AutoRemove to true so that stopped container goes away by itself
		config.AutoRemove = true
		config.RestartPolicy = docker.RestartPolicy{
			Name: "no",
		}
	})
	if err != nil {
		log.Fatalf("Could not start resource: %s", err)
	}

	// exponential backoff-retry, because the application in the container might not be ready to accept connections yet
	err = pool.Retry(func() error {
		var err error
		dbClient, err = mongo.Connect(
			context.TODO(),
			options.Client().ApplyURI(
				fmt.Sprintf("mongodb://%s:%s@localhost:%s", Username, Password, resource.GetPort("27017/tcp")),
			),
		)
		if err != nil {
			return err
		}
		return dbClient.Ping(context.TODO(), nil)
	})

	if err != nil {
		log.Fatalf("Could not connect to docker: %s", err)
	}

	defer func() {
		// When you're done, kill and remove the container
		if err = pool.Purge(resource); err != nil {
			log.Fatalf("Could not purge resource: %s", err)
		}

		// disconnect mongodb client
		if err = dbClient.Disconnect(context.TODO()); err != nil {
			panic(err)
		}
	}()

	m.Run()
}

func TestCreateTodo(t *testing.T) {
	repo := todo.NewRepository(dbClient)
	srv := todo.NewService(repo)

	newTodo := todo.Todo{
		ID:     primitive.NewObjectID(),
		Title:  "Hello",
		IsDone: true,
	}

	err := srv.CreateTodo(context.Background(), newTodo)
	if err != nil {
		t.Errorf("got error %s", err.Error())
	}
}
