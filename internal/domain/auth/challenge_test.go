package auth

import (
	"math/big"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

var (
	testUUID   = uuid.New()
	testBigInt = big.NewInt(1234)
)

func TestChallenge_IsValid(t *testing.T) {
	type fields struct {
		userID    string
		authID    uuid.UUID
		c         *big.Int
		r1        int64
		r2        int64
		timestamp int64
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Test Valid Challenge",
			fields: fields{
				userID: "test_user",
				authID: testUUID,
				c:      testBigInt,
			},
			want: true,
		},
		{
			name: "Test Invalid Challenge with empty userID",
			fields: fields{
				userID: "",
				authID: testUUID,
				c:      testBigInt,
			},
			want: false,
		},
		{
			name: "Test Invalid Challenge with nil authID",
			fields: fields{
				userID: "test_user",
				authID: uuid.Nil,
				c:      testBigInt,
			},
			want: false,
		},
		{
			name: "Test Invalid Challenge with nil bigInt",
			fields: fields{
				userID: "test_user",
				authID: testUUID,
				c:      nil,
			},
			want: false,
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			c := Challenge{
				userID:    tt.fields.userID,
				authID:    tt.fields.authID,
				c:         tt.fields.c,
				r1:        tt.fields.r1,
				r2:        tt.fields.r2,
				timestamp: tt.fields.timestamp,
			}
			if got := c.IsValid(); got != tt.want {
				t.Errorf("IsValid() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestNewChallenge(t *testing.T) {
	type args struct {
		c         *big.Int
		userID    string
		timestamp int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "Test New Challenge Successfully",
			args: args{
				c:      testBigInt,
				userID: "test_user",
			},
			wantErr: false,
		},
		{
			name: "Test New Challenge with nil bigInt",
			args: args{
				c:      nil,
				userID: "test_user",
			},
			wantErr: true,
		},
		{
			name: "Test New Challenge with empty userID",
			args: args{
				c:      testBigInt,
				userID: "",
			},
			wantErr: true,
		}}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got, err := NewChallenge(tt.args.c, tt.args.userID, 2, 4, tt.args.timestamp)
			if tt.wantErr {
				require.Error(t, err, "NewChallenge() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			require.NoError(t, err, "NewChallenge() error = %v, wantErr %v", err, tt.wantErr)
			require.Equal(t, tt.args.c, got.c, "NewChallenge() got = %v, want %v", got.c, tt.args.c)
			require.Equal(t, tt.args.userID, got.userID, "NewChallenge() got = %v, want %v", got.userID, tt.args.userID)
		})
	}
}
