package impl

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"lab3/internal/infra/abs"
	"time"
)

type MongoDBRepository struct {
	client     *mongo.Client
	database   *mongo.Database
	collection *mongo.Collection
}

func NewMongoDBRepository(collectionName string) (abs.Repository, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI("mongodb://user:pass@localhost:27017"))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	database := client.Database("file_storage")
	collection := database.Collection(collectionName)

	return &MongoDBRepository{client: client, database: database, collection: collection}, nil
}

func (r *MongoDBRepository) Insert(data interface{}) error {
	_, err := r.collection.InsertOne(context.TODO(), data)
	return err
}

func (r *MongoDBRepository) FindOneByField(field string, value interface{}, result interface{}) error {
	filter := bson.D{{field, value}}
	return r.collection.FindOne(context.TODO(), filter).Decode(result)

}
