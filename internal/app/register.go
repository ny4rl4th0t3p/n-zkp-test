package app

import (
	"log/slog"
	"math/big"

	"practical-case-test/config"
)

// RegisterExecuter is an interface that defines the `Exec` method for executing the registration process.
// The `Exec` method takes a `cfg` configuration object and a `userPassword` of type `*big.Int`.
// It returns two big integers of type `*big.Int` and an error.
type RegisterExecuter interface {
	Exec(cfg *config.Config, userPassword *big.Int) (*big.Int, *big.Int, error)
}

// Register represents a type used for registration functionality.
type Register struct{}

// NewRegister returns a new instance of RegisterExecuter.
func NewRegister() RegisterExecuter {
	return &Register{}
}

// Exec executes the registration process by calculating y1 and y2 values based on the provided
// config and user password. The y1 and y2 values are returned along with any error that occurred.
// The function logs the registration values before returning.
func (ru Register) Exec(cfg *config.Config, userPassword *big.Int) (
	*big.Int,
	*big.Int,
	error,
) {
	y1, y2, err := calculateYs(cfg, userPassword)
	if err != nil {
		return nil, nil, err
	}

	slog.Info("registering with values", "y1", y1, "y2", y2)

	return y1, y2, nil
}
