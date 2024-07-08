package auth

import (
	"errors"
)

var (
	ErrInvalidUser = errors.New("invalid user")
)

type User struct {
	userID string
	y1     int64
	y2     int64
}

func (u User) UserID() string {
	return u.userID
}

func (u User) Y1() int64 {
	return u.y1
}

func (u User) Y2() int64 {
	return u.y2
}

func NewUser(user string, y1 int64, y2 int64) (*User, error) {
	u := &User{userID: user, y1: y1, y2: y2}
	if !u.IsValid() {
		return nil, ErrInvalidUser
	}
	return u, nil
}

func (u User) IsValid() bool {
	if u.userID == "" || u.y1 < 0 || u.y2 < 0 {
		return false
	}
	return true
}
