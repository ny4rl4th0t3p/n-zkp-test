package grpc

import (
	"context"
	"errors"
	"math/big"

	"practical-case-test/config"
	"practical-case-test/internal/app"
	"practical-case-test/internal/domain/auth"
	interactor "practical-case-test/internal/interactor/proto"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"google.golang.org/grpc"
)

type MockCreateAuthenticationChallenge struct {
	mock.Mock
}

func (m *MockCreateAuthenticationChallenge) Exec(context.Context, *interactor.AuthenticationChallengeRequest) (
	*auth.Challenge, error) {
	args := m.Called()
	return args.Get(0).(*auth.Challenge), args.Error(1)
}

type MockRegisterUser struct {
	mock.Mock
}

func (m *MockRegisterUser) Exec(ctx context.Context, req *interactor.RegisterRequest) error {
	args := m.Called(ctx, req)
	return args.Error(0)
}

type MockVerifyAuthExecuterSuccess struct{}

func (m *MockVerifyAuthExecuterSuccess) Exec(_ context.Context, _ *config.Config,
	_ *interactor.AuthenticationAnswerRequest) (*auth.Session, error) {
	session, _ := auth.NewSession(uuid.New(), "userId", 1234)
	return session, nil
}

type MockVerifyAuthExecuterFail struct{}

func (m *MockVerifyAuthExecuterFail) Exec(_ context.Context, _ *config.Config, _ *interactor.AuthenticationAnswerRequest) (*auth.Session, error) {
	return nil, errors.New("authentication failed")
}

type MockAuthClient struct {
	RegisterResponse                *interactor.RegisterResponse
	RegisterError                   error
	AuthenticationChallengeResponse *interactor.AuthenticationChallengeResponse
	AuthenticationChallengeError    error
	AuthenticationAnswerResponse    *interactor.AuthenticationAnswerResponse
	AuthenticationAnswerError       error
}

func (m *MockAuthClient) Register(_ context.Context, _ *interactor.RegisterRequest, _ ...grpc.CallOption) (*interactor.RegisterResponse,
	error) {
	return m.RegisterResponse, m.RegisterError
}

func (m *MockAuthClient) CreateAuthenticationChallenge(_ context.Context, _ *interactor.AuthenticationChallengeRequest,
	_ ...grpc.CallOption) (*interactor.AuthenticationChallengeResponse, error) {
	return m.AuthenticationChallengeResponse, m.AuthenticationChallengeError
}

func (m *MockAuthClient) VerifyAuthentication(_ context.Context, _ *interactor.AuthenticationAnswerRequest,
	_ ...grpc.CallOption) (*interactor.AuthenticationAnswerResponse, error) {
	return m.AuthenticationAnswerResponse, m.AuthenticationAnswerError
}

type MockRegisterExecuter struct {
	Y1      *big.Int
	Y2      *big.Int
	ExecErr error
}

func (m *MockRegisterExecuter) Exec(_ *config.Config, _ *big.Int) (*big.Int, *big.Int, error) {
	return m.Y1, m.Y2, m.ExecErr
}

type MockCommitmentExecuter struct {
	Result *app.CommitmentResult
	Err    error
}

func (m *MockCommitmentExecuter) Exec(_ *config.Config) (*app.CommitmentResult, error) {
	return m.Result, m.Err
}

type MockComputeSExecuter struct {
	Result *big.Int
	Err    error
}

func (m *MockComputeSExecuter) Exec(_ *config.Config, _, _ *big.Int, _ *interactor.AuthenticationChallengeResponse) (*big.Int, error) {
	return m.Result, m.Err
}
