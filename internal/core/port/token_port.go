package port

// TokenService handles the generation and validation of authentication tokens.
type TokenService interface {
	// GenerateToken creates a new token for the specified user ID.
	GenerateToken(userID string) (string, error)

	// ValidateToken checks the validity of a token and returns the associated user ID if valid.
	ValidateToken(token string) (string, error)
}
