package queue

import (
	"encoding/json"
	"fmt"
	"strings"

	"github.com/lalabuy948/linkopus/backend/pkg/linkopus/entity"

	"github.com/lalabuy948/linkopus/backend/pkg/linkopus/data"

	"go.mongodb.org/mongo-driver/bson"

	"github.com/segmentio/nsq-go"
)

// Worker holds write manager and consumer config.
type Worker struct {
	manager        *data.Manager
	consumerConfig *nsq.ConsumerConfig
}

// NewWorker returns instance of the background worker.
func NewWorker(m *data.Manager, c *nsq.ConsumerConfig) *Worker {
	return &Worker{m, c}
}

// Consume starts consuming messages from message broker.
func (w *Worker) Consume() {
	fmt.Println("start consuming...")
	consumer, _ := nsq.StartConsumer(*w.consumerConfig)

	for msg := range consumer.Messages() {
		var linkVew entity.LinkView

		err := json.Unmarshal(msg.Body[:], &linkVew)
		if err == nil {
			if err := w.updateLinkView(&linkVew); err != nil {
				fmt.Println(err)
			}
		}

		msg.Finish()
	}
}

// updateLinkView upsert views count by link and date.
func (w *Worker) updateLinkView(linkVew *entity.LinkView) error {
	date := linkVew.Date
	dateSplit := strings.Split(date, "-")
	year := dateSplit[0]
	month := dateSplit[1]
	day := dateSplit[2]

	criteria := bson.D{
		{Key: "link", Value: linkVew.Link},
		{Key: "linkHash", Value: linkVew.LinkHash},
		{Key: "date", Value: date},
		{Key: "year", Value: year},
		{Key: "month", Value: month},
		{Key: "day", Value: day},
	}

	action := bson.D{
		{Key: "$inc", Value: bson.D{
			{Key: "amount", Value: 1},
		},
		},
	}

	return w.manager.UpdateLinkView(criteria, action)
}
