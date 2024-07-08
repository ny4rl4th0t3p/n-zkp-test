package app

import (
	"log/slog"
	"math/big"

	"practical-case-test/config"
	interactor "practical-case-test/internal/interactor/proto"
)

// ComputeSExecuter is an interface that defines the `Exec` method for executing the computational logic of `ComputeS` operation.
//
// The `Exec` method takes a `cfg` configuration object, `x` and `k` big integers, and a pointer to an `AuthenticationChallengeResponse` object.
// It returns a big integer and an error.
type ComputeSExecuter interface {
	Exec(cfg *config.Config, x, k *big.Int, res *interactor.AuthenticationChallengeResponse) (*big.Int, error)
}

// ComputeS represents a type that computes the value of S based on the provided inputs.
type ComputeS struct{}

// NewComputeS returns a new instance of ComputeSExecuter.
func NewComputeS() ComputeSExecuter {
	return &ComputeS{}
}

// Exec calculates the value of s by using the given configuration, x, k, and res parameters.
// It calculates s using the formula: s = (k - (c * x)) mod q, where c is obtained from res.GetC().
// The function returns the calculated value of s and an error, if any.
// If the configuration is nil, it returns nil and an error indicating that the config cannot be nil.
// If the value of q in the configuration is zero, it returns nil and an error indicating that q cannot be zero.
// The function logs the received value of c, k, and res before invoking the calculateS function to calculate s.
func (ru ComputeS) Exec(cfg *config.Config, x, k *big.Int, res *interactor.AuthenticationChallengeResponse) (
	*big.Int,
	error,
) {
	c := new(big.Int).SetInt64(res.GetC())

	slog.Info("received c", "c", c, "k", k, "res", res)

	s, err := calculateS(cfg, c, x, k)
	if err != nil {
		return nil, err
	}

	return s, nil
}
