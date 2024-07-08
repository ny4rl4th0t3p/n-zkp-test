package repository

import (
	"context"

	authDomain "practical-case-test/internal/domain/auth"

	"github.com/google/uuid"
)

type AuthRepository interface {
	StoreUserRegistration(ctx context.Context, userID authDomain.User) error
	GetUserRegistration(ctx context.Context, userID string) (*authDomain.User, error)
	StoreAuthenticationChallenge(ctx context.Context, challenge authDomain.Challenge) error
	GetAuthenticationChallenge(ctx context.Context, authID string) (*authDomain.Challenge, error)
	StoreSession(ctx context.Context, session authDomain.Session) error
	GetSession(ctx context.Context, userID string, sessionID uuid.UUID) (*authDomain.Session, error)
}
