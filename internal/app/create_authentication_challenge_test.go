package app

import (
	"context"
	"errors"
	"testing"

	"practical-case-test/internal/domain/auth"
	interactor "practical-case-test/internal/interactor/proto"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestCreateAuthenticationChallenge_Exec(t *testing.T) {
	testCases := []struct {
		name          string
		req           *interactor.AuthenticationChallengeRequest
		mockGetUser   *auth.User
		mockGetError  error
		mockStoreErr  error
		expectedError string
	}{
		{
			name:         "Create auth challenge, successful case",
			req:          &interactor.AuthenticationChallengeRequest{User: "testUser", R1: 0, R2: 0},
			mockGetUser:  &auth.User{}, // valid user
			mockGetError: nil,
			mockStoreErr: nil,
		},
		{
			name:          "GetUserRegistration returns error",
			req:           &interactor.AuthenticationChallengeRequest{User: "testUser", R1: 0, R2: 0},
			mockGetUser:   nil,
			mockGetError:  errors.New("get user error"),
			mockStoreErr:  nil,
			expectedError: "get user error",
		},
		{
			name:          "StoreAuthenticationChallenge returns error",
			req:           &interactor.AuthenticationChallengeRequest{User: "testUser", R1: 0, R2: 0},
			mockGetUser:   &auth.User{}, // valid user
			mockGetError:  nil,
			mockStoreErr:  errors.New("store auth challenge error"),
			expectedError: "store auth challenge error",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ar := new(mockAuthRepository)
			creator := NewCreateAuthenticationChallenge(ar)
			ar.On("GetUserRegistration", mock.Anything, tt.req.GetUser()).Return(tt.mockGetUser, tt.mockGetError)
			if tt.mockGetError == nil {
				ar.On("StoreAuthenticationChallenge", mock.Anything, mock.Anything).Return(tt.mockStoreErr)
			}

			_, err := creator.Exec(context.Background(), tt.req)

			if tt.expectedError != "" {
				require.Error(t, err, "Expected an error but got none")
				require.Equal(t, tt.expectedError, err.Error(), "Expected error of type %v, but got type %v", tt.expectedError, err.Error())
			} else {
				require.NoError(t, err, "Did not expect an error but got %v", err)
			}
			ar.AssertExpectations(t)
		})
	}
}
