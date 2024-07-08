package memory

import (
	"context"
	"errors"
	"sync"

	authDomain "practical-case-test/internal/domain/auth"

	"github.com/google/uuid"
)

var (
	ErrCastUser          = errors.New("error casting user")
	ErrCastChallenge     = errors.New("error casting challenge")
	ErrCastSession       = errors.New("error casting session")
	ErrUserIDNotFound    = errors.New("userID not found")
	ErrAuthIDNotFound    = errors.New("AuthID not found")
	ErrSessionNotFound   = errors.New("session not found")
	ErrUserAlreadyExists = errors.New("user already exists")
)

type InMemAuthRepository struct {
	userRegistration sync.Map
	authChallenge    sync.Map
	sessions         sync.Map
}

func NewInMemAuthRepository() *InMemAuthRepository {
	return &InMemAuthRepository{
		userRegistration: sync.Map{},
		authChallenge:    sync.Map{},
	}
}

// StoreUserRegistration stores the user registration in the InMemAuthRepository.
// It first checks if the context has an error and returns the error if it exists.
// Then it validates the user object and returns authDomain.ErrInvalidUser if it is invalid.
// Next, it checks if the user already exists in the repository and returns ErrUserAlreadyExists if it does.
// Finally, it stores the user registration in the repository using the user's UserID as the key.
// It returns nil if the user registration is successfully stored.
func (repo *InMemAuthRepository) StoreUserRegistration(ctx context.Context, user authDomain.User) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if !user.IsValid() {
		return authDomain.ErrInvalidUser
	}
	if _, loaded := repo.userRegistration.Load(user.UserID()); loaded {
		return ErrUserAlreadyExists
	}
	repo.userRegistration.Store(user.UserID(), user)
	return nil
}

// GetUserRegistration retrieves the user registration based on the provided userID.
// It returns the user registration if found, otherwise it returns an error.
func (repo *InMemAuthRepository) GetUserRegistration(ctx context.Context, userID string) (*authDomain.User, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	val, loadOk := repo.userRegistration.Load(userID)
	if !loadOk {
		return nil, ErrUserIDNotFound
	}
	user, castOk := val.(authDomain.User)
	if !castOk {
		return nil, ErrCastUser
	}
	return &user, nil
}

// StoreAuthenticationChallenge stores an authentication challenge in the in-memory repository.
// It validates the challenge and returns an error if the challenge is invalid.
// If the context is canceled or has timed out, it returns an error.
// The challenge is stored in the authChallenge map using its AuthID as the key.
// If a challenge with the same AuthID already exists, it is overwritten.
// The method returns nil if the challenge is stored successfully.
func (repo *InMemAuthRepository) StoreAuthenticationChallenge(ctx context.Context, challenge authDomain.Challenge) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if !challenge.IsValid() {
		return authDomain.ErrInvalidChallenge
	}
	repo.authChallenge.Store(challenge.AuthID().String(), challenge)
	return nil
}

// GetAuthenticationChallenge retrieves the authentication challenge based on the provided authID.
// It returns the challenge and nil error if the challenge is found, otherwise it returns nil challenge
// and the ErrAuthIDNotFound error. If the challenge cannot be casted to authDomain.Challenge,
// it returns nil challenge and the ErrCastChallenge error.
func (repo *InMemAuthRepository) GetAuthenticationChallenge(ctx context.Context, authID string) (*authDomain.Challenge, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	val, loadOk := repo.authChallenge.Load(authID)
	if !loadOk {
		return nil, ErrAuthIDNotFound
	}
	challenge, castOk := val.(authDomain.Challenge)
	if !castOk {
		return nil, ErrCastChallenge
	}
	return &challenge, nil
}

// StoreSession stores the given session in the in-memory repository.
// It first checks if the context has an error, and returns the error if present.
// Then it checks if the session is valid using the IsValid method of the session.
// If the session is not valid, it returns the error ErrInvalidSession from the authDomain package.
// It then generates a session key using the generateSessionKey function, passing the user ID and session ID.
// If an error occurs during the generation of the session key, it is returned.
// Finally, it stores the session in the sessions map of the repository using the generated session key.
// It returns nil if the session is stored successfully.
func (repo *InMemAuthRepository) StoreSession(ctx context.Context, session authDomain.Session) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	if !session.IsValid() {
		return authDomain.ErrInvalidSession
	}
	sessionKey, err := generateSessionKey(session.UserID(), session.ID().String())
	if err != nil {
		return err
	}
	repo.sessions.Store(sessionKey, session)
	return nil
}

// GetSession retrieves a session from the InMemAuthRepository based on the provided userID and sessionID.
// It returns the session if found and valid, otherwise returns an error.
// The session key is generated using the userID and sessionID strings.
// The session key is used to load the session from the sessions sync.Map.
// If the session is not found, ErrSessionNotFound is returned.
// If the loaded value is not of type *authDomain.Session, ErrCastSession is returned.
// If the loaded session is not valid, authDomain.ErrInvalidSession is returned.
// Finally, the method returns the session if it is found and valid, or an error otherwise.
func (repo *InMemAuthRepository) GetSession(ctx context.Context, userID string, sessionID uuid.UUID) (*authDomain.Session, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	sessionKey, err := generateSessionKey(userID, sessionID.String())
	if err != nil {
		return nil, err
	}
	val, loadOk := repo.sessions.Load(sessionKey)
	if !loadOk {
		return nil, ErrSessionNotFound
	}
	session, castOk := val.(*authDomain.Session)
	if !castOk {
		return nil, ErrCastSession
	}
	if !session.IsValid() {
		return nil, authDomain.ErrInvalidSession
	}
	return session, nil
}

// generateSessionKey takes a userID and sessionID as input and generates a session key
// by concatenating userID and sessionID with a colon ":" delimiter. It returns the generated
// session key and any error that occurred during the process. If the userID or sessionID is empty,
// an error is returned with a corresponding error message.
func generateSessionKey(userID, sessionID string) (string, error) {
	if userID == "" {
		return "", errors.New("userID should not be empty")
	}
	if sessionID == "" {
		return "", errors.New("sessionID should not be empty")
	}
	return userID + ":" + sessionID, nil
}
