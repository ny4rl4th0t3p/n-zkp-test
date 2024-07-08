package main

import (
	"log"
	"net"

	"practical-case-test/config"
	"practical-case-test/internal/app"
	igrpc "practical-case-test/internal/interactor/grpc"
	interactor "practical-case-test/internal/interactor/proto"
	"practical-case-test/internal/repository/memory"

	"google.golang.org/grpc"
)

// main is the entry point of the application. It starts a gRPC server and registers
// the authentication server handlers. It also initializes the necessary dependencies, such as
// the authentication repository and the interactor. It uses the config loaded from LoadConfig
// function. The server listens on port 50051 for incoming connections.
func main() {
	listener, err := net.Listen("tcp", "0.0.0.0:50051")
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

	if err = s.Serve(listener); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
