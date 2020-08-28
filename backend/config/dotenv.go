package config

import (
	"fmt"
	"log"

	"github.com/lalabuy948/genvutils"
)

// Config contains env values.
type Config struct {
	Environment string `genv:"ENVIRONMENT, DEV"`
	ServerPort  string `genv:"SERVER_PORT, 8080"`
	RedisUrl    string `genv:"REDIS_URL, 6379"`
	NSQUrl      string `genv:"NSQ_URL, localhost:4150"`
	MongoDBUrl  string `genv:"MONGO_DB_URL, mongodb://localhost:27017"`
}

// ParseEnv firstly loads env from dotenv.
// Secondly checks env from cli.
// Finally adds default values if not found in first two steps and return constructed config.
func ParseEnv() *Config {
	if err := genvutils.Load(); err != nil {
		fmt.Println(err)
	}

	var cfg Config
	if err := genvutils.Parse(&cfg); err != nil {
		log.Fatal("Error parsing dotenv")
	}

	return &cfg
}
