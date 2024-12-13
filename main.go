package main

import (
	"context"
	"crypto/tls"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
	"github.com/imrishuroy/legal-referral/api"
	"github.com/imrishuroy/legal-referral/chat"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	"github.com/imrishuroy/legal-referral/util"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/redis/go-redis/v9"
	"time"

	"github.com/rs/zerolog/log"
)

func main() {

	log.Info().Msg("Welcome to LegalReferral")

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Err(err).Msg("cannot load config : " + err.Error())
	}

	// db connection
	connPool, err := pgxpool.New(context.Background(), config.DBSource)

	if err != nil {
		fmt.Println("cannot connect to db:", err)
	}
	defer connPool.Close()

	store := db.NewStore(connPool)

	hub := chat.NewHub(store)
	go hub.Run()

	//setup producer
	conf := kafka.ConfigMap{
		// User-specific properties that you must set
		"bootstrap.servers": config.BootStrapServers,
		"sasl.username":     config.SASLUsername,
		"sasl.password":     config.SASLPassword,

		// Fixed properties
		"security.protocol": "SASL_SSL",
		"sasl.mechanisms":   "PLAIN",
		"acks":              "all"}

	producer, err := kafka.NewProducer(&conf)
	if err != nil {
		log.Error().Err(err).Msg("cannot create producer")
	}

	ctx := context.Background()

	rdb := GetRedisClient(config)

	pong, err := rdb.Ping(ctx).Result()
	if err != nil {
		log.Error().Err(err).Msg("cannot connect to redis")
	} else {
		log.Info().Msg("Connected to Redis with TLS: " + pong)
	}

	// api server setup
	server, err := api.NewServer(config, store, hub, producer, rdb)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server:")
	}

	//go server.CreateConsumer(context.Background())

	// start the server
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server:")
	}

}

func GetRedisClient(config util.Config) api.RedisClient {
	redisURL := fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort)
	if config.Env == "prod" {
		clusterClient := redis.NewClusterClient(&redis.ClusterOptions{
			Addrs:           []string{redisURL},
			Password:        "",
			PoolSize:        50,
			MinIdleConns:    20,
			DialTimeout:     3 * time.Second,
			ReadTimeout:     1 * time.Second,
			WriteTimeout:    1 * time.Second,
			PoolTimeout:     2 * time.Second,
			MaxRetries:      3,
			MinRetryBackoff: 8 * time.Millisecond,
			MaxRetryBackoff: 256 * time.Millisecond,
			TLSConfig: &tls.Config{
				InsecureSkipVerify: false,
			},
			ReadOnly:       false,
			RouteByLatency: true,  // Prioritize low-latency nodes
			RouteRandomly:  false, // Avoid random routing to improve predictability
		})
		return clusterClient
	} else {
		client := redis.NewClient(&redis.Options{
			Addr:     redisURL,
			Password: "",
			DB:       0,
		})
		return client
	}
}
