package utils

import (
	"reflect"
	"testing"
)

func TestGetNetworkInterfaces(t *testing.T) {
	type args struct {
		filter []string
	}
	tests := []struct {
		name    string
		args    args
		want    map[string]map[string]string
		wantErr bool
	}{
		{
			name:    "GetNetworkInterfaces",
			args:    args{filter: []string{"eth0"}},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetNetworkInterfaces(tt.args.filter)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetNetworkInterfaces() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetNetworkInterfaces() got = %v, want %v", got, tt.want)
			}
		})
	}
}
