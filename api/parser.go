package api

import (
	"fmt"
	"hash/fnv"
	"lru/src"
	"strings"
	"unsafe"
)

type Command struct {
	operation string
	key       uint64
	value     []byte
}

func hashString(s string) uint64 {
	h := fnv.New64a()
	h.Write([]byte(s))
	return h.Sum64()
}

func Parse(input string) (*Command, error) {
	args := strings.Fields(input)
	if len(args) == 0 {
		return nil, fmt.Errorf("empty command")
	}

	op := strings.ToUpper(args[0])
	cmd := &Command{operation: op}

	switch op {
	case "PUT", "SET":
		if len(args) != 3 {
			return nil, fmt.Errorf("usage: PUT/SET <key> <value>")
		}
		cmd.key = hashString(args[1])
		cmd.value = unsafe.Slice(unsafe.StringData(args[2]), len(args[2]))

	case "GET":
		if len(args) != 2 {
			return nil, fmt.Errorf("usage: GET <key>")
		}
		cmd.key = hashString(args[1])

	case "EJECT", "DEL":
		if len(args) != 2 {
			return nil, fmt.Errorf("usage: EJECT/DEL <key>")
		}
		cmd.key = hashString(args[1])

	case "PRINT", "CLEAR", "QUIT", "HELP":
		// Commands without arguments
		break

	default:
		return nil, fmt.Errorf("unknown command: %s", op)
	}

	return cmd, nil
}

func Execute(cache *src.LRUMap, cmd *Command) (string, error) {
	switch cmd.operation {
	case "PUT", "SET":
		cache.Put(cmd.key, cmd.value)
		return "OK", nil

	case "GET":
		if value := cache.Get(cmd.key); value != nil {
			return fmt.Sprintf("%v", value), nil
		}
		return "", fmt.Errorf("key not found")

	case "EJECT", "DEL":
		cache.Eject(cmd.key)
		return "OK", nil

	case "PRINT":
		if output := cache.Print(); output == "" {
			return "Cache is empty", nil
		} else {
			return output, nil
		}

	case "CLEAR":
		cache.Clear()
		return "OK", nil

	case "HELP":
		return "Available commands: PUT/SET <key> <value>, GET <key>, EJECT/DEL <key>, PRINT, CLEAR, QUIT", nil
	}

	return "", fmt.Errorf("unhandled command: %s", cmd.operation)
}
