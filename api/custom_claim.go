package api

import "github.com/golang-jwt/jwt/v5"

// CustomClaims contains custom data we want from the token.
type CustomClaims struct {
	// Name         string `json:"name"`
	// Username     string `json:"username"`
	// ShouldReject bool   `json:"shouldReject,omitempty"`
	Scope string `json:"scope"`
	Email string `json:"email"`
}

func (c CustomClaims) GetExpirationTime() (*jwt.NumericDate, error) {
	//TODO implement me
	panic("implement me")
}

func (c CustomClaims) GetIssuedAt() (*jwt.NumericDate, error) {
	//TODO implement me
	panic("implement me")
}

func (c CustomClaims) GetNotBefore() (*jwt.NumericDate, error) {
	//TODO implement me
	panic("implement me")
}

func (c CustomClaims) GetIssuer() (string, error) {
	//TODO implement me
	panic("implement me")
}

func (c CustomClaims) GetSubject() (string, error) {
	//TODO implement me
	panic("implement me")
}

func (c CustomClaims) GetAudience() (jwt.ClaimStrings, error) {
	//TODO implement me
	panic("implement me")
}

func (c CustomClaims) GetEmail() (string, error) {
	return c.Email, nil
}
