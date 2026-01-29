package utils

import (
	"testing"
	"time"
)

func TestNewSSHManager(t *testing.T) {
	type args struct {
		host    string
		port    string
		user    string
		keyPath string
	}
	tests := []struct {
		name    string
		args    args
		want    *SSHManager
		wantErr bool
	}{
		{
			name: "测试连接 SSH",
			args: args{
				host:    "192.168.106.0",
				port:    "22",
				user:    "root",
				keyPath: "id_rsa",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			connect, err := NewSSHManager(tt.args.host, tt.args.port, tt.args.user, []byte("keyinfo"), 30*time.Second)
			if (err != nil) != tt.wantErr {
				t.Errorf("NewSSHManager() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			defer connect.Close()
			//测试超时
			output, err := connect.ExecuteCommand("sleep 5 && ps -ef |grep ios", time.Second*7)
			if err != nil {
				t.Errorf("ExecuteCommand() error = %v", err)
				return
			}
			t.Logf("ExecuteCommand() output = %v", output)
		})
	}
}
