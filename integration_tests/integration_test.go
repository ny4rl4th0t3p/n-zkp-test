//go:build integration

package integration_test

import (
	"context"
	"log"
	"math/big"
	"net"
	"testing"
	"time"

	"practical-case-test/internal/repository/memory"

	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"

	"practical-case-test/config"
	"practical-case-test/internal/app"
	igrpc "practical-case-test/internal/interactor/grpc"
	interactor "practical-case-test/internal/interactor/proto"
)

// runServer starts the server.
//
// It listens for incoming connections on the specified address and creates a new gRPC server.
// It creates an in-memory authentication repository and initializes the required application executer instances.
// Finally, it registers the AuthenticationServer with the gRPC server and starts serving incoming requests.
//
// Parameters:
// - address: The address to listen on for incoming connections (e.g., "localhost:50051").
func runServer(address string) {
	lis, err := net.Listen("tcp", address)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	cfg := config.LoadConfig()
	s := grpc.NewServer()

	ar := memory.NewInMemAuthRepository()
	ru := app.NewRegisterUser(ar)
	ca := app.NewCreateAuthenticationChallenge(ar)
	va := app.NewVerifyAuthentication(ar)

	interactor.RegisterAuthServer(s, igrpc.NewAuthenticationServer(cfg, ru, ca, va))

	if err = s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}

// Test_FuncTestScenario1 tests the successful register and login scenario.
//
// It first starts the server by calling the runServer function in a separate goroutine.
// Then, it creates a new client using the igrpc.NewClient function.
// It registers a user using the client's Register method.
// Finally, it tests the Login method using the registered user's credentials.
func Test_FuncTestScenario1(t *testing.T) {
	go runServer("localhost:50051")
	time.Sleep(time.Second)

	cfg := config.LoadConfig()

	client, err := igrpc.NewClient(
		"localhost:50051",
		cfg,
		app.NewRegister(),
		app.NewCommitment(),
		app.NewComputeS(),
	)
	require.NoError(t, err)

	userName := "testUser1"
	userPassword := big.NewInt(123)

	err = client.Register(context.Background(), userName, userPassword)
	require.NoError(t, err)

	_, err = client.Login(context.Background(), userName, userPassword)
	require.NoError(t, err)

	err = client.Close()
	require.NoError(t, err)
}

// Test_FuncTestScenario2 tests the scenario of logging in with the wrong password.
//
// It starts the server by calling the runServer function in a goroutine.
// Then, it creates a client and registers a user with the correct password.
// Next, it tries to login with the same user but with a wrong password, expecting an error.
// Finally, it closes the client connection.
func Test_FuncTestScenario2(t *testing.T) {
	go runServer("localhost:50052")
	time.Sleep(time.Second)

	cfg := config.LoadConfig()

	client, err := igrpc.NewClient(
		"localhost:50051",
		cfg,
		app.NewRegister(),
		app.NewCommitment(),
		app.NewComputeS(),
	)
	require.NoError(t, err)

	userName := "testUser2"
	correctPassword := big.NewInt(456)
	wrongPassword := big.NewInt(3)

	err = client.Register(context.Background(), userName, correctPassword)
	require.NoError(t, err)

	_, err = client.Login(context.Background(), userName, wrongPassword)
	require.Error(t, err)

	err = client.Close()
	require.NoError(t, err)
}

func Test_FuncTestScenario3(t *testing.T) {
	go runServer("localhost:50053")
	time.Sleep(time.Second)

	cfg := config.LoadConfig()
	cfg.Q = big.NewInt(222)

	client, err := igrpc.NewClient(
		"localhost:50051",
		cfg,
		app.NewRegister(),
		app.NewCommitment(),
		app.NewComputeS(),
	)
	require.NoError(t, err)

	userName := "testUser3"
	correctPassword := big.NewInt(456)

	err = client.Register(context.Background(), userName, correctPassword)
	require.NoError(t, err)

	_, err = client.Login(context.Background(), userName, correctPassword)
	require.Error(t, err)

	err = client.Close()
	require.NoError(t, err)
}

func Test_FuncTestScenario4(t *testing.T) {
	go runServer("localhost:50054")
	time.Sleep(time.Second)

	cfg := config.LoadConfig()
	cfg.H = big.NewInt(222)

	client, err := igrpc.NewClient(
		"localhost:50051",
		cfg,
		app.NewRegister(),
		app.NewCommitment(),
		app.NewComputeS(),
	)
	require.NoError(t, err)

	userName := "testUser4"
	correctPassword := big.NewInt(456)

	err = client.Register(context.Background(), userName, correctPassword)
	require.NoError(t, err)

	_, err = client.Login(context.Background(), userName, correctPassword)
	require.Error(t, err)

	err = client.Close()
	require.NoError(t, err)
}

func Test_FuncTestScenario5(t *testing.T) {
	go runServer("localhost:50055")
	time.Sleep(time.Second)

	cfg := config.LoadConfig()
	cfg.G = big.NewInt(1)

	client, err := igrpc.NewClient(
		"localhost:50051",
		cfg,
		app.NewRegister(),
		app.NewCommitment(),
		app.NewComputeS(),
	)
	require.NoError(t, err)

	userName := "testUser5"
	correctPassword := big.NewInt(456)

	err = client.Register(context.Background(), userName, correctPassword)
	require.NoError(t, err)

	_, err = client.Login(context.Background(), userName, correctPassword)
	require.Error(t, err)

	err = client.Close()
	require.NoError(t, err)
}
