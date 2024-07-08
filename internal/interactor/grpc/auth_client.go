package grpc

import (
	"context"
	"fmt"
	"log/slog"
	"math/big"

	"practical-case-test/config"
	"practical-case-test/internal/app"
	interactor "practical-case-test/internal/interactor/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type AuthenticationClient struct {
	cc   *grpc.ClientConn
	cfg  *config.Config
	auth interactor.AuthClient
	re   app.RegisterExecuter
	co   app.CommitmentExecuter
	cs   app.ComputeSExecuter
}

func NewClient(address string, cfg *config.Config, re app.RegisterExecuter, co app.CommitmentExecuter, cs app.ComputeSExecuter) (*AuthenticationClient, error) {
	conn, err := grpc.NewClient(
		address,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to dial server %s, err: %w", address, err)
	}

	return &AuthenticationClient{
		cc:   conn,
		cfg:  cfg,
		auth: interactor.NewAuthClient(conn),
		re:   re,
		co:   co,
		cs:   cs,
	}, nil
}

// Register sends a register request to the authentication server for the given user credentials.
func (c *AuthenticationClient) Register(ctx context.Context, userName string,
	userPassword *big.Int) error {
	y1, y2, err := c.re.Exec(c.cfg, userPassword)
	if err != nil {
		return fmt.Errorf("could not calculate y1 and y2, err: %w", err)
	}
	_, err = c.auth.Register(ctx, &interactor.RegisterRequest{User: userName, Y1: y1.Int64(), Y2: y2.Int64()})
	if err != nil {
		return fmt.Errorf("register request failed for user %s, err: %w", userName, err)
	}

	slog.Info("registered user", "user", userName, "password", userPassword, "y1", y1, "y2", y2)

	return nil
}

// Login performs the login process for a user.
//
// The login process involves the following steps:
//  1. Generate data for commitment.
//  2. Send the commitment data to the server.
//  3. Execute the challenge response.
//  4. Verify authentication with the server.
//
// Upon successful authentication, the method returns the session ID.
// Otherwise, it returns an error.
func (c AuthenticationClient) Login(ctx context.Context, userName string, userPassword *big.Int) (string, error) {
	slog.Info("start login process")

	slog.Info("generating data for commitment")
	commitment, err := c.co.Exec(c.cfg)
	if err != nil {
		return "", err
	}

	challengeResp, err := c.auth.CreateAuthenticationChallenge(ctx, &interactor.AuthenticationChallengeRequest{
		User: userName,
		R1:   commitment.R1.Int64(),
		R2:   commitment.R2.Int64(),
	})
	if err != nil {
		return "", fmt.Errorf("create authentication challenge failed for user %s, err: %w", userName, err)
	}

	slog.Info("commitment sent successfully", "challenge response", challengeResp)

	slog.Info("processing challenge response")
	s, err := c.cs.Exec(c.cfg, userPassword, commitment.K, challengeResp)
	if err != nil {
		return "", err
	}

	slog.Info("verifying authentication with the server.")
	authResp, err := c.auth.VerifyAuthentication(ctx, &interactor.AuthenticationAnswerRequest{
		AuthId: challengeResp.GetAuthId(),
		S:      s.Int64(),
	})
	if err != nil {
		return "", fmt.Errorf("verify authentication failed for user %s, err: %w", userName, err)
	}

	slog.Info("user authenticated successfully", "session id", authResp.GetSessionId())

	return authResp.GetSessionId(), nil
}

// Close closes the client connection. If the connection is not nil,
// it calls the Close method on the underlying grpc.ClientConn.
// It returns nil if the connection is successfully closed or if the connection is nil.
func (c *AuthenticationClient) Close() error {
	if c.cc != nil {
		return c.cc.Close()
	}
	return nil
}
