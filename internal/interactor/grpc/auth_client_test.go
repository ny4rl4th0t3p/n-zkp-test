package grpc

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"practical-case-test/config"
	"practical-case-test/internal/app"
	interactor "practical-case-test/internal/interactor/proto"
)

func TestAuthenticationClient_Login(t *testing.T) {
	type args struct {
		ctx          context.Context
		userName     string
		userPassword *big.Int
	}
	tests := []struct {
		name    string
		args    args
		cfg     *config.Config
		auth    *MockAuthClient
		re      *MockRegisterExecuter
		co      *MockCommitmentExecuter
		cs      *MockComputeSExecuter
		wantErr bool
	}{
		{
			name: "Test Case 1: Successful Login",
			auth: &MockAuthClient{
				AuthenticationChallengeResponse: &interactor.AuthenticationChallengeResponse{AuthId: "authId", C: 0},
				AuthenticationAnswerResponse:    &interactor.AuthenticationAnswerResponse{SessionId: "sessionId"},
			},
			co: &MockCommitmentExecuter{
				Result: &app.CommitmentResult{R1: big.NewInt(1), R2: big.NewInt(1), K: big.NewInt(1)},
			},
			cs: &MockComputeSExecuter{
				Result: big.NewInt(1),
			},
			cfg: &config.Config{},
			args: args{
				ctx:          context.Background(),
				userName:     "test",
				userPassword: big.NewInt(12345),
			},
			wantErr: false,
		},
		{
			name: "Test Case 2: Failed Login due to Authentication Challenge error",
			auth: &MockAuthClient{
				AuthenticationChallengeError: errors.New("authentication challenge error"),
			},
			co: &MockCommitmentExecuter{
				Result: &app.CommitmentResult{R1: big.NewInt(1), R2: big.NewInt(1), K: big.NewInt(1)},
			},
			cs: &MockComputeSExecuter{
				Result: big.NewInt(1),
			},
			cfg: &config.Config{},
			args: args{
				ctx:          context.Background(),
				userName:     "test",
				userPassword: big.NewInt(12345),
			},
			wantErr: true,
		},
		{
			name: "Test Case 3: Failed Login due to Verify Authentication error",
			auth: &MockAuthClient{
				AuthenticationChallengeResponse: &interactor.AuthenticationChallengeResponse{AuthId: "authId", C: 0},
				AuthenticationAnswerError:       errors.New("verify authentication error"),
			},
			co: &MockCommitmentExecuter{
				Result: &app.CommitmentResult{R1: big.NewInt(1), R2: big.NewInt(1), K: big.NewInt(1)},
			},
			cs: &MockComputeSExecuter{
				Result: big.NewInt(1),
			},
			cfg: &config.Config{},
			args: args{
				ctx:          context.Background(),
				userName:     "test",
				userPassword: big.NewInt(12345),
			},
			wantErr: true,
		},
		{
			name: "Test Case 4: Failed Login due to Commitment Exec error",
			auth: &MockAuthClient{
				AuthenticationChallengeResponse: &interactor.AuthenticationChallengeResponse{AuthId: "authId", C: 0},
			},
			co: &MockCommitmentExecuter{
				Err: errors.New("commitment exec error"),
			},
			cs: &MockComputeSExecuter{
				Result: big.NewInt(1),
			},
			cfg: &config.Config{},
			args: args{
				ctx:          context.Background(),
				userName:     "test",
				userPassword: big.NewInt(12345),
			},
			wantErr: true,
		},
		{
			name: "Test Case 5: Failed Login due to ComputeS Exec error",
			auth: &MockAuthClient{
				AuthenticationChallengeResponse: &interactor.AuthenticationChallengeResponse{AuthId: "authId", C: 0},
			},
			co: &MockCommitmentExecuter{
				Result: &app.CommitmentResult{R1: big.NewInt(1), R2: big.NewInt(1), K: big.NewInt(1)},
			},
			cs: &MockComputeSExecuter{
				Err: errors.New("computeS exec error"),
			},
			cfg: &config.Config{},
			args: args{
				ctx:          context.Background(),
				userName:     "test",
				userPassword: big.NewInt(12345),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			address := ":50051"

			// Create client with mocks
			c, err := NewClient(address, tt.cfg, tt.re, tt.co, tt.cs)
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}

			// Replace auth client with a mock
			c.auth = tt.auth

			_, err = c.Login(tt.args.ctx, tt.args.userName, tt.args.userPassword)
			if (err != nil) != tt.wantErr {
				t.Errorf("Login() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestAuthenticationClient_Register(t *testing.T) {
	type args struct {
		ctx          context.Context
		userName     string
		userPassword *big.Int
	}
	tests := []struct {
		name    string
		args    args
		cfg     *config.Config
		auth    *MockAuthClient
		re      *MockRegisterExecuter
		co      *MockCommitmentExecuter
		cs      *MockComputeSExecuter
		wantErr bool
	}{
		{
			name: "Test Case 1: Successful Register",
			auth: &MockAuthClient{
				RegisterResponse:                &interactor.RegisterResponse{},
				AuthenticationChallengeResponse: &interactor.AuthenticationChallengeResponse{AuthId: "authId", C: 0},
				AuthenticationAnswerResponse:    &interactor.AuthenticationAnswerResponse{SessionId: "sessionId"},
			},
			re: &MockRegisterExecuter{
				Y1: big.NewInt(1),
				Y2: big.NewInt(1),
			},
			co: &MockCommitmentExecuter{
				Result: &app.CommitmentResult{R1: big.NewInt(1), R2: big.NewInt(1), K: big.NewInt(1)},
			},
			cs: &MockComputeSExecuter{
				Result: big.NewInt(1),
			},
			cfg: &config.Config{},
			args: args{
				ctx:          context.Background(),
				userName:     "test",
				userPassword: big.NewInt(12345),
			},
			wantErr: false,
		},
		{
			name: "Test Case 2: Failed Register",
			auth: &MockAuthClient{
				RegisterError: errors.New("register error"),
			},
			re: &MockRegisterExecuter{
				Y1:      big.NewInt(1),
				Y2:      big.NewInt(1),
				ExecErr: errors.New("exec error"),
			},
			co: &MockCommitmentExecuter{
				Err: errors.New("exec error"),
			},
			cs: &MockComputeSExecuter{
				Result: big.NewInt(1),
				Err:    errors.New("exec error"),
			},
			cfg: &config.Config{},
			args: args{
				ctx:          context.Background(),
				userName:     "test",
				userPassword: big.NewInt(12345),
			},
			wantErr: true,
		},
		{
			name: "Test Case 3: Failed Register due to c.auth.Register error",
			auth: &MockAuthClient{
				RegisterError: errors.New("register error"),
			},
			re: &MockRegisterExecuter{
				Y1: big.NewInt(1),
				Y2: big.NewInt(1),
			},
			co: &MockCommitmentExecuter{
				Result: &app.CommitmentResult{R1: big.NewInt(1), R2: big.NewInt(1), K: big.NewInt(1)},
			},
			cs: &MockComputeSExecuter{
				Result: big.NewInt(1),
			},
			cfg: &config.Config{},
			args: args{
				ctx:          context.Background(),
				userName:     "test",
				userPassword: big.NewInt(12345),
			},
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			address := ":50051"

			// Create client with mocks
			c, err := NewClient(address, tt.cfg, tt.re, tt.co, tt.cs)
			if err != nil {
				t.Fatalf("failed to create client: %v", err)
			}

			// Replace auth client with a mock
			c.auth = tt.auth

			if err = c.Register(tt.args.ctx, tt.args.userName, tt.args.userPassword); (err != nil) != tt.wantErr {
				t.Errorf("Register() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
