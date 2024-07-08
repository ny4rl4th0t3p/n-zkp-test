package auth

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestNewSession(t *testing.T) {
	tests := []struct {
		name    string
		args    Session
		wantErr bool
	}{
		{
			name: "Valid: regular condition",
			args: Session{
				id:             uuid.New(),
				userID:         "user_id",
				loginTimestamp: 1598896296,
			},
			wantErr: false,
		},
		{
			name: "Invalid: empty uuid",
			args: Session{
				id:             uuid.UUID{},
				userID:         "user_id",
				loginTimestamp: 1598896296,
			},
			wantErr: true,
		},
		{
			name: "Invalid: empty userID",
			args: Session{
				id:             uuid.New(),
				userID:         "",
				loginTimestamp: 1598896296,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			_, err := NewSession(tt.args.id, tt.args.userID, tt.args.loginTimestamp)
			if tt.wantErr {
				require.Error(t, err, "NewSession() error = %v, wantErr %v", err, tt.wantErr)
			} else {
				require.NoError(t, err, "NewSession() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestSession_IsValid(t *testing.T) {
	type fields struct {
		id             uuid.UUID
		userID         string
		loginTimestamp int64
	}
	tests := []struct {
		name   string
		fields fields
		want   bool
	}{
		{
			name: "Valid: regular condition",
			fields: fields{
				id:             uuid.New(),
				userID:         "user_id",
				loginTimestamp: 1598896296,
			},
			want: true,
		},
		{
			name: "Invalid: empty uuid",
			fields: fields{
				id:             uuid.UUID{},
				userID:         "user_id",
				loginTimestamp: 1598896296,
			},
			want: false,
		},
		{
			name: "Invalid: empty userID",
			fields: fields{
				id:             uuid.New(),
				userID:         "",
				loginTimestamp: 1598896296,
			},
			want: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			s := Session{
				id:             tt.fields.id,
				userID:         tt.fields.userID,
				loginTimestamp: tt.fields.loginTimestamp,
			}
			got := s.IsValid()
			require.Equal(t, tt.want, got, "IsValid() = %v, want %v", got, tt.want)
		})
	}
}
