package auth

import (
	"errors"
	"math/big"

	"github.com/google/uuid"
)

var (
	ErrInvalidChallenge = errors.New("invalid challenge")
)

type Challenge struct {
	userID    string
	authID    uuid.UUID
	c         *big.Int
	r1        int64
	r2        int64
	timestamp int64
}

func (c Challenge) R1() int64 {
	return c.r1
}

func (c Challenge) R2() int64 {
	return c.r2
}

func (c Challenge) UserID() string {
	return c.userID
}

func (c Challenge) AuthID() uuid.UUID {
	return c.authID
}

func (c Challenge) C() *big.Int {
	return c.c
}

func (c Challenge) Timestamp() int64 {
	return c.timestamp
}

func NewChallenge(c *big.Int, userID string, r1, r2, timestamp int64) (*Challenge, error) {
	authID := uuid.New()
	ch := &Challenge{
		userID:    userID,
		authID:    authID,
		c:         c,
		r1:        r1,
		r2:        r2,
		timestamp: timestamp,
	}
	if !ch.IsValid() {
		return nil, ErrInvalidChallenge
	}
	return ch, nil
}

func (c Challenge) IsValid() bool {
	return !(c.userID == "" || c.authID == uuid.Nil || c.c == nil || c.c.Cmp(big.NewInt(0)) == 0)
}
