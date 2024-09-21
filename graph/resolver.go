package graph

import (
	db "github.com/imrishuroy/legal-referral/db/sqlc"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

const (
	authorizationHeaderKey  = "authorization"
	authorizationTypeBearer = "bearer"
	authorizationPayloadKey = "authorization_payload"
)

type Resolver struct {
	Store db.Store
}
