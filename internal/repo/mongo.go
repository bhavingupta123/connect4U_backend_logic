package repo

import (
	"context"
	"fmt"
	"ludo_backend_refactored/internal/model"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mongoRepository struct {
	collection *mongo.Collection
}

func NewMongoRepository(uri string) (Repository, error) {
	clientOpts := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		return nil, err
	}

	if err := client.Ping(context.TODO(), nil); err != nil {
		return nil, fmt.Errorf("MongoDB ping failed: %w", err)
	}
	coll := client.Database("connect4").Collection("match_stats")
	return &mongoRepository{collection: coll}, nil
}

func (r *mongoRepository) SaveResult(result model.MatchResult) error {
	_, err := r.collection.InsertOne(context.TODO(), result)
	return err
}

func (r *mongoRepository) GetStatsForPlayer(playerName string) ([]model.MatchResult, error) {
	cursor, err := r.collection.Find(context.TODO(), map[string]interface{}{"player": playerName})
	if err != nil {
		return nil, err
	}
	defer cursor.Close(context.TODO())

	var results []model.MatchResult
	for cursor.Next(context.TODO()) {
		var res model.MatchResult
		if err := cursor.Decode(&res); err == nil {
			results = append(results, res)
		}
	}
	return results, nil
}
