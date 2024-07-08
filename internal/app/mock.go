package app

import (
	"context"

	"practical-case-test/internal/domain/auth"

	"github.com/google/uuid"
	"github.com/stretchr/testify/mock"
)

type mockAuthRepository struct {
	mock.Mock
}

func (m *mockAuthRepository) StoreUserRegistration(ctx context.Context, user auth.User) error {
	args := m.Called(ctx, user)
	return args.Error(0)
}

func (m *mockAuthRepository) GetAuthenticationChallenge(ctx context.Context, authID string) (*auth.Challenge, error) {
	args := m.Called(ctx, authID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.Challenge), args.Error(1)
}

func (m *mockAuthRepository) StoreSession(ctx context.Context, session auth.Session) error {
	args := m.Called(ctx, session)
	return args.Error(0)
}

func (m *mockAuthRepository) GetSession(_ context.Context, _ string, _ uuid.UUID) (*auth.Session, error) {
	panic("implement me")
}

func (m *mockAuthRepository) GetUserRegistration(ctx context.Context, userID string) (*auth.User, error) {
	args := m.Called(ctx, userID)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*auth.User), args.Error(1)
}

func (m *mockAuthRepository) StoreAuthenticationChallenge(ctx context.Context, challenge auth.Challenge) error {
	args := m.Called(ctx, challenge)
	return args.Error(0)
}
