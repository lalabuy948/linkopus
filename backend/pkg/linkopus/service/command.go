package service

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/lalabuy948/linkopus/backend/pkg/linkopus/database"

	"github.com/segmentio/nsq-go"
	"go.mongodb.org/mongo-driver/bson"
)

// Command represents command service. Holds manager for synchronous writes and producer for asynchronous message passing.
type Command struct {
	manager  *database.Manager
	producer *nsq.Producer
}

// NewCommand return instance of linkopus Command service.
func NewCommand(m *database.Manager, p *nsq.Producer) *Command {
	return &Command{m, p}
}

// SYNC

// saveLinkMap stores link and link hash into the storage.
func (c *Command) saveLinkMap(link string, linkHash string) error {
	date := time.Now().Format("2006-01-02")
	dateSplit := strings.Split(date, "-")
	year := dateSplit[0]
	month := dateSplit[1]
	day := dateSplit[2]

	// this is incorrect way of doing commands
	return c.manager.InsertLinkMap(bson.D{
		{Key: "link", Value: link},
		{Key: "linkHash", Value: linkHash},
		{Key: "date", Value: date},
		{Key: "year", Value: year},
		{Key: "month", Value: month},
		{Key: "day", Value: day},
	})
}

// DeleteLinkMapByLink deletes LinkMap from the the storage by given link.
func (c *Command) DeleteLinkMapByLink(link string) error {
	return c.manager.DeleteLinkMap(bson.M{
		"link": link,
	})
}

// DeleteLinkMapByTime deletes LinkMap from the the storage by given year, month, day.
func (c *Command) DeleteLinkMapByTime(year string, month string, day string) error {
	return c.manager.DeleteLinkMap(bson.M{
		"$and": []bson.M{
			{"year": bson.M{"$eq": year}},
			{"month": bson.M{"$lte": month}},
			{"day": bson.M{"$lte": day}},
		},
	})
}

// DeleteLinkViewByLink deletes LinkView from the the storage by given link.
func (c *Command) DeleteLinkViewByLink(link string) error {
	return c.manager.DeleteLinkView(bson.M{
		"link": link,
	})
}

// DeleteLinkViewByLink deletes LinkView from the the storage by given year, month, day.
func (c *Command) DeleteLinkViewByTime(year string, month string, day string) error {
	return c.manager.DeleteLinkView(bson.M{
		"$and": []bson.M{
			{"year": bson.M{"$eq": year}},
			{"month": bson.M{"$lte": month}},
			{"day": bson.M{"$lte": day}},
		},
	})
}

// ASYNC

// SaveLinkView adds LinkView to the storage if not exists, otherwise will increment the amount of views.
func (c *Command) SaveLinkView(link string) {
	todayDate := time.Now().Format("2006-01-02")

	messageBody, _ := json.Marshal(&LinkView{
		Link: link,
		Date: todayDate,
	})

	_ = c.producer.Publish(messageBody) //nolint
}
