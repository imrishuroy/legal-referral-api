package util

import (
	"github.com/google/uuid"
)

func GenerateUUID() (string, error) {
	// Generate a random UUID
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	// Convert the UUID to string and return
	return id.String(), nil
}
