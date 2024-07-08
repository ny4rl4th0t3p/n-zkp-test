package memory

import (
	"context"
	"math/big"
	"sync"
	"testing"
	"time"

	authDomain "practical-case-test/internal/domain/auth"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestInMemAuthRepository_GetAuthenticationChallenge(t *testing.T) {
	repo := &InMemAuthRepository{
		userRegistration: sync.Map{},
		authChallenge:    sync.Map{},
		sessions:         sync.Map{},
	}

	challenge, err := authDomain.NewChallenge(big.NewInt(2), "user-id-1", 0, 2, time.Now().Unix())
	require.NoError(t, err)

	err = repo.StoreAuthenticationChallenge(context.Background(), *challenge)
	require.NoError(t, err)

	existingAuthID := challenge.AuthID().String()

	type args struct {
		ctx    context.Context
		authID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nonexistent auth challenge",
			args: args{
				ctx:    context.Background(),
				authID: "nonexistentAuthId",
			},
			wantErr: true,
		},
		{
			name: "existing auth challenge",
			args: args{
				ctx:    context.Background(),
				authID: existingAuthID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var got *authDomain.Challenge
			got, err = repo.GetAuthenticationChallenge(tt.args.ctx, tt.args.authID)
			if tt.wantErr {
				require.Error(t, err, "GetAuthenticationChallenge() should return an error")
				return
			}
			require.Equal(t, tt.args.authID, got.AuthID().String(), "GetAuthenticationChallenge() got.AuthID().String() = %v, "+
				"tt.args.authId %v", got.AuthID().String(), tt.args.authID)
		})
	}
}

func TestInMemAuthRepository_GetSession(t *testing.T) {
	repo := &InMemAuthRepository{
		userRegistration: sync.Map{},
		authChallenge:    sync.Map{},
		sessions:         sync.Map{},
	}

	sessionID, _ := uuid.Parse("550e8400-e29b-41d4-a716-446655440000")
	testSession, err := authDomain.NewSession(sessionID, "existingUser", time.Now().Unix())
	require.NoError(t, err)

	sessionKey, err := generateSessionKey("existingUser", sessionID.String())
	require.NoError(t, err)

	repo.sessions.Store(sessionKey, testSession)

	type args struct {
		ctx       context.Context
		userID    string
		sessionID uuid.UUID
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "nonexistent session",
			args: args{
				ctx:       context.Background(),
				userID:    "nonexistentUser",
				sessionID: uuid.New(), // non existing uuid
			},
			wantErr: true,
		},
		{
			name: "existing session",
			args: args{
				ctx:       context.Background(),
				userID:    "existingUser",
				sessionID: sessionID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var got *authDomain.Session
			got, err = repo.GetSession(tt.args.ctx, tt.args.userID, tt.args.sessionID)
			if tt.wantErr {
				require.Error(t, err, "GetSession() should return an error for non-existing session")
				return
			}
			require.NoError(t, err, "GetSession() should not return an error for existing session")
			require.Equal(t, tt.args.sessionID, got.ID(), "GetSession() fetched session has unexpected ID. got = %v, want = %v", got.ID(), tt.args.sessionID)
			require.Equal(t, tt.args.userID, got.UserID(), "GetSession() fetched session has unexpected UserID. got = %v, want = %v", got.UserID(), tt.args.userID)
		})
	}
}
func TestInMemAuthRepository_GetUserRegistration(t *testing.T) {
	repo := &InMemAuthRepository{
		userRegistration: sync.Map{},
		authChallenge:    sync.Map{},
		sessions:         sync.Map{},
	}

	userID := "existingUser"
	testUser, _ := authDomain.NewUser(userID, 10, 20)
	// repo.userRegistration.Store(userID, testUser)

	err := repo.StoreUserRegistration(context.Background(), *testUser)
	require.NoError(t, err)

	type args struct {
		ctx    context.Context
		userID string
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "non-existing user",
			args: args{
				ctx:    context.Background(),
				userID: "nonexistentUser",
			},
			wantErr: true,
		},
		{
			name: "existing user",
			args: args{
				ctx:    context.Background(),
				userID: userID,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			var got *authDomain.User
			got, err = repo.GetUserRegistration(tt.args.ctx, tt.args.userID)
			if tt.wantErr {
				require.Error(t, err, "GetUserRegistration() should return an error for non-existing user")
				return
			}
			require.NoError(t, err, "GetUserRegistration() should not return an error for an existing user")
			require.Equal(t, tt.args.userID, got.UserID(), "GetUserRegistration() got.ID() = %v, want = %v", got.UserID(), tt.args.userID)
		})
	}
}

func TestInMemAuthRepository_StoreAuthenticationChallenge(t *testing.T) {
	repo := &InMemAuthRepository{
		userRegistration: sync.Map{},
		authChallenge:    sync.Map{},
		sessions:         sync.Map{},
	}

	challenge1, err := authDomain.NewChallenge(big.NewInt(2), "user-id-1", 0, 1, time.Now().Unix())
	require.NoError(t, err)

	challenge2, err := authDomain.NewChallenge(big.NewInt(3), "user-id-2", 1, 5, time.Now().Unix())
	require.NoError(t, err)

	type args struct {
		ctx       context.Context
		challenge authDomain.Challenge
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "storing auth challenge for user 1",
			args: args{
				ctx:       context.Background(),
				challenge: *challenge1,
			},
			wantErr: false,
		},
		{
			name: "storing auth challenge for user 2",
			args: args{
				ctx:       context.Background(),
				challenge: *challenge2,
			},
			wantErr: false,
		},
		{
			name: "storing the same auth challenge again",
			args: args{
				ctx:       context.Background(),
				challenge: *challenge1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err = repo.StoreAuthenticationChallenge(tt.args.ctx, tt.args.challenge)
			if tt.wantErr {
				require.Error(t, err, "StoreAuthenticationChallenge() should return an error")
				return
			}
			require.NoError(t, err, "StoreAuthenticationChallenge() should not return an error")

			storedChallenge, ok := repo.authChallenge.Load(tt.args.challenge.AuthID().String())
			require.True(t, ok, "Stored challenge not found in the authChallenge map")
			require.IsType(t, tt.args.challenge, storedChallenge, "Stored challenge in the map is not of the correct type")

			storedChallengeTyped, ok := storedChallenge.(authDomain.Challenge)
			require.True(t, ok)
			require.Equal(t, tt.args.challenge, storedChallengeTyped, "Stored challenge not equal to the given challenge. got = %v, want = %v", storedChallengeTyped, tt.args.challenge)
		})
	}
}

func TestInMemAuthRepository_StoreSession(t *testing.T) {
	repo := &InMemAuthRepository{
		userRegistration: sync.Map{},
		authChallenge:    sync.Map{},
		sessions:         sync.Map{},
	}

	session1, err := authDomain.NewSession(uuid.New(), "user-id-1", time.Now().Add(5*time.Hour).Unix())
	require.NoError(t, err)

	session2, err := authDomain.NewSession(uuid.New(), "user-id-2", time.Now().Add(2*time.Hour).Unix())
	require.NoError(t, err)

	type args struct {
		ctx     context.Context
		session authDomain.Session
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "storing session for user 1",
			args: args{
				ctx:     context.Background(),
				session: *session1,
			},
			wantErr: false,
		},
		{
			name: "storing session for user 2",
			args: args{
				ctx:     context.Background(),
				session: *session2,
			},
			wantErr: false,
		},
		{
			name: "storing the same session again",
			args: args{
				ctx:     context.Background(),
				session: *session1,
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel() // Run test cases in parallel
			err = repo.StoreSession(tt.args.ctx, tt.args.session)
			if tt.wantErr {
				require.Error(t, err, "StoreSession() should return an error")
				return
			}
			require.NoError(t, err, "StoreSession() should not return an error")
			var sessionKey string
			sessionKey, err = generateSessionKey(tt.args.session.UserID(), tt.args.session.ID().String())
			require.NoError(t, err)
			storedSession, ok := repo.sessions.Load(sessionKey)
			require.True(t, ok, "Stored session not found in sessions map")
			require.IsType(t, tt.args.session, storedSession, "Stored session in map is not of correct type")

			storedSessionTyped, ok := storedSession.(authDomain.Session)
			require.True(t, ok)
			require.Equal(t, tt.args.session, storedSessionTyped, "Stored session not equal to the given session. got = %v, want = %v", storedSessionTyped, tt.args.session)
		})
	}
}

func TestInMemAuthRepository_StoreUserRegistration(t *testing.T) {
	repo := &InMemAuthRepository{
		userRegistration: sync.Map{},
		authChallenge:    sync.Map{},
		sessions:         sync.Map{},
	}

	user1, err := authDomain.NewUser("user-id-1", 1, 2)
	require.NoError(t, err)

	user2, err := authDomain.NewUser("user-id-2", 3, 4)
	require.NoError(t, err)

	type args struct {
		ctx  context.Context
		user *authDomain.User
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "storing registration for user 1",
			args: args{
				ctx:  context.Background(),
				user: user1,
			},
			wantErr: false,
		},
		{
			name: "storing registration for user 2",
			args: args{
				ctx:  context.Background(),
				user: user2,
			},
			wantErr: false,
		},
		{
			name: "storing the same registration again",
			args: args{
				ctx:  context.Background(),
				user: user1,
			},
			wantErr: true, // or false depending on whether your StoreUserRegistration method allows storing the same user again
		},
	}
	for _, tt := range tests {
		// tt := tt // Use local version of tt to avoid race condition.
		t.Run(tt.name, func(t *testing.T) {
			// t.Parallel() // Run test cases in parallel
			err = repo.StoreUserRegistration(tt.args.ctx, *tt.args.user)
			if tt.wantErr {
				require.Error(t, err, "StoreUserRegistration() should return an error")
				return
			}
			require.NoError(t, err, "StoreUserRegistration() should not return an error")

			storedUser, ok := repo.userRegistration.Load(tt.args.user.UserID())
			require.True(t, ok, "Stored user not found in userRegistration map")
			require.IsType(t, *tt.args.user, storedUser, "Stored user in map is not of correct type")

			storedUserTyped, ok := storedUser.(authDomain.User)
			require.True(t, ok)
			require.Equal(t, tt.args.user, &storedUserTyped, "Stored user not equal to the given user. got = %v, want = %v", storedUserTyped, tt.args.user)
		})
	}
}

func TestNewInMemAuthRepository(t *testing.T) {
	tests := []struct {
		name string
		want *InMemAuthRepository
	}{
		{
			name: "new in-memory auth repository initialization",
			want: &InMemAuthRepository{
				userRegistration: sync.Map{},
				authChallenge:    sync.Map{},
				sessions:         sync.Map{},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := NewInMemAuthRepository()

			// Test if userRegistration is empty
			got.userRegistration.Range(func(_, _ interface{}) bool {
				t.Error("userRegistration map should be empty")
				return false
			})

			// Test if authChallenge is empty
			got.authChallenge.Range(func(_, _ interface{}) bool {
				t.Error("authChallenge map should be empty")
				return false
			})

			// Test if sessions is empty
			got.sessions.Range(func(_, _ interface{}) bool {
				t.Error("sessions map should be empty")
				return false
			})
		})
	}
}

func Test_generateSessionKey(t *testing.T) {
	tests := []struct {
		name      string
		userID    string
		sessionID string
		want      string
		wantErr   bool
	}{
		{
			name:      "both userID and sessionID are not empty",
			userID:    "user1",
			sessionID: "session1",
			want:      "user1:session1",
			wantErr:   false,
		},
		{
			name:      "userID is empty but sessionID is not",
			userID:    "",
			sessionID: "session2",
			want:      "",
			wantErr:   true,
		},
		{
			name:      "userID is not empty but sessionID is empty",
			userID:    "user3",
			sessionID: "",
			want:      "",
			wantErr:   true,
		},
		{
			name:      "both userID and sessionID are empty",
			userID:    "",
			sessionID: "",
			want:      "",
			wantErr:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := generateSessionKey(tt.userID, tt.sessionID)
			if (err != nil) != tt.wantErr {
				t.Errorf("generateSessionKey() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("generateSessionKey() = %v, want %v", got, tt.want)
			}
		})
	}
}
