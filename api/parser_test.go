package api

import (
	"lru/src"
	"testing"
)

func TestExecute(t *testing.T) {
	cm := src.NewCacheManager()
	tests := []struct {
		name    string
		cmd     *Command
		want    string
		wantErr bool
	}{
		{
			name: "CREATE cache",
			cmd: &Command{
				operation: Cmd_CREATE,
				cacheName: "test-cache",
				value:     []byte("5"),
			},
			want:    "OK",
			wantErr: false,
		},
		{
			name: "SET command",
			cmd: &Command{
				operation: Cmd_SET,
				cacheName: "test-cache",
				key:       hashString("test"),
				value:     []byte("value"),
			},
			want:    "OK",
			wantErr: false,
		},
		{
			name: "GET existing key",
			cmd: &Command{
				operation: Cmd_GET,
				cacheName: "test-cache",
				key:       hashString("test"),
			},
			want:    "value",
			wantErr: false,
		},
		{
			name: "GET non-existing key",
			cmd: &Command{
				operation: Cmd_GET,
				cacheName: "test-cache",
				key:       hashString("nonexistent"),
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "GET from non-existing cache",
			cmd: &Command{
				operation: Cmd_GET,
				cacheName: "nonexistent-cache",
				key:       hashString("test"),
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "LIST caches",
			cmd: &Command{
				operation: Cmd_LIST,
			},
			want:    "test-cache",
			wantErr: false,
		},
		{
			name: "DEL command",
			cmd: &Command{
				operation: Cmd_DEL,
				cacheName: "test-cache",
				key:       hashString("test"),
			},
			want:    "OK",
			wantErr: false,
		},
		{
			name: "DESTROY cache",
			cmd: &Command{
				operation: Cmd_DESTROY,
				cacheName: "test-cache",
			},
			want:    "OK",
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := Execute(cm, tt.cmd)
			if (err != nil) != tt.wantErr {
				t.Errorf("Execute() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("Execute() = %v, want %v", got, tt.want)
			}
		})
	}
}
