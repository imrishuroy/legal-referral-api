package util

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/valkey-io/valkey-go"
)

// CreateValkeyClient creates either a real Valkey client for production or a nil client for local development
func CreateValkeyClient(config Config) (valkey.Client, error) {
	// Check if we're in local development mode
	if config.Env == "local" || config.Env == "dev" || config.Env == "development" {
		log.Println("Running in local mode - skipping Valkey connection")
		return nil, nil
	}

	// For production, create real Valkey client
	valkeyURL := fmt.Sprintf("%s:%s", config.ValKeyHost, config.ValKeyPort)
	log.Printf("Connecting to Valkey at: %s", valkeyURL)

	vkClient, err := valkey.NewClient(valkey.ClientOption{
		InitAddress: []string{valkeyURL},
		Password:    "",
		// Remove TLS for testing - you can add it back if needed
		// TLSConfig: &tls.Config{
		// 	InsecureSkipVerify: false,
		// },
		DisableCache: true,
	})

	if err != nil {
		return nil, fmt.Errorf("cannot create valkey client: %w", err)
	}

	// Test the connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Simple ping test
	result := vkClient.Do(ctx, vkClient.B().Ping().Build())
	if result.Error() != nil {
		vkClient.Close()
		return nil, fmt.Errorf("cannot connect to valkey: %w", result.Error())
	}

	log.Println("Successfully connected to Valkey")
	return vkClient, nil
}
