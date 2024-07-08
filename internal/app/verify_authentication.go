package app

import (
	"context"
	"errors"
	"log/slog"
	"time"

	"practical-case-test/config"
	"practical-case-test/internal/domain/auth"
	interactor "practical-case-test/internal/interactor/proto"
	"practical-case-test/internal/repository"

	"github.com/google/uuid"
)

// VerifyAuthenticationExecuter is an interface that defines the contract for executing
// the verification of authentication information.
type VerifyAuthenticationExecuter interface {
	Exec(ctx context.Context, cfg *config.Config, req *interactor.AuthenticationAnswerRequest) (*auth.Session, error)
}

// VerifyAuthentication is a type that is responsible for verifying the authentication
// of a user using an AuthRepository.
type VerifyAuthentication struct {
	ar repository.AuthRepository
}

// NewVerifyAuthentication creates a new instance of VerifyAuthenticationExecuter
// with the provided AuthRepository
func NewVerifyAuthentication(ar repository.AuthRepository) VerifyAuthenticationExecuter {
	return &VerifyAuthentication{ar: ar}
}

// Exec retrieves the authentication challenge for the given authID,
// verifies the user's response, and creates a new session for the user.
// It returns the newly created session, or an error if any operation fails.
func (va VerifyAuthentication) Exec(ctx context.Context, cfg *config.Config,
	req *interactor.AuthenticationAnswerRequest) (*auth.Session, error) {
	authID := req.GetAuthId()
	s := req.GetS()

	challenge, err := va.ar.GetAuthenticationChallenge(ctx, authID)
	if err != nil {
		return nil, err
	}

	slog.Info("challenge loaded", "challenge", challenge)

	user, err := va.ar.GetUserRegistration(ctx, challenge.UserID())
	if err != nil {
		return nil, err
	}

	if ok := verifyS(cfg, challenge, user, s); !ok {
		return nil, errors.New("verification failed, Invalid s")
	}

	session, err := auth.NewSession(uuid.New(), user.UserID(), time.Now().Unix())
	if err != nil {
		return nil, err
	}

	err = va.ar.StoreSession(ctx, *session)
	if err != nil {
		return nil, err
	}

	slog.Info("session initiated", "user", user.UserID(), "session", session.ID())

	return session, nil
}
