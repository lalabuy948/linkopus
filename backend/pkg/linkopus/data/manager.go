package data

import (
	"context"
	"fmt"

	"go.mongodb.org/mongo-driver/mongo/options"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

// Manager holds write database client.
type Manager struct {
	mongoClient *mongo.Client
}

// NewManager returns instance of linkopus WRITE manager.
func NewManager(mongoClient *mongo.Client) *Manager {
	return &Manager{mongoClient}
}

// InsertLinkMap inserts given bson.D document into links collection in app database.
func (m *Manager) InsertLinkMap(document bson.D) error {
	appDatabase := m.mongoClient.Database("app")
	linksCollection := appDatabase.Collection("links")

	_, err := linksCollection.InsertOne(context.Background(), document)

	return err
}

// DeleteLinkMap deletes given bson.D document in links collection in app database.
func (m *Manager) DeleteLinkMap(criteria bson.M) error {
	appDatabase := m.mongoClient.Database("app")
	linksCollection := appDatabase.Collection("links")

	res, err := linksCollection.DeleteOne(context.Background(), criteria)
	fmt.Println(res, err)

	return err
}

// UpdateLinkView updates given bson.D document in views collection in app database.
func (m *Manager) UpdateLinkView(criteria bson.D, action bson.D) error {
	appDatabase := m.mongoClient.Database("app")
	viewsCollection := appDatabase.Collection("views")

	opts := options.Update().SetUpsert(true)
	_, err := viewsCollection.UpdateOne(context.Background(), criteria, action, opts)

	return err
}

// DeleteLinkMap deletes LinkView by given bson.D in links collection in app database.
func (m *Manager) DeleteLinkView(criteria bson.M) error {
	appDatabase := m.mongoClient.Database("app")
	viewsCollection := appDatabase.Collection("views")

	res, err := viewsCollection.DeleteOne(context.Background(), criteria)
	fmt.Println(res, err)

	return err
}
