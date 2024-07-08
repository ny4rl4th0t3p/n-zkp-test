package app

import (
	"context"
	"errors"
	"math/big"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"

	"practical-case-test/config"
	"practical-case-test/internal/domain/auth"
	interactor "practical-case-test/internal/interactor/proto"
)

func TestVerifyAuthentication_Exec(t *testing.T) {
	cfg := &config.Config{
		G: big.NewInt(3),
		H: big.NewInt(5),
		Q: big.NewInt(13),
	}

	uID := "UserID1"
	authID := "AuthId1"

	var y1, y2 int64 = 6, 8
	var r1, r2 int64 = 7, 5
	var c = big.NewInt(11)
	var s int64 = 4

	user, _ := auth.NewUser(uID, y1, y2)
	challenge, _ := auth.NewChallenge(c, uID, r1, r2, 123456789)

	req := &interactor.AuthenticationAnswerRequest{
		AuthId: authID,
		S:      s,
	}

	testCases := []struct {
		name    string
		request *interactor.AuthenticationAnswerRequest
		setup   func(ar *mockAuthRepository)
		check   func(*auth.Session, error)
	}{
		{
			name:    "Successful Path",
			request: req,
			setup: func(ar *mockAuthRepository) {
				ar.On("GetUserRegistration", context.Background(), uID).Return(user, nil)
				ar.On("GetAuthenticationChallenge", context.Background(), authID).Return(challenge, nil)
				ar.On("StoreSession", context.Background(), mock.Anything).Return(nil)
			},
			check: func(s *auth.Session, err error) {
				require.NoError(t, err)
				require.NotNil(t, s)
			},
		},
		{
			name:    "GetAuthenticationChallenge fails",
			request: req,
			setup: func(ar *mockAuthRepository) {
				ar.On("GetAuthenticationChallenge", context.Background(), authID).Return(nil, errors.New("GetAuthenticationChallenge error"))
			},
			check: func(_ *auth.Session, err error) {
				require.Error(t, err)
			},
		},
		{
			name:    "GetUserRegistration fails",
			request: req,
			setup: func(ar *mockAuthRepository) {
				ar.On("GetAuthenticationChallenge", context.Background(), authID).Return(challenge, nil)
				ar.On("GetUserRegistration", context.Background(), challenge.UserID()).Return(nil, errors.New("GetUserRegistration error"))
			},
			check: func(_ *auth.Session, err error) {
				require.Error(t, err)
			},
		},
		{
			name:    "StoreSession fails",
			request: req,
			setup: func(ar *mockAuthRepository) {
				ar.On("GetUserRegistration", context.Background(), uID).Return(user, nil)
				ar.On("GetAuthenticationChallenge", context.Background(), authID).Return(challenge, nil)
				ar.On("StoreSession", context.Background(), mock.Anything).Return(errors.New("Store Session Error"))
			},
			check: func(_ *auth.Session, err error) {
				require.Error(t, err)
			},
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			ar := new(mockAuthRepository)
			va := NewVerifyAuthentication(ar)
			tt.setup(ar)
			sess, err := va.Exec(context.Background(), cfg, tt.request)
			tt.check(sess, err)
			ar.AssertExpectations(t)
		})
	}
}
