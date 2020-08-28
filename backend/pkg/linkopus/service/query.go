package service

import (
	"errors"
	"strings"
	"time"

	"github.com/lalabuy948/linkopus/backend/pkg/linkopus/database"

	"go.mongodb.org/mongo-driver/mongo"

	"go.mongodb.org/mongo-driver/bson"
)

// Query holds represents query service. Holds read repository.
type Query struct {
	repository *database.Repository
}

// NewQuery returns instance of linkopus query service.
func NewQuery(r *database.Repository) *Query {
	return &Query{r}
}

// QueryLinkMap will build query based on given link or linkHash, query data storage and return LinkMap.
func (q *Query) QueryLinkMap(link string, linkHash string) (*LinkMap, error) {
	var criteria bson.M

	if link != "" {
		criteria = bson.M{"link": link}
	} else if linkHash != "" {
		criteria = bson.M{"linkHash": linkHash}
	}

	result, err := q.repository.FetchOneLinkMap(criteria)
	if err != nil {
		return nil, err
	}

	var linkMap LinkMap

	bsonBytes, _ := bson.Marshal(result)
	err = bson.Unmarshal(bsonBytes, &linkMap)
	if err != nil {
		return nil, err
	}

	return &linkMap, err
}

// QueryLinkMap will build query based on given link or date, query data storage and return list of LinkView.
func (q *Query) QueryLinkViews(link string, date string) (*[]LinkView, error) {

	if link == "" {
		return nil, errors.New("query: link criteria is invalid")
	}

	if date == "" {
		date = time.Now().Format("2006-01-02")
	}

	dateSplit := strings.Split(date, "-")
	year := dateSplit[0]
	month := dateSplit[1]

	criteria := bson.M{
		"link": link,
		"$and": []bson.M{
			{"year": bson.M{"$eq": year}},
			{"month": bson.M{"$eq": month}},
		},
	}

	linkViewsB, err := q.repository.FetchLinkViews(criteria, 31)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	var linkViews []LinkView
	for _, linkViewB := range linkViewsB {
		var linkView LinkView
		bsonBytes, _ := bson.Marshal(linkViewB)
		err = bson.Unmarshal(bsonBytes, &linkView)
		if err == nil {
			linkViews = append(linkViews, linkView)
		}
	}

	return &linkViews, nil
}

// QueryTodayTopLinksViews will build query based on given year, month and day.
// Then query data storage and return list of top LinkView of the current day.
func (q *Query) QueryTodayTopLinksViews(year string, month string, day string) (*[]LinkView, error) {
	criteria := bson.M{
		"$and": []bson.M{
			{"year": bson.M{"$eq": year}},
			{"month": bson.M{"$eq": month}},
			{"day": bson.M{"$eq": day}},
		},
	}

	linkViewsB, err := q.repository.FetchTopLinksViews(criteria, 10)
	if err != nil && err != mongo.ErrNoDocuments {
		return nil, err
	}

	var linkViews []LinkView
	for _, linkViewB := range linkViewsB {
		var linkView LinkView
		bsonBytes, _ := bson.Marshal(linkViewB)
		err = bson.Unmarshal(bsonBytes, &linkView)
		if err == nil {
			linkViews = append(linkViews, linkView)
		}
	}

	return &linkViews, nil
}
