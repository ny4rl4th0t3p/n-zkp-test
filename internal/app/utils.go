package app

import (
	"crypto/rand"
	"errors"
	"log/slog"
	"math"
	"math/big"
	mrand "math/rand"
	"sync"
	"time"

	"practical-case-test/config"
	"practical-case-test/internal/domain/auth"
)

// ErrZeroQ is an error indicating that the value of q should not be 0.
// ErrNilConfig is an error indicating that the config should not be nil.
var (
	ErrZeroQ     = errors.New("q cannot be 0")
	ErrNilConfig = errors.New("config cannot be nil")
)

// charset is a constant string that contains all the lowercase and uppercase alphabets.
const charset = "abcdefghijklmnopqrstuvwxyz" +
	"ABCDEFGHIJKLMNOPQRSTUVWXYZ"

// RandString generates a random string of the specified length using the characters defined in 'charset'
func RandString(length int) string {
	r := mrand.New(mrand.NewSource(time.Now().UnixNano()))
	b := make([]byte, length)
	for i := range b {
		b[i] = charset[r.Intn(len(charset))]
	}
	return string(b)
}

// RandomPassword generates a random password as a big.Int using the math.MaxInt16 as the upper bound.
func RandomPassword() (*big.Int, error) {
	return rand.Int(rand.Reader, big.NewInt(math.MaxInt16))
}

// calculateYs calculates y1 and y2 values based on the provided config and user password.
// It returns the calculated y1 and y2 values along with any error that occurred.
// If the given config is nil, it returns nil for both y1 and y2 and an error with the message "config cannot be nil".
// If the length of config's Q value is 0, it returns nil for both y1 and y2 and an error with the message "q cannot be 0".
// Otherwise, it calculates y1 as cfg.G^userPassword mod cfg.Q and y2 as cfg.H^userPassword mod cfg.Q.
// The function does not perform any logging or additional operations beyond the calculations.
// It is important to note that the function assumes the correctness of the provided input parameters.
func calculateYs(cfg *config.Config, userPassword *big.Int) (y1, y2 *big.Int, err error) {
	if cfg == nil {
		return nil, nil, ErrNilConfig
	}
	if len(cfg.Q.Bits()) == 0 {
		return nil, nil, ErrZeroQ
	}
	y1Tmp := new(big.Int).Exp(cfg.G, userPassword, nil)
	y1 = new(big.Int).Mod(y1Tmp, cfg.Q)

	y2Tmp := new(big.Int).Exp(cfg.H, userPassword, nil)
	y2 = new(big.Int).Mod(y2Tmp, cfg.Q)

	return
}

// calculateCommitment calculates the commitment values (r1, r2, and k) based on the provided config.
// It generates a random number k using math.MaxInt16 as the maximum value.
// It then adds 1 to k and logs the generated value.
// It calculates r1 and r2 by exponentiating the values cfg.G and cfg.H to the power of k, respectively.
// The results r1 and r2 are then computed modulo cfg.Q.
// The function runs the calculations concurrently using goroutines and waits for them to finish using a WaitGroup.
// If an error occurs during the generation of k, it returns nil for all values and the error.
// Otherwise, it returns the calculated values r1, r2, and k, along with nil error.
func calculateCommitment(cfg *config.Config) (r1, r2, k *big.Int, err error) {
	k, err = rand.Int(rand.Reader, big.NewInt(math.MaxInt16))
	if err != nil {
		return nil, nil, nil, err
	}
	k.Add(k, big.NewInt(1))

	slog.Info("generating commitment", "k", k)

	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		r1Tmp := new(big.Int).Exp(cfg.G, k, nil)
		r1 = new(big.Int).Mod(r1Tmp, cfg.Q)
	}()

	go func() {
		defer wg.Done()
		r2Tmp := new(big.Int).Exp(cfg.H, k, nil)
		r2 = new(big.Int).Mod(r2Tmp, cfg.Q)
	}()

	wg.Wait()

	return
}

// calculateS calculates the value of s based on the given configuration, c, x, and k.
// It calculates s as (k - (c * x)) % q, where q is obtained from cfg.Q.
// The function returns the calculated value of s and an error, if any.
// If the configuration is nil, it returns nil and an error indicating that the config cannot be nil.
// If the value of q in the configuration is zero, it returns nil and an error indicating that q cannot be zero.
func calculateS(cfg *config.Config, c, x, k *big.Int) (*big.Int, error) {
	if cfg == nil {
		return nil, ErrNilConfig
	}
	if len(cfg.Q.Bits()) == 0 {
		return nil, ErrZeroQ
	}
	cx := new(big.Int).Mul(c, x)
	kSubCx := new(big.Int).Sub(k, cx)
	s := new(big.Int).Mod(kSubCx, cfg.Q)

	return s, nil
}

// verifyS verifies the authenticity of the user's response to the authentication challenge
// by performing calculations and comparing the results with the challenge's expected values.
func verifyS(cfg *config.Config, challenge *auth.Challenge, user *auth.User, s int64) bool {
	cBigInt := challenge.C()
	y1BigInt := big.NewInt(user.Y1())
	y2BigInt := big.NewInt(user.Y2())
	sBigInt := big.NewInt(s)

	gBigInt := cfg.G
	hBigInt := cfg.H

	leftR1 := new(big.Int).Exp(gBigInt, sBigInt, nil)
	rightR1 := new(big.Int).Exp(y1BigInt, cBigInt, nil)
	r1Tmp := new(big.Int).Mul(leftR1, rightR1)
	r1 := new(big.Int).Mod(r1Tmp, cfg.Q)

	leftR2 := new(big.Int).Exp(hBigInt, sBigInt, nil)
	rightR2 := new(big.Int).Exp(y2BigInt, cBigInt, nil)
	r2Tmp := new(big.Int).Mul(leftR2, rightR2)
	r2 := new(big.Int).Mod(r2Tmp, cfg.Q)

	slog.Info("r1", "r1 calculated locally", r1, "r1 received", challenge.R1())
	slog.Info("r2", "r2 calculated locally", r2, "r2 received", challenge.R2())

	if r1.Cmp(big.NewInt(challenge.R1())) != 0 || r2.Cmp(big.NewInt(challenge.R2())) != 0 {
		return false
	}

	return true
}
