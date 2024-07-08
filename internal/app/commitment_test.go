package app

import (
	"testing"

	"practical-case-test/config"

	"github.com/stretchr/testify/require"
)

func TestCommitment_Exec(t *testing.T) {
	cfg := config.LoadConfig()

	testCases := []struct {
		name    string
		cfg     *config.Config
		wantErr bool
	}{
		{
			name:    "Valid case",
			cfg:     cfg,
			wantErr: false,
		},
		{
			name:    "Invalid case: no Config",
			cfg:     nil,
			wantErr: true,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			commitment := Commitment{}
			gotResult, err := commitment.Exec(tt.cfg)
			if tt.wantErr {
				require.Error(t, err, "Exec() should return an error.")
			} else {
				require.NoError(t, err, "Exec() should not return an error.")
				require.NotNil(t, gotResult, "Exec() should return a valid result.")
			}
		})
	}
}

func Test_generateRndCommitment(t *testing.T) {
	cfg := config.LoadConfig()

	testCases := []struct {
		name    string
		cfg     *config.Config
		wantErr bool
	}{
		{
			name:    "Valid case",
			cfg:     cfg,
			wantErr: false,
		},
		{
			name:    "Invalid case: no Config",
			cfg:     nil,
			wantErr: true,
		},
	}
	for _, tt := range testCases {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			gotR1, gotR2, gotK, err := generateRndCommitment(tt.cfg)

			if tt.wantErr {
				require.Error(t, err, "generateRndCommitment() should return an error.")
			} else {
				require.NoError(t, err, "generateRndCommitment() should not return an error.")
				require.NotNil(t, gotR1, "generateRndCommitment() should return a valid r1.")
				require.NotNil(t, gotR2, "generateRndCommitment() should return a valid r2.")
				require.NotNil(t, gotK, "generateRndCommitment() should return a valid k.")
			}
		})
	}
}
