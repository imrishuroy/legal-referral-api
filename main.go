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
	"github.com/redis/go-redis/v9"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func main() {

	log.Info().Msg("Welcome to LegalReferral")

	config, err := util.LoadConfig(".")
	if err != nil {
		// print the error message
		log.Fatal().Err(err).Msg("cannot load config : " + err.Error())
	}

	// db connection
	connPool, err := pgxpool.New(context.Background(), config.DBSource)

	if err != nil {
		fmt.Println("cannot connect to db:", err)
	}
	defer connPool.Close() // close db connection

	store := db.NewStore(connPool)

	//	redisClient := redis.NewRedis()

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

	redisURL := fmt.Sprintf("%s:%s", config.RedisHost, config.RedisPort)
	log.Info().Msg("Connecting to Redis at " + redisURL)

	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:        []string{redisURL},
		Password:     "",
		PoolSize:     10,
		MinIdleConns: 10,

		DialTimeout:  5 * time.Second,
		ReadTimeout:  3 * time.Second,
		WriteTimeout: 3 * time.Second,
		PoolTimeout:  4 * time.Second,

		//IdleCheckFrequency: 60 * time.Second,
		//IdleTimeout:        5 * time.Minute,
		//MaxConnAge:         0 * time.Second,

		MaxRetries:      10,
		MinRetryBackoff: 8 * time.Millisecond,
		MaxRetryBackoff: 512 * time.Millisecond,

		TLSConfig: &tls.Config{
			InsecureSkipVerify: true,
		},

		ReadOnly:       false,
		RouteRandomly:  false,
		RouteByLatency: false,
	})

	pong, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		log.Error().Err(err).Msg("cannot connect to redis")
	} else {
		log.Info().Msg("Connected to Redis with TLS:" + pong)
	}

	// api server setup
	server, err := api.NewServer(config, store, hub, producer)
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
