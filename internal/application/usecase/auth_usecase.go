package usecase

import (
	"context"
	"errors"
	"fmt"
	"log"
	"time"

	"github.com/bytedance/sonic"
	"github.com/google/uuid"

	"github.com/irfantrue/web-application-backend-go/internal/application/dto"
	"github.com/irfantrue/web-application-backend-go/internal/core/entity"
	"github.com/irfantrue/web-application-backend-go/internal/core/port"
)

var (
	// ErrInvalidCredentials is returned when login fails.
	ErrInvalidCredentials = errors.New("invalid email or password")

	// ErrEmailExists is returned when registration uses a taken email.
	ErrEmailExists = errors.New("email already registered")

	// ErrTokenExpired is returned when the refresh token has passed its expiry.
	ErrTokenExpired = errors.New("token has expired")

	// ErrTokenRevoked is returned when the token or its family has been revoked.
	ErrTokenRevoked = errors.New("token has been revoked")

	// ErrWrongPassword is returned when current password doesn't match.
	ErrWrongPassword = errors.New("current password is incorrect")

	// ErrResetTokenInvalid is returned when the reset token is invalid or expired.
	ErrResetTokenInvalid = errors.New("reset token is invalid or expired")
)

// AuthUsecase defines the business logic interface for auth operations.
type AuthUsecase interface {
	Register(ctx context.Context, req dto.RegisterRequest) error
	Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error)
	Logout(ctx context.Context, refreshToken string) error
	Refresh(ctx context.Context, refreshToken string) (*dto.AuthResponse, error)
	ForgotPassword(ctx context.Context, email string) error
	ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) error
	ChangePassword(ctx context.Context, userID uuid.UUID, req dto.ChangePasswordRequest) error
	Me(ctx context.Context, userID uuid.UUID) (*dto.MeResponse, error)
}

type authUsecaseImpl struct {
	uow          port.UnitOfWork
	userRepo     port.UserRepository
	refreshRepo  port.RefreshTokenSessionRepository
	crypto       port.CryptographyService
	tokenService port.TokenService
	broker       port.MessageBroker
}

// NewAuthUsecase creates a new instance of AuthUsecase.
func NewAuthUsecase(
	uow port.UnitOfWork,
	userRepo port.UserRepository,
	refreshRepo port.RefreshTokenSessionRepository,
	crypto port.CryptographyService,
	tokenService port.TokenService,
	broker port.MessageBroker,
) AuthUsecase {
	return &authUsecaseImpl{
		uow:          uow,
		userRepo:     userRepo,
		refreshRepo:  refreshRepo,
		crypto:       crypto,
		tokenService: tokenService,
		broker:       broker,
	}
}

func (u *authUsecaseImpl) Register(ctx context.Context, req dto.RegisterRequest) error {
	_, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err == nil {
		return ErrEmailExists
	}

	hashedPassword, err := u.crypto.HashPassword(req.Password)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	user, err := entity.NewUser(req.Email, hashedPassword)
	if err != nil {
		return fmt.Errorf("failed to create user entity: %w", err)
	}

	if err := u.userRepo.Create(ctx, user); err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	return nil
}

func (u *authUsecaseImpl) Login(ctx context.Context, req dto.LoginRequest) (*dto.AuthResponse, error) {
	user, err := u.userRepo.FindByEmail(ctx, req.Email)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	if err := u.crypto.ComparePassword(req.Password, user.Password); err != nil {
		return nil, ErrInvalidCredentials
	}

	var res *dto.AuthResponse
	err = u.uow.Execute(ctx, func(txCtx context.Context) error {
		rawToken, err := u.createRefreshSession(txCtx, user.ID)
		if err != nil {
			return err
		}

		// Generate access token
		payloadToken := &entity.TokenPayload{
			UserID:       user.ID.String(),
			Email:        user.Email,
			TokenVersion: user.TokenVersion,
		}

		accessToken, err := u.tokenService.GenerateAccessToken(payloadToken)
		if err != nil {
			return fmt.Errorf("generate access token: %w", err)
		}

		if err := u.userRepo.UpdateLastLogin(txCtx, user.ID, time.Now()); err != nil {
			return fmt.Errorf("update last login: %w", err)
		}

		res = dto.NewAuthResponse(accessToken.Token, rawToken, time.Until(accessToken.ExpiresAt))

		return nil
	})

	return res, err
}

func (u *authUsecaseImpl) Logout(ctx context.Context, refreshToken string) error {
	hash := u.crypto.HashString(refreshToken)

	session, err := u.refreshRepo.FindByTokenHash(ctx, hash)
	if err != nil {
		return nil // already gone, idempotent
	}

	now := time.Now()
	return u.refreshRepo.RevokeByID(ctx, session.ID, now)
}

func (u *authUsecaseImpl) Refresh(ctx context.Context, refreshToken string) (*dto.AuthResponse, error) {
	hash := u.crypto.HashString(refreshToken)

	session, err := u.refreshRepo.FindByTokenHash(ctx, hash)
	if err != nil {
		return nil, ErrInvalidCredentials
	}

	// Family already revoked
	if session.RevokedAt != nil {
		return nil, ErrTokenRevoked
	}

	var res *dto.AuthResponse
	err = u.uow.Execute(ctx, func(txCtx context.Context) error {
		// REUSE DETECTION
		if session.RotatedAt != nil {
			if err := u.refreshRepo.RevokeFamily(txCtx, session.FamilyID); err != nil {
				return fmt.Errorf("revoke by family: %w", err)
			}

			if err := u.refreshRepo.RevokeByUserID(txCtx, session.ID); err != nil {
				return fmt.Errorf("revoke by user id: %w", err)
			}

			if err := u.userRepo.IncrementTokenVersion(txCtx, session.UserID); err != nil {
				return fmt.Errorf("increment token version: %w", err)
			}

			return ErrTokenRevoked
		}

		// Expiry check
		if time.Now().After(session.ExpiresAt) {
			return ErrTokenExpired
		}

		// Rotate
		rawToken, err := u.createRefreshSession(txCtx, session.UserID)
		if err != nil {
			return err
		}

		// Mark old as rotated
		now := time.Now()
		if err := u.refreshRepo.MarkRotated(txCtx, session.ID, now); err != nil {
			return fmt.Errorf("mark rotated: %w", err)
		}

		// Generate new access token
		user, err := u.userRepo.FindByID(txCtx, session.UserID)
		if err != nil {
			return err
		}

		tokenPayload := &entity.TokenPayload{
			UserID:       user.ID.String(),
			Email:        user.Email,
			TokenVersion: user.TokenVersion,
		}

		accessToken, err := u.tokenService.GenerateAccessToken(tokenPayload)
		if err != nil {
			return fmt.Errorf("generate access token: %w", err)
		}

		res = dto.NewAuthResponse(accessToken.Token, rawToken, time.Until(accessToken.ExpiresAt))

		return nil
	})

	return res, err
}

func (u *authUsecaseImpl) ForgotPassword(ctx context.Context, email string) error {
	user, err := u.userRepo.FindByEmail(ctx, email)
	if err != nil {
		// Don't reveal whether email exists
		return nil
	}

	resetToken, err := u.crypto.GenerateSecureString(32)
	if err != nil {
		return fmt.Errorf("generate reset token: %w", err)
	}

	resetHash := u.crypto.HashString(resetToken)
	expiresAt := time.Now().Add(15 * time.Minute)

	if err := u.userRepo.SetResetToken(ctx, user.ID, resetHash, expiresAt); err != nil {
		return fmt.Errorf("set reset token: %w", err)
	}

	// Send email (fire-and-forget or queue)
	eventData := map[string]any{
		"event":      "reset_password",
		"url":        fmt.Sprintf("https://yourapp.com/reset-password?token=%s", resetToken),
		"user_email": user.Email,
		"timestamp":  time.Now(),
	}

	eventBytes, err := sonic.Marshal(eventData)
	if err != nil {
		return fmt.Errorf("failed marshaling event data reset password: %v", err)
	}

	go func() {
		bgCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
		defer cancel()

		err := u.broker.Publish(bgCtx, "", "notification.reset_password", eventBytes)
		if err != nil {
			log.Printf("Failed to publish product_created event: %v", err)
		}
	}()

	return nil
}

func (u *authUsecaseImpl) ResetPassword(ctx context.Context, req dto.ResetPasswordRequest) error {
	resetHash := u.crypto.HashString(req.Token)

	user, err := u.userRepo.FindByResetToken(ctx, resetHash)
	if err != nil {
		return ErrResetTokenInvalid
	}

	if user.ResetTokenExpiresAt == nil || time.Now().After(*user.ResetTokenExpiresAt) {
		return ErrResetTokenInvalid
	}

	hashedPassword, err := u.crypto.HashPassword(req.Password)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	return u.uow.Execute(ctx, func(txCtx context.Context) error {
		if err := u.userRepo.UpdatePassword(txCtx, user.ID, hashedPassword); err != nil {
			return fmt.Errorf("update password: %w", err)
		}

		// Clear reset token
		if err := u.userRepo.ClearResetToken(txCtx, user.ID); err != nil {
			return fmt.Errorf("clear reset token: %w", err)
		}

		// Revoke all sessions (password changed)
		if err := u.refreshRepo.RevokeByUserID(txCtx, user.ID); err != nil {
			return fmt.Errorf("revoke by user id: %w", err)
		}

		if err := u.userRepo.IncrementTokenVersion(ctx, user.ID); err != nil {
			return fmt.Errorf("increment token version: %w", err)
		}

		return nil
	})
}

func (u *authUsecaseImpl) ChangePassword(ctx context.Context, userID uuid.UUID, req dto.ChangePasswordRequest) error {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return err
	}

	if err := u.crypto.ComparePassword(req.CurrentPassword, user.Password); err != nil {
		return ErrWrongPassword
	}

	hashedPassword, err := u.crypto.HashPassword(req.NewPassword)
	if err != nil {
		return fmt.Errorf("hash password: %w", err)
	}

	return u.uow.Execute(ctx, func(txCtx context.Context) error {
		if err := u.userRepo.UpdatePassword(txCtx, user.ID, hashedPassword); err != nil {
			return fmt.Errorf("update password: %w", err)
		}

		// Revoke all sessions
		if err := u.refreshRepo.RevokeByUserID(txCtx, user.ID); err != nil {
			return fmt.Errorf("revoke by user id: %w", err)
		}

		if err := u.userRepo.IncrementTokenVersion(txCtx, user.ID); err != nil {
			return fmt.Errorf("increment token version: %w", err)
		}

		return nil
	})
}

func (u *authUsecaseImpl) Me(ctx context.Context, userID uuid.UUID) (*dto.MeResponse, error) {
	user, err := u.userRepo.FindByID(ctx, userID)
	if err != nil {
		return nil, err
	}

	return dto.NewMeResponse(user), nil
}

// createRefreshSession generates a new refresh token and persists the session.
// It returns the raw token (for the client) and its expiry.
func (u *authUsecaseImpl) createRefreshSession(ctx context.Context, userID uuid.UUID) (string, error) {
	rawToken, err := u.crypto.GenerateSecureString(32)
	if err != nil {
		return "", fmt.Errorf("generate refresh token: %w", err)
	}

	expiresAt := time.Now().Add(30 * 24 * time.Hour)
	tokenHash := u.crypto.HashString(rawToken)

	session := &entity.RefreshTokenSession{
		ID:        uuid.New(),
		FamilyID:  uuid.New(),
		UserID:    userID,
		TokenHash: tokenHash,
		ExpiresAt: expiresAt,
		CreatedAt: time.Now(),
	}

	if err := u.refreshRepo.Create(ctx, session); err != nil {
		return "", fmt.Errorf("create refresh session: %w", err)
	}

	return rawToken, nil
}
