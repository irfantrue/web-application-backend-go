// Package security provides authentication, token generation, and security-related utilities.
package security

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/irfantrue/web-application-backend-go/internal/config"
	"github.com/irfantrue/web-application-backend-go/internal/core/port"
)

type jwtService struct {
	secretKey []byte
	expireHr  int
}

// CustomClaims represents the JWT claims containing the user ID and standard registered claims.
type CustomClaims struct {
	UserID string `json:"user_id"`
	jwt.RegisteredClaims
}

// NewJwtService creates and returns a new instance of port.TokenService using JWT.
func NewJwtService(cfg *config.Config) port.TokenService {
	return &jwtService{
		secretKey: []byte(cfg.App.JwtSecret),
		expireHr:  cfg.App.JwtExpire,
	}
}

// GenerateToken creates a new token for the specified user ID.
func (j *jwtService) GenerateToken(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
		UserID:    userID,
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Duration(j.expireHr) * time.Hour)),
		IssuedAt:  jwt.NewNumericDate(time.Now()),
	})

	return token.SignedString(j.secretKey)
}

// ValidateToken checks the validity of a token and returns the associated user ID if valid.
func (j *jwtService) ValidateToken(tokenString string) (string, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token algorithm")
		}
		return j.secretKey, nil
	})
	if err != nil {
		return "", err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if ok && token.Valid {
		return claims.UserID, nil
	}

	return "", errors.New("the token is invalid or has expired")
}
