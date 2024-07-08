package auth

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNewUser(t *testing.T) {
	type args struct {
		user string
		y1   int64
		y2   int64
	}
	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name:    "Valid Case",
			args:    args{user: "valid_user", y1: 10, y2: 20},
			wantErr: false,
		},
		{
			name:    "Empty username",
			args:    args{user: "", y1: 10, y2: 20},
			wantErr: true,
		},
		{
			name:    "Negative y1",
			args:    args{user: "valid_user", y1: -10, y2: 20},
			wantErr: true,
		},
		{
			name:    "Negative y2",
			args:    args{user: "valid_user", y1: 10, y2: -20},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := NewUser(tt.args.user, tt.args.y1, tt.args.y2)
			if tt.wantErr {
				require.Error(t, err, "NewUser() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				require.NoError(t, err, "NewUser() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestUser_IsValid(t *testing.T) {
	tests := []struct {
		name string
		user User
		want bool
	}{
		{
			name: "Valid Case",
			user: User{userID: "valid_user", y1: 10, y2: 20},
			want: true,
		},
		{
			name: "Empty Username",
			user: User{userID: "", y1: 10, y2: 20},
			want: false,
		},
		{
			name: "Negative y1",
			user: User{userID: "valid_user", y1: -10, y2: 20},
			want: false,
		},
		{
			name: "Negative y2",
			user: User{userID: "valid_user", y1: 10, y2: -20},
			want: false,
		},
		{
			name: "Negative y1 and y2",
			user: User{userID: "valid_user", y1: -10, y2: -20},
			want: false,
		},
		{
			name: "Empty username and negative y1/y2",
			user: User{userID: "", y1: -10, y2: -20},
			want: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			got := tt.user.IsValid()
			require.Equal(t, tt.want, got, tt.want, "IsValid() = %v, want %v", got, tt.want)
		})
	}
}
