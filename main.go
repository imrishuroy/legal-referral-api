package main

import (
	"context"
	"fmt"
	"github.com/imrishuroy/legal-referral/api"
	"github.com/imrishuroy/legal-referral/chat"
	db "github.com/imrishuroy/legal-referral/db/sqlc"
	"github.com/imrishuroy/legal-referral/util"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

func main() {

	log.Info().Msg("Welcome to LegalReferral")

	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal().Msg("cannot connect to db:")
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

	// api server setup
	server, err := api.NewServer(config, store, hub)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server:")
	}

	// start the server
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal().Err(err).Msg("cannot create server:")
	}
}
