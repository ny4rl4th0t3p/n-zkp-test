package app

import (
	"math/big"
	"testing"

	"practical-case-test/config"

	"github.com/stretchr/testify/require"
)

func TestRegister_Exec(t *testing.T) {
	r := new(Register)
	testCases := []struct {
		name        string
		cfg         *config.Config
		password    *big.Int
		expectedY1  *big.Int
		expectedY2  *big.Int
		expectedErr error
	}{
		{
			name:        "Nil Config Test",
			cfg:         nil,
			password:    big.NewInt(10),
			expectedY1:  nil,
			expectedY2:  nil,
			expectedErr: ErrNilConfig,
		},
		{
			name: "Zero Q Config Test",
			cfg: &config.Config{
				Q: new(big.Int),
			},
			password:    big.NewInt(10),
			expectedY1:  nil,
			expectedY2:  nil,
			expectedErr: ErrZeroQ,
		},
		{
			name: "Positive Number Test",
			cfg: &config.Config{
				Q: big.NewInt(10),
				G: big.NewInt(2),
				H: big.NewInt(2),
			},
			password:    big.NewInt(10),
			expectedY1:  big.NewInt(4),
			expectedY2:  big.NewInt(4),
			expectedErr: nil,
		},
		{
			name: "Zero Number Test",
			cfg: &config.Config{
				Q: big.NewInt(10),
				G: big.NewInt(2),
				H: big.NewInt(2),
			},
			password:    big.NewInt(0),
			expectedY1:  big.NewInt(1),
			expectedY2:  big.NewInt(1),
			expectedErr: nil,
		},
		{
			name: "Negative Number Test",
			cfg: &config.Config{
				Q: big.NewInt(10),
				G: big.NewInt(2),
				H: big.NewInt(2),
			},
			password:    big.NewInt(-10),
			expectedY1:  big.NewInt(1),
			expectedY2:  big.NewInt(1),
			expectedErr: nil,
		},
	}

	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			y1, y2, err := r.Exec(tt.cfg, tt.password)
			require.Equal(t, tt.expectedErr, err)
			require.Equal(t, tt.expectedY1, y1)
			require.Equal(t, tt.expectedY2, y2)
		})
	}
}
