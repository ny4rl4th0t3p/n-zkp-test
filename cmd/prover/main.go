package main

import (
	"context"
	"log"
	"log/slog"
	"time"

	"practical-case-test/config"
	"practical-case-test/internal/app"
	interactor "practical-case-test/internal/interactor/grpc"
)

// main initializes the client and performs the following actions:
// 1. Generates a random userName and userPassword.
// 2. Registers the user with the client using the generated userName and userPassword.
// 3. Logs in with the registered user credentials.
// 4. Prints the session ID of the successful login.
// 5. Sleeps for 60 seconds before ending the program.
func main() {
	cfg := config.LoadConfig()

	client, err := interactor.NewClient(
		cfg.VerifierURL,
		cfg,
		app.NewRegister(),
		app.NewCommitment(),
		app.NewComputeS(),
	)
	if err != nil {
		log.Fatalf("Failed to create client: %v", err)
	}
	defer func(client *interactor.AuthenticationClient) {
		_ = client.Close()
	}(client)

	var userName = app.RandString(10)

	userPassword, err := app.RandomPassword()
	if err != nil {
		slog.Error("error creating random password", "error", err)
		return
	}

	err = client.Register(context.Background(), userName, userPassword)
	if err != nil {
		slog.Error("register failed", "error", err)
		return
	}

	sessionID, err := client.Login(context.Background(), userName, userPassword)
	if err != nil {
		slog.Error("login failed", "error", err)
		slog.Info("sleeping for 30 seconds after error")
		time.Sleep(30 * time.Second)
		return
	}

	slog.Info("successfully logged", "session id", sessionID)

	slog.Info("sleeping for 60 seconds before killing prover")

	time.Sleep(60 * time.Second)
}
