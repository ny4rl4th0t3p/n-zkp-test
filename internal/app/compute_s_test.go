package app

import (
	"math"
	"math/big"
	"testing"

	"github.com/stretchr/testify/require"

	"practical-case-test/config"
	interactor "practical-case-test/internal/interactor/proto"
)

func TestComputeS_Exec(t *testing.T) {
	cfg := config.LoadConfig()
	faultyCfg := &config.Config{Q: big.NewInt(0)}

	compute := ComputeS{}

	bigNum := new(big.Int)
	bigNum.SetString("10000000000000000000", 10)

	// Define the test cases
	testCases := []struct {
		name        string
		x           *big.Int
		k           *big.Int
		c           int64
		cfg         *config.Config
		expectedErr error
	}{
		{
			name:        "Test zero values",
			x:           big.NewInt(0),
			k:           big.NewInt(0),
			c:           0,
			cfg:         cfg,
			expectedErr: nil,
		},
		{
			name:        "Test negative values",
			x:           big.NewInt(-5),
			k:           big.NewInt(-7),
			c:           -11,
			cfg:         cfg,
			expectedErr: nil,
		},
		{
			name:        "Test with x value being larger",
			x:           big.NewInt(999),
			k:           big.NewInt(7),
			c:           12,
			cfg:         cfg,
			expectedErr: nil,
		},
		{
			name:        "Test with k value being larger",
			x:           big.NewInt(11),
			k:           big.NewInt(99),
			c:           33,
			cfg:         cfg,
			expectedErr: nil,
		},
		{
			name:        "Test large values",
			x:           bigNum,
			k:           bigNum,
			c:           math.MaxInt64,
			cfg:         cfg,
			expectedErr: nil,
		},
		{
			name:        "Test when cfg.Q equals to 0",
			cfg:         faultyCfg,
			x:           big.NewInt(11),
			k:           big.NewInt(99),
			c:           33,
			expectedErr: ErrZeroQ,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			mockRes := &interactor.AuthenticationChallengeResponse{
				C: tt.c,
			}
			_, actualErr := compute.Exec(tt.cfg, tt.x, tt.k, mockRes)
			if tt.expectedErr != nil {
				require.ErrorIs(t, actualErr, tt.expectedErr, "Expected error of type %v, but got %v", tt.expectedErr, actualErr)
				return
			}
			require.NoError(t, actualErr, "Expected no error, but got %v", actualErr)
		})
	}
}
