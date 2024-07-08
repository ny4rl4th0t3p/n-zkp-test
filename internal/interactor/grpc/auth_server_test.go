package grpc

import (
	"context"
	"errors"
	"math/big"
	"testing"
	"time"

	"practical-case-test/config"
	"practical-case-test/internal/app"
	"practical-case-test/internal/domain/auth"
	interactor "practical-case-test/internal/interactor/proto"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAuthenticationServer_CreateAuthenticationChallenge(t *testing.T) {
	// Create an instance of a request.
	req := &interactor.AuthenticationChallengeRequest{User: "test user"}

	mockChallenge, err := auth.NewChallenge(big.NewInt(1234), req.GetUser(), 1, 3, time.Now().Unix())
	require.NoError(t, err)

	testCases := []struct {
		name           string
		mockExecReturn *auth.Challenge
		mockExecError  error
		wantErr        bool
	}{
		{
			name:           "Successful challenge creation",
			mockExecReturn: mockChallenge,
			mockExecError:  nil,
			wantErr:        false,
		},
		{
			name:           "Failed challenge creation",
			mockExecReturn: nil,
			mockExecError:  errors.New("mocked error"),
			wantErr:        true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			fields := struct {
				cac *MockCreateAuthenticationChallenge
			}{
				cac: &MockCreateAuthenticationChallenge{},
			}
			fields.cac.On("Exec", mock.Anything, mock.Anything).Return(tt.mockExecReturn, tt.mockExecError)
			server := &AuthenticationServer{cac: fields.cac}
			_, err = server.CreateAuthenticationChallenge(context.Background(), req)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}
			fields.cac.AssertCalled(t, "Exec", mock.Anything, mock.Anything)
		})
	}
}

func TestAuthenticationServer_Register(t *testing.T) {
	uID := "userId"
	testCases := []struct {
		name    string
		request *interactor.RegisterRequest
		execErr error
		wantErr bool
	}{
		{
			name:    "Successful registration",
			request: &interactor.RegisterRequest{User: uID},
			execErr: nil,
			wantErr: false,
		},
		{
			name:    "Failed registration",
			request: &interactor.RegisterRequest{User: uID},
			execErr: errors.New("failed to register user"),
			wantErr: true,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockRegisterUser := new(MockRegisterUser)
			mockRegisterUser.On("Exec", context.Background(), tt.request).Return(tt.execErr)

			as := &AuthenticationServer{
				cfg: &config.Config{},
				ru:  mockRegisterUser,
			}

			_, err := as.Register(context.Background(), tt.request)
			if tt.wantErr {
				require.Error(t, err)
			} else {
				require.NoError(t, err)
			}

			mockRegisterUser.AssertExpectations(t)
		})
	}
}

func TestAuthenticationServer_VerifyAuthentication(t *testing.T) {
	tests := []struct {
		name        string
		verifyAuth  app.VerifyAuthenticationExecuter
		request     *interactor.AuthenticationAnswerRequest
		wantResp    *interactor.AuthenticationAnswerResponse
		expectError bool
	}{
		{
			name:        "Successful authentication",
			verifyAuth:  &MockVerifyAuthExecuterSuccess{},
			request:     &interactor.AuthenticationAnswerRequest{AuthId: "validAuthID"},
			expectError: false,
		},
		{
			name:        "Failed authentication",
			verifyAuth:  &MockVerifyAuthExecuterFail{},
			request:     &interactor.AuthenticationAnswerRequest{AuthId: "invalidAuthID"},
			expectError: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			server := NewAuthenticationServer(nil, nil, nil, tt.verifyAuth)
			resp, err := server.VerifyAuthentication(context.TODO(), tt.request)

			if tt.expectError {
				require.Error(t, err, "Expected error but got nil")
				return
			}

			require.NoError(t, err, "Got unexpected error")
			require.NoError(t, uuid.Validate(resp.GetSessionId()), "Expected uuid formed ID, got %v", resp.GetSessionId())
		})
	}
}
