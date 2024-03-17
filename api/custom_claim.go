package api

// CustomClaims contains custom data we want from the token.
type CustomClaims struct {
	// Name         string `json:"name"`
	// Username     string `json:"username"`
	// ShouldReject bool   `json:"shouldReject,omitempty"`
	Scope string `json:"scope"`
}
