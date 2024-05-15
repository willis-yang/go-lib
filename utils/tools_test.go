package utils

import (
	"context"
	"testing"
)

func TestGetRemoteIp(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{
			name: "test",
			args: args{
				ctx: context.Background(),
			},
			want:    "",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetRemoteIp(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetRemoteIp() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetRemoteIp() got = %v, want %v", got, tt.want)
			}
		})
	}
}
