package data

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Repository holds read database client.
type Repository struct {
	mongoClient *mongo.Client
}

// NewRepository returns instance of linkopus READ repository.
func NewRepository(mongoClient *mongo.Client) *Repository {
	return &Repository{mongoClient}
}

// FetchOneLinkMap will lookup for exactly one match for given bson.M as criteria.
func (r *Repository) FetchOneLinkMap(criteria bson.M) (bson.M, error) {
	appDatabase := r.mongoClient.Database("app")
	linksCollection := appDatabase.Collection("links")

	var linkMap bson.M
	if err := linksCollection.FindOne(context.Background(), criteria).Decode(&linkMap); err != nil {
		return nil, err
	}

	return linkMap, nil
}

// FetchLinkViews will lookup for LinkViews match for given bson.M as criteria and limit as int.
func (r *Repository) FetchLinkViews(criteria bson.M, limit int) ([]bson.M, error) {
	appDatabase := r.mongoClient.Database("app")
	linksCollection := appDatabase.Collection("views")

	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))

	cursor, err := linksCollection.Find(context.Background(), criteria)
	if err != nil {
		return nil, err
	}

	var linkMaps []bson.M
	for cursor.Next(context.TODO()) {
		var linkMap bson.M
		err := cursor.Decode(&linkMap)
		if err != nil {
			return nil, err
		}

		linkMaps = append(linkMaps, linkMap)
	}

	return linkMaps, nil
}

// FetchTopLinksViews will lookup for top LinkViews match for given bson.M as criteria and limit as int.
func (r *Repository) FetchTopLinksViews(criteria bson.M, limit int) ([]bson.M, error) {
	appDatabase := r.mongoClient.Database("app")
	linksCollection := appDatabase.Collection("views")

	findOptions := options.Find()
	findOptions.SetLimit(int64(limit))
	findOptions.SetSort(bson.D{{Key: "amount", Value: -1}})

	cursor, err := linksCollection.Find(context.Background(), criteria, findOptions)
	if err != nil {
		return nil, err
	}

	var linkMaps []bson.M
	for cursor.Next(context.TODO()) {
		var linkMap bson.M
		err := cursor.Decode(&linkMap)
		if err != nil {
			return nil, err
		}

		linkMaps = append(linkMaps, linkMap)
	}

	return linkMaps, nil
}
