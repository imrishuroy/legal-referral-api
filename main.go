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
	"github.com/rs/zerolog/log"
	"time"
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

	ctx := context.Background()
	rdb := redis.NewClusterClient(&redis.ClusterOptions{
		Addrs:    []string{redisURL},
		Password: "",

		// Increased pool size and idle connections for high-concurrency handling
		PoolSize:     50,
		MinIdleConns: 20,

		// Reduced timeouts for faster fail over in case of delays
		DialTimeout:  3 * time.Second,
		ReadTimeout:  1 * time.Second,
		WriteTimeout: 1 * time.Second,
		PoolTimeout:  2 * time.Second,

		// Connection idle and age management for better connection freshness
		//IdleCheckFrequency: 30 * time.Second,
		//IdleTimeout:        2 * time.Minute,
		//MaxConnAge:         5 * time.Minute,

		// Adjusted retry settings for optimal backoff strategy
		MaxRetries:      3,
		MinRetryBackoff: 8 * time.Millisecond,
		MaxRetryBackoff: 256 * time.Millisecond,

		// TLS configuration based on your security needs
		TLSConfig: &tls.Config{
			InsecureSkipVerify: false, // Better to set to false for production, assuming valid certificates
		},

		// Routing adjustments for balanced load with low latency focus
		ReadOnly:       false,
		RouteByLatency: true,  // Prioritize low-latency nodes
		RouteRandomly:  false, // Avoid random routing to improve predictability
	})

	//rdb := redis.NewClient(&redis.Options{
	//	//Addr: "localhost:6379",
	//	Addr:     redisURL,
	//	Password: "", // No password set
	//	DB:       0,  // Use default DB
	//	Protocol: 2,  // Connection protocol
	//})

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
