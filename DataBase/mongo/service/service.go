package service

import (
	"SDT_ApiServices/DataBase/mongo/connection"
	"SDT_ApiServices/DataBase/mongo/models"
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoRepository struct{}

func NewMongoRepository() mongoRepository {
	return &mongoRepository{}
}

func (r *mongoRepository) Execute(req models.MongoDynamicRequest) (interface{}, error) {

	client, err := connection.NewMongoClient(req.DBConnection.URI)
	if err != nil {
		return nil, err
	}

	db := client.Database(req.DBConnection.Database)
	collection := db.Collection(req.Collection)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M(req.Filter)

	switch req.Action {

	case "find":
		opts := options.Find()

		if req.Limit > 0 {
			opts.SetLimit(int64(req.Limit))
		}

		if req.SortBy != "" {
			order := 1
			if req.Order == "desc" {
				order = -1
			}
			opts.SetSort(bson.D{{Key: req.SortBy, Value: order}})
		}

		cursor, err := collection.Find(ctx, filter, opts)
		if err != nil {
			return nil, err
		}

		var results []bson.M
		if err = cursor.All(ctx, &results); err != nil {
			return nil, err
		}

		return results, nil

	case "insert":
		if req.IsMultiple {
			return collection.InsertMany(ctx, req.Data["documents"].([]interface{}))
		}
		return collection.InsertOne(ctx, req.Data)

	case "update":
		update := bson.M{"$set": req.Data}
		if req.IsMultiple {
			return collection.UpdateMany(ctx, filter, update)
		}
		return collection.UpdateOne(ctx, filter, update)

	case "delete":
		if req.IsMultiple {
			return collection.DeleteMany(ctx, filter)
		}
		return collection.DeleteOne(ctx, filter)
	}

	return nil, mongo.ErrClientDisconnected
}
