package di

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/lalabuy948/linkopus/backend/config"
	"github.com/lalabuy948/linkopus/backend/pkg/linkopus/database"
	"github.com/lalabuy948/linkopus/backend/pkg/linkopus/service"

	"github.com/segmentio/nsq-go"

	"github.com/VictoriaMetrics/fastcache"
	"github.com/go-redis/cache/v8"
	"github.com/go-redis/redis/v8"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Container contains constructed linkopus services.
type Container struct {
	QueryService   *service.Query
	CommandService *service.Command
	FacadeService  *service.Facade
	WorkerService  *service.Worker
}

// Compile constructing service container and all necessary DIs.
func Compile(cfg *config.Config) *Container {
	mongoClient := connectMongoClient(cfg.MongoDBUrl)
	cacheClient := connectRedisCluster(cfg.RedisUrl)

	nqsClient := connectNQSProducer(cfg.NSQUrl)
	nqsConsumer := connectNQSConsumer(cfg.NSQUrl)

	cacheService := service.NewCache(cacheClient)

	repository := database.NewRepository(mongoClient)
	manager := database.NewManager(mongoClient)

	queryService := service.NewQuery(repository)
	commandService := service.NewCommand(manager, nqsClient)
	facadeService := service.NewFacade(queryService, commandService, cacheService)

	workerService := service.NewWorker(manager, nqsConsumer)

	return &Container{
		queryService,
		commandService,
		facadeService,
		workerService,
	}
}

func connectMongoClient(mongoUrl string) *mongo.Client {
	fmt.Println("connecting to db... ", mongoUrl)
	client, err := mongo.NewClient(options.Client().ApplyURI(mongoUrl))
	if err != nil {
		log.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	if err = client.Connect(ctx); err != nil {
		log.Fatal(err)
	}

	return client
}

func connectRedisCluster(redisUrl string) *cache.Cache {
	fmt.Println("connecting to redis... ", redisUrl)
	ring := redis.NewRing(&redis.RingOptions{
		Addrs: map[string]string{
			"redis": redisUrl,
		},
	})

	return cache.New(&cache.Options{
		Redis:        ring,
		LocalCache:   fastcache.New(100 << 20),
		StatsEnabled: false,
	})
}

func connectNQSProducer(nqsUrl string) *nsq.Producer {
	fmt.Println("connecting to nqs (producer)... ", nqsUrl)
	producer, err := nsq.StartProducer(nsq.ProducerConfig{
		Topic:   "views",
		Address: nqsUrl,
	})

	if err != nil {
		log.Fatal(err)
	}

	return producer
}

func connectNQSConsumer(nqsUrl string) *nsq.ConsumerConfig {
	return &nsq.ConsumerConfig{
		Topic:       "views",
		Channel:     "consumer1",
		Address:     nqsUrl,
		MaxInFlight: 250,
	}
}
