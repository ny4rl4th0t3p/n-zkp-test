package auth

import (
	"errors"

	"github.com/google/uuid"
)

var (
	ErrInvalidSession = errors.New("invalid session")
)

type Session struct {
	id             uuid.UUID
	userID         string
	loginTimestamp int64
}

func (s Session) ID() uuid.UUID {
	return s.id
}

func (s Session) UserID() string {
	return s.userID
}

func (s Session) LoginTimestamp() int64 {
	return s.loginTimestamp
}

func NewSession(id uuid.UUID, userID string, loginTimestamp int64) (*Session, error) {
	s := &Session{
		id:             id,
		userID:         userID,
		loginTimestamp: loginTimestamp,
	}
	if !s.IsValid() {
		return nil, ErrInvalidSession
	}
	return s, nil
}

func (s Session) IsValid() bool {
	return s.id != uuid.Nil && s.userID != ""
}
