package api

import (
	"fmt"
	"lrue/src"
	"testing"
)

func TestExecute(t *testing.T) {
	cm := src.NewCacheManager[uint8, uint64, []byte]()
	tests := []struct {
		name    string
		cmd     *Command[uint64, []byte]
		want    string
		wantErr bool
	}{
		{
			name: "CREATE cache",
			cmd: &Command[uint64, []byte]{
				operation: Cmd_CREATE,
				mapTitle:  "test-cache",
				mapKey:    hash[uint64]([]byte("test-cache")),
				value:     []byte("5"),
			},
			want:    "OK",
			wantErr: false,
		},
		{
			name: "SET command",
			cmd: &Command[uint64, []byte]{
				operation: Cmd_SET,
				mapKey:    hash[uint64]([]byte("test-cache")),
				key:       hash[uint64]([]byte("test")),
				value:     []byte("value"),
			},
			want:    "OK",
			wantErr: false,
		},
		{
			name: "GET existing key",
			cmd: &Command[uint64, []byte]{
				operation: Cmd_GET,
				mapKey:    hash[uint64]([]byte("test-cache")),
				key:       hash[uint64]([]byte("test")),
			},
			want:    "value",
			wantErr: false,
		},
		{
			name: "GET non-existing key",
			cmd: &Command[uint64, []byte]{
				operation: Cmd_GET,
				mapKey:    hash[uint64]([]byte("test-cache")),
				key:       hash[uint64]([]byte("nonexistent")),
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "GET from non-existing cache",
			cmd: &Command[uint64, []byte]{
				operation: Cmd_GET,
				mapKey:    hash[uint64]([]byte("nonexistent-cache")),
				key:       hash[uint64]([]byte("test")),
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "LIST caches",
			cmd: &Command[uint64, []byte]{
				operation: Cmd_LIST,
			},
			want:    fmt.Sprintf("test-cache\nKey: %d, Value: %v", hash[uint64]([]byte("test")), []byte("value")),
			wantErr: false,
		},
		{
			name: "DEL command",
			cmd: &Command[uint64, []byte]{
				operation: Cmd_DEL,
				mapKey:    hash[uint64]([]byte("test-cache")),
				key:       hash[uint64]([]byte("test")),
			},
			want:    "OK",
			wantErr: false,
		},
		{
			name: "DESTROY cache",
			cmd: &Command[uint64, []byte]{
				operation: Cmd_DESTROY,
				mapKey:    hash[uint64]([]byte("test-cache")),
			},
			want:    "OK",
			wantErr: false,
		},
		{
			name: "CLEAR_ALL caches",
			cmd: &Command[uint64, []byte]{
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
