package port

// CryptographyService provides methods for cryptographic operations.
type CryptographyService interface {
	// GenerateSecureString creates a random secure string of the specified length.
	GenerateSecureString(length int) (string, error)

	// HashString generates a cryptographic hash of the given input string.
	HashString(input string) string

	// HashPassword hashes a password using Argon2id and returns the encoded PHC string.
	HashPassword(password string) (string, error)

	// ComparePassword verifies a password against an Argon2id-encoded hash.
	ComparePassword(hashedPassword, password string) error
}
