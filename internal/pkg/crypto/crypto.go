// Package crypto provides cryptographic utilities such as secure string generation and hashing.
package crypto

import (
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"

	"github.com/ccoveille/go-safecast/v2"
	"github.com/irfantrue/web-application-backend-go/internal/core/port"
	"golang.org/x/crypto/argon2"
)

// Ensure CryptographyService implements port.CryptographyService at compile time.
var _ port.CryptographyService = (*CryptographyService)(nil)

type params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

// CryptographyService provides concrete implementations for cryptographic operations.
type CryptographyService struct {
	p *params
}

// NewCryptographyService creates and returns a new instance of CryptographyService.
func NewCryptographyService() *CryptographyService {
	return &CryptographyService{
		p: &params{
			memory:      64 * 1024, // 64 MB memory
			iterations:  3,         // 3 iteration CPU
			parallelism: 2,         // 2 threads/cores
			saltLength:  16,        // 16 bytes length salt
			keyLength:   32,        // 32 bytes length hash
		},
	}
}

const alnum = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789"

// GenerateSecureString creates a random, URL-safe alphanumeric string of the specified length.
func (s *CryptographyService) GenerateSecureString(length int) (string, error) {
	if length <= 0 {
		return "", errors.New("length must be positive")
	}

	buf := make([]byte, length)
	if _, err := rand.Read(buf); err != nil {
		return "", err
	}
	for i, b := range buf {
		buf[i] = alnum[int(b)%len(alnum)]
	}
	return string(buf), nil
}

// HashString generates a SHA-256 hash of the given input string and returns it as a hexadecimal string.
// It returns an empty string if the input is empty.
func (s *CryptographyService) HashString(input string) string {
	if input == "" {
		return ""
	}
	h := sha256.Sum256([]byte(input))
	return hex.EncodeToString(h[:])
}

// HashPassword hashes a password using Argon2id and returns the encoded PHC string.
func (s *CryptographyService) HashPassword(password string) (string, error) {
	salt := make([]byte, s.p.saltLength)
	if _, err := rand.Read(salt); err != nil {
		return "", err
	}

	hash := argon2.IDKey([]byte(password), salt, s.p.iterations, s.p.memory, s.p.parallelism, s.p.keyLength)

	return fmt.Sprintf(
		"$argon2id$v=%d$m=%d,t=%d,p=%d$%s$%s",
		argon2.Version, s.p.memory, s.p.iterations, s.p.parallelism,
		base64.RawStdEncoding.EncodeToString(salt),
		base64.RawStdEncoding.EncodeToString(hash),
	), nil
}

// ComparePassword verifies a password against an Argon2id-encoded hash.
func (s *CryptographyService) ComparePassword(hashedPassword, password string) error {
	vals := strings.Split(hashedPassword, "$")
	if len(vals) != 6 {
		return errors.New("invalid hash format")
	}
	if vals[1] != "argon2id" {
		return errors.New("unsupported algorithm, expected argon2id")
	}

	var version int
	if _, err := fmt.Sscanf(vals[2], "v=%d", &version); err != nil || version != argon2.Version {
		return errors.New("argon2 version mismatch")
	}

	var m, t, p int
	if _, err := fmt.Sscanf(vals[3], "m=%d,t=%d,p=%d", &m, &t, &p); err != nil {
		return errors.New("invalid argon2 parameters")
	}

	memory, err := safecast.Convert[uint32](m)
	if err != nil {
		return errors.New("invalid argon2 parameters")
	}
	iterations, err := safecast.Convert[uint32](t)
	if err != nil {
		return errors.New("invalid argon2 parameters")
	}
	parallelism, err := safecast.Convert[uint8](p)
	if err != nil {
		return errors.New("invalid argon2 parameters")
	}

	salt, err := base64.RawStdEncoding.DecodeString(vals[4])
	if err != nil {
		return errors.New("corrupted salt")
	}

	hash, err := base64.RawStdEncoding.DecodeString(vals[5])
	if err != nil {
		return errors.New("corrupted hash")
	}

	keyLength, err := safecast.Convert[uint32](len(hash))
	if err != nil {
		return errors.New("corrupted hash")
	}

	comparisonHash := argon2.IDKey([]byte(password), salt, iterations, memory, parallelism, keyLength)

	if subtle.ConstantTimeCompare(hash, comparisonHash) == 1 {
		return nil
	}
	return errors.New("incorrect password")
}
