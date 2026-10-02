// Package security provides authentication, token generation, and security-related utilities.
package security

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/irfantrue/web-application-backend-go/internal/config"
	"github.com/irfantrue/web-application-backend-go/internal/core/entity"
	"github.com/irfantrue/web-application-backend-go/internal/core/port"
)

type jwtService struct {
	secretKey []byte
	expireHr  int
}

// CustomClaims represents the JWT claims containing the user ID and standard registered claims.
type CustomClaims struct {
	entity.TokenPayload
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
func (j *jwtService) GenerateAccessToken(payload *entity.TokenPayload) (*entity.TokenResponse, error) {
	expirationTime := time.Now().Add(time.Duration(j.expireHr) * time.Hour)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, CustomClaims{
		TokenPayload: *payload,
		ExpiresAt:    jwt.NewNumericDate(expirationTime),
		IssuedAt:     jwt.NewNumericDate(time.Now()),
	})

	tokenString, err := token.SignedString(j.secretKey)
	if err != nil {
		return nil, err
	}

	return &entity.TokenResponse{
		Token:     tokenString,
		ExpiresAt: expirationTime,
	}, nil
}

// ValidateToken checks the validity of a token and returns the associated user ID if valid.
func (j *jwtService) ValidateToken(tokenString string) (*entity.TokenPayload, error) {
	token, err := jwt.ParseWithClaims(tokenString, &CustomClaims{}, func(token *jwt.Token) (any, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("invalid token algorithm")
		}
		return j.secretKey, nil
	})
	if err != nil {
		return nil, err
	}

	claims, ok := token.Claims.(*CustomClaims)
	if ok && token.Valid {
		return &claims.TokenPayload, nil
	}

	return nil, errors.New("the token is invalid or has expired")
}
