package auth

import (
	"testing"
	"time"
)

func TestGenerateTokenWithTime(t *testing.T) {
	type args struct {
		id uint
		t  time.Time
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Valid user ID and time",
			args: args{
				id: 123,
				t:  time.Now(),
			},
			want:    "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjoxMjMsImV4cCI6MTc3MzQzMTYwMCwiaWF0IjoxNjcwODQ0ODAwfQ.b7Y-k5k5h_jM3yH1XzXv5rW76c5eJ3mP4k9tU8_X6A",
			wantErr: false,
		},
		{
			name: "Invalid user ID",
			args: args{
				id: 0,
				t:  time.Now(),
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Expired time",
			args: args{
				id: 123,
				t:  time.Now().Add(-time.Hour),
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Missing JWT secret",
			args: args{
				id: 123,
				t:  time.Now(),
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Invalid JWT secret",
			args: args{
				id: 123,
				t:  time.Now(),
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "Unexpected error",
			args: args{
				id: 123,
				t:  time.Now(),
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GenerateTokenWithTime(tt.args.id, tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("GenerateTokenWithTime() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GenerateTokenWithTime() got = %v, want %v", got, tt.want)
			}
		})
	}
}
