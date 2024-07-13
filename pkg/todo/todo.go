package todo

import "go.mongodb.org/mongo-driver/bson/primitive"

type Todo struct {
	ID     primitive.ObjectID `bson:"_id"`
	Title  string             `bson:"title"`
	IsDone bool               `bson:"is_done"`
}
