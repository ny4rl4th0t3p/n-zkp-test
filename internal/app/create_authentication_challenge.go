package app

import (
	"context"
	"crypto/rand"
	"log/slog"
	"math"
	"math/big"
	"time"

	"practical-case-test/internal/domain/auth"
	interactor "practical-case-test/internal/interactor/proto"
	"practical-case-test/internal/repository"
)

// CreateAuthenticationChallengeExecuter is an interface that defines the method for executing the creation of an authentication challenge.
type CreateAuthenticationChallengeExecuter interface {
	Exec(ctx context.Context, req *interactor.AuthenticationChallengeRequest) (*auth.Challenge, error)
}

// CreateAuthenticationChallenge is a type responsible for creating an authentication challenge.
type CreateAuthenticationChallenge struct {
	ar repository.AuthRepository
}

// NewCreateAuthenticationChallenge creates a new instance of CreateAuthenticationChallengeExecuter
// with the provided AuthRepository as a dependency. It returns a pointer to the CreateAuthenticationChallenge struct.
func NewCreateAuthenticationChallenge(ar repository.AuthRepository) CreateAuthenticationChallengeExecuter {
	return &CreateAuthenticationChallenge{ar: ar}
}

// Exec creates an Authentication Challenge for a user based on the provided request.
// It generates a random challenge value and logs it.
// The challenge request is then validated and stored in the repository.
// The generated challenge and any error that occurs during the process are returned.
// If an error occurs during the process, the returned challenge will be nil.
func (ru CreateAuthenticationChallenge) Exec(ctx context.Context, req *interactor.AuthenticationChallengeRequest) (
	*auth.Challenge,
	error,
) {
	userID := req.GetUser()

	slog.Info("creating authentication challenge for user", "user", userID)

	_, err := ru.ar.GetUserRegistration(ctx, userID)
	if err != nil {
		return nil, err
	}

	c, _ := rand.Int(rand.Reader, big.NewInt(math.MaxInt16))
	c.Add(c, big.NewInt(1))

	slog.Info("random challenge generated", "c", c)

	slog.Info("challenge request", "request", req)

	challenge, err := auth.NewChallenge(c, userID, req.GetR1(), req.GetR2(), time.Now().Unix())
	if err != nil {
		return nil, err
	}

	err = ru.ar.StoreAuthenticationChallenge(ctx, *challenge)
	if err != nil {
		return nil, err
	}

	return challenge, nil
}
