package app

import (
	"context"
	"errors"
	"testing"

	data "practical-case-test/internal/interactor/proto"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

// mockAuthRepository was defined in the previous message

func TestRegisterUser_Exec(t *testing.T) {
	testCases := []struct {
		name          string
		req           *data.RegisterRequest
		mockStoreErr  error
		expectedError string
	}{
		{
			name:         "Register user successful case",
			req:          &data.RegisterRequest{User: "testUser", Y1: 0, Y2: 0},
			mockStoreErr: nil,
		},
		{
			name:          "StoreUserRegistration returns error",
			req:           &data.RegisterRequest{User: "testUser", Y1: 0, Y2: 0},
			mockStoreErr:  errors.New("store user registration error"),
			expectedError: "store user registration error",
		},
		{
			name:          "Invalid user - Empty username",
			req:           &data.RegisterRequest{User: "", Y1: 0, Y2: 0},
			expectedError: "invalid user",
		},
		{
			name:          "Invalid user - Negative Y1",
			req:           &data.RegisterRequest{User: "testUser", Y1: -1, Y2: 0},
			expectedError: "invalid user",
		},
		{
			name:          "Invalid user - Negative Y2",
			req:           &data.RegisterRequest{User: "testUser", Y1: 0, Y2: -1},
			expectedError: "invalid user",
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ar := new(mockAuthRepository)
			registrar := NewRegisterUser(ar)

			if tt.expectedError != "invalid user" {
				ar.On("StoreUserRegistration", mock.Anything, mock.Anything).Return(tt.mockStoreErr)
			}

			err := registrar.Exec(context.Background(), tt.req)

			if tt.expectedError != "" {
				require.Error(t, err, "Expected an error but got none")
				require.Equal(t, tt.expectedError, err.Error(), "Expected error of type %v, but got type %v", tt.expectedError, err.Error())
			} else {
				require.NoError(t, err, "Did not expect an error but got one")
			}
			ar.AssertExpectations(t)
		})
	}
}
