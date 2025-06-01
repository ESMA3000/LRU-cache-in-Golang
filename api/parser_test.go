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
				mapTitle:  "test-cache",
				mapKey:    hash([]byte("test-cache")),
				value:     []byte("5"),
			},
			want:    "OK",
			wantErr: false,
		},
		{
			name: "SET command",
			cmd: &Command{
				operation: Cmd_SET,
				mapKey:    hash([]byte("test-cache")),
				key:       hash([]byte("test")),
				value:     []byte("value"),
			},
			want:    "OK",
			wantErr: false,
		},
		{
			name: "GET existing key",
			cmd: &Command{
				operation: Cmd_GET,
				mapKey:    hash([]byte("test-cache")),
				key:       hash([]byte("test")),
			},
			want:    "value",
			wantErr: false,
		},
		{
			name: "GET non-existing key",
			cmd: &Command{
				operation: Cmd_GET,
				mapKey:    hash([]byte("test-cache")),
				key:       hash([]byte("nonexistent")),
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "GET from non-existing cache",
			cmd: &Command{
				operation: Cmd_GET,
				mapKey:    hash([]byte("nonexistent-cache")),
				key:       hash([]byte("test")),
			},
			want:    "",
			wantErr: true,
		},
		/* 		{
			name: "LIST caches",
			cmd: &Command{
				operation: Cmd_LIST,
			},
			want:    "test-cache\nKey: 18007334074686647077, Value: [118 97 108 117 101]",
			wantErr: false,
		}, */
		{
			name: "DEL command",
			cmd: &Command{
				operation: Cmd_DEL,
				mapKey:    hash([]byte("test-cache")),
				key:       hash([]byte("test")),
			},
			want:    "OK",
			wantErr: false,
		},
		{
			name: "DESTROY cache",
			cmd: &Command{
				operation: Cmd_DESTROY,
				mapKey:    hash([]byte("test-cache")),
			},
			want:    "OK",
			wantErr: false,
		},
		{
			name: "CLEAR_ALL caches",
			cmd: &Command{
				operation: Cmd_CLEAR_ALL,
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
