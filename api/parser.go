package api

import (
	"bytes"
	"fmt"
	"hash/fnv"
	"lru/src"
	"strconv"
	"strings"
)

type Cmd string
type Command struct {
	operation Cmd
	mapTitle  string
	mapKey    uint64
	key       uint64
	value     []byte
}

const (
	Cmd_CREATE    Cmd = "CREATE"
	Cmd_DESTROY   Cmd = "DESTROY"
	Cmd_LIST      Cmd = "LIST"
	Cmd_SET       Cmd = "SET"
	Cmd_GET       Cmd = "GET"
	Cmd_DEL       Cmd = "DEL"
	Cmd_PRINT     Cmd = "PRINT"
	Cmd_CLEAR     Cmd = "CLEAR"
	Cmd_CLEAR_ALL Cmd = "CLEAR_ALL"
	Cmd_HELP      Cmd = "HELP"
)

func hash(data []byte) uint64 {
	h := fnv.New64a()
	h.Write(data)
	return h.Sum64()
}

func splitBytes(input []byte) [][]byte {
	var args [][]byte
	start := 0
	for i := 0; i < len(input); i++ {
		if input[i] == ' ' {
			if start != i {
				args = append(args, input[start:i])
			}
			start = i + 1
		}
	}
	if start < len(input) {
		args = append(args, input[start:])
	}
	return args
}

func Parse(input []byte) (*Command, error) {
	args := splitBytes(input)
	if len(args) == 0 {
		return nil, fmt.Errorf("empty command")
	}

	opStr := strings.ToUpper(string(args[0]))
	cmd := &Command{operation: Cmd(opStr)}

	switch cmd.operation {
	case Cmd_CREATE:
		if len(args) != 3 {
			return nil, fmt.Errorf("usage: CREATE <cache_name> <capacity>")
		}
		cmd.mapTitle = string(args[1])
		cmd.mapKey = hash(args[1])
		cmd.value = []byte(args[2])

	case Cmd_LIST:
		if len(args) != 1 {
			return nil, fmt.Errorf("usage: LIST")
		}

	case Cmd_SET, Cmd_GET, Cmd_DEL, Cmd_PRINT, Cmd_CLEAR:
		if len(args) < 2 {
			return nil, fmt.Errorf("usage: %s <cache_name> [args...]", cmd.operation)
		}
		cmd.mapKey = hash(args[1])

		switch cmd.operation {
		case Cmd_SET:
			if len(args) < 4 {
				return nil, fmt.Errorf("usage: SET <cache_name> <key> <value>")
			}
			cmd.key = hash(args[2])
			cmd.value = bytes.Join(args[3:], []byte(""))
		case Cmd_GET, Cmd_DEL:
			if len(args) != 3 {
				return nil, fmt.Errorf("usage: %s <cache_name> <key>", cmd.operation)
			}
			cmd.key = hash(args[2])
		}

	case Cmd_HELP, Cmd_CLEAR_ALL:
		// No arguments

	default:
		return nil, fmt.Errorf("unknown command: %s", cmd.operation)
	}

	return cmd, nil
}

func Execute(cm *src.CacheManager, cmd *Command) (string, error) {
	switch cmd.operation {
	case Cmd_CREATE:
		capacity, err := strconv.ParseUint(string(cmd.value), 10, 8)
		if err != nil {
			return "", fmt.Errorf("invalid capacity: %s", cmd.value)
		}
		cm.CreateCache(cmd.mapTitle, cmd.mapKey, uint8(capacity))
		return "OK", nil

	case Cmd_DESTROY:
		if cache := cm.GetCache(cmd.mapKey); cache == nil {
			return "", fmt.Errorf("cache not found: %s", cmd.mapTitle)
		}
		cm.DestroyCache(cmd.mapKey)
		return "OK", nil

	case Cmd_LIST:
		names := cm.ListCaches()
		if len(names) == 0 {
			return "No caches", nil
		}
		return strings.Join(names, "\n"), nil

	case Cmd_SET, Cmd_GET, Cmd_DEL, Cmd_PRINT, Cmd_CLEAR:
		cache := cm.GetCache(cmd.mapKey)
		if cache == nil {
			return "", fmt.Errorf("cache not found: %s", cmd.mapTitle)
		}

		switch cmd.operation {
		case Cmd_SET:
			cache.Put(cmd.key, cmd.value)
			return "OK", nil
		case Cmd_GET:
			if value := cache.Get(cmd.key); value != nil {
				return string(value), nil
			}
			return "", fmt.Errorf("key not found")
		case Cmd_DEL:
			cache.Eject(cmd.key)
			return "OK", nil
		case Cmd_PRINT:
			return cache.Print(), nil
		case Cmd_CLEAR:
			cache.Clear()
			return "OK", nil
		}

	case Cmd_CLEAR_ALL:
		cm.ClearAllCaches()
		return "OK", nil

	case Cmd_HELP:
		return `Available commands:
CREATE <cache_name> <capacity>
DESTROY <cache_name>
LIST
SET <cache_name> <key> <value>
GET <cache_name> <key>
DEL <cache_name> <key>
PRINT <cache_name>
CLEAR <cache_name>
CLEAR_ALL
QUIT`, nil
	}

	return "", fmt.Errorf("unhandled command: %s", cmd.operation)
}
