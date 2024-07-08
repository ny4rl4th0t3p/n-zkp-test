package app

import (
	"errors"
	"log/slog"
	"math/big"

	"practical-case-test/config"
)

// ErrConfigNil represents an error indicating that the configuration is nil.
var (
	ErrConfigNil = errors.New("config is nil")
)

// CommitmentResult represents the result of a commitment execution.
type CommitmentResult struct {
	R1, R2, K *big.Int
}

// CommitmentExecuter is an interface that defines the `Exec` method for executing the
// commitment step, which generates a CommitmentResult.
//
// The `Exec` method takes a `cfg` configuration object and returns a CommitmentResult
// and an error.
type CommitmentExecuter interface {
	Exec(cfg *config.Config) (*CommitmentResult, error)
}

// Commitment represents a commitment object used for generating random commitment values (r1, r2, and k)
// based on the provided configuration. The `Exec` method is used to generate the commitment values.
type Commitment struct{}

func NewCommitment() CommitmentExecuter {
	return &Commitment{}
}

// Exec generates random commitment values (r1, r2, and k) based on the provided configuration.
// If the config argument is nil, it returns nil for the result and an error of ErrConfigNil.
// Otherwise, it calls the generateRndCommitment function to calculate r1, r2, and k based on the config.
// If an error occurs during the calculation, it returns nil for the result and the error.
// Otherwise, it returns the calculated values r1, r2, and k as a CommitmentResult pointer, along with nil error.
func (ru Commitment) Exec(cfg *config.Config) (
	*CommitmentResult,
	error,
) {
	if cfg == nil {
		return nil, ErrConfigNil
	}
	cr1, cr2, ck, err := generateRndCommitment(cfg)
	if err != nil {
		return nil, err
	}

	return &CommitmentResult{
		R1: cr1,
		R2: cr2,
		K:  ck,
	}, nil
}

// generateRndCommitment generates random commitment values (r1, r2, and k) based on the provided config.
// If the config argument is nil, it returns nil for all values and an error of ErrConfigNil.
// Otherwise, it calls the calculateCommitment function to calculate r1, r2, and k based on the config.
// If an error occurs during the calculation, it returns nil for all values and the error.
// Otherwise, it returns the calculated values r1, r2, and k, along with nil error.
func generateRndCommitment(cfg *config.Config) (*big.Int, *big.Int, *big.Int, error) {
	if cfg == nil {
		return nil, nil, nil, ErrConfigNil
	}

	slog.Info("generating commitment, cfg", "cfg", cfg)

	r1, r2, k, err := calculateCommitment(cfg)
	if err != nil {
		return nil, nil, nil, err
	}

	slog.Info("generating commitment", "r1", r1, "r2", r2)

	return r1, r2, k, nil
}
