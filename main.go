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
	log.Info().Msg("Connecting to Redis at " + redisURL)

	ctx := context.Background()
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

		ReadOnly: false,
		//RouteRandomly:  false,
		//RouteByLatency: false,
		RouteByLatency: true,
		RouteRandomly:  true,
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
		log.Info().Msg("Connected to Redis with TLS:" + pong)
	}

	//Set a key-value pair in Redis
	err = rdb.Set(ctx, "mykey", "myvalue", 0).Err()
	if err != nil {
		log.Error().Err(err).Msg("Failed to set key")
	}

	// Retrieve the value of the key
	val, err := rdb.Get(ctx, "mykey").Result()
	if err != nil {
		log.Error().Err(err).Msg("Failed to get key")
	}

	fmt.Println("mykey:", val)

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
