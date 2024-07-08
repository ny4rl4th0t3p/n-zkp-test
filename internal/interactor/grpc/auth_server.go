package grpc

import (
	"context"
	"fmt"
	"log/slog"

	"practical-case-test/config"
	"practical-case-test/internal/app"
	interactor "practical-case-test/internal/interactor/proto"
)

type AuthenticationServer struct {
	interactor.UnimplementedAuthServer
	cfg *config.Config
	ru  app.RegisterUserExecuter
	cac app.CreateAuthenticationChallengeExecuter
	va  app.VerifyAuthenticationExecuter
}

func NewAuthenticationServer(cfg *config.Config, ru app.RegisterUserExecuter, cac app.CreateAuthenticationChallengeExecuter,
	va app.VerifyAuthenticationExecuter) *AuthenticationServer {
	return &AuthenticationServer{cfg: cfg, ru: ru, cac: cac, va: va}
}

func (a *AuthenticationServer) Register(ctx context.Context, in *interactor.RegisterRequest) (*interactor.RegisterResponse, error) {
	user := in.GetUser()
	slog.Info("received registration request", "user", user)

	err := a.ru.Exec(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("failed to register user %q: %w", user, err)
	}
	return &interactor.RegisterResponse{}, nil
}

// CreateAuthenticationChallenge creates an authentication challenge for the user specified in the request.
// It logs the user ID of the user making the request and calls the Execute method of the `cac` (CreateAuthenticationChallengeExecuter) field of the AuthenticationServer struct.
// If there is an error executing the challenge, it returns the error.
// Otherwise, it creates an AuthenticationChallengeResponse with the AuthId and C values from the challenge, and returns it along with nil error.
func (a *AuthenticationServer) CreateAuthenticationChallenge(ctx context.Context, in *interactor.AuthenticationChallengeRequest) (*interactor.AuthenticationChallengeResponse, error) {
	userID := in.GetUser()
	slog.Info("received challenge request", "user", userID)

	challenge, err := a.cac.Exec(ctx, in)
	if err != nil {
		return nil, fmt.Errorf("user %s failed challenge: %w", userID, err)
	}

	return &interactor.AuthenticationChallengeResponse{
		AuthId: challenge.AuthID().String(),
		C:      challenge.C().Int64(),
	}, nil
}

// VerifyAuthentication verifies the authentication based on the provided AuthenticationAnswerRequest.
// It retrieves the authID from the request and logs the received request.
// It then executes the VerifyAuthenticationExecuter to verify the session.
// If an error occurs, it logs the failure, constructs an error message, and returns it.
// Otherwise, it constructs an AuthenticationAnswerResponse with the session ID and returns it.
func (a *AuthenticationServer) VerifyAuthentication(ctx context.Context, in *interactor.AuthenticationAnswerRequest) (*interactor.AuthenticationAnswerResponse, error) {
	authID := in.GetAuthId()
	slog.Info("received verify authentication", "authID", authID)

	session, err := a.va.Exec(ctx, a.cfg, in)
	if err != nil {
		slog.Error("failed to verify session", "authID", authID, "error", err)
		return nil, fmt.Errorf("failed to authenticate %s: %w", authID, err)
	}

	return &interactor.AuthenticationAnswerResponse{
		SessionId: session.ID().String(),
	}, nil
}
