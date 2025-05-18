package api

import (
	"fmt"
	"lru/src"
	"strings"
)

type Command struct {
	Operation string
	Key       string
	Value     string
}

func Parse(input string) (*Command, error) {
	args := strings.Fields(input)
	if len(args) == 0 {
		return nil, fmt.Errorf("empty command")
	}

	op := strings.ToUpper(args[0])
	cmd := &Command{Operation: op}

	switch op {
	case "PUT", "SET":
		if len(args) != 3 {
			return nil, fmt.Errorf("usage: PUT/SET <key> <value>")
		}
		cmd.Key = args[1]
		cmd.Value = args[2]

	case "GET":
		if len(args) != 2 {
			return nil, fmt.Errorf("usage: GET <key>")
		}
		cmd.Key = args[1]

	case "EJECT", "DEL":
		if len(args) != 2 {
			return nil, fmt.Errorf("usage: EJECT/DEL <key>")
		}
		cmd.Key = args[1]

	case "PRINT", "CLEAR", "QUIT", "HELP":
		// Commands without arguments
		break

	default:
		return nil, fmt.Errorf("unknown command: %s", op)
	}

	return cmd, nil
}

func Execute(cache *src.LRUCache, cmd *Command) (string, error) {
	switch cmd.Operation {
	case "PUT", "SET":
		cache.Put(cmd.Key, []byte(cmd.Value))
		return "OK", nil

	case "GET":
		if value := cache.Get(cmd.Key); value != nil {
			return fmt.Sprintf("%v", value), nil
		}
		return "", fmt.Errorf("key not found")

	case "EJECT", "DEL":
		cache.Eject(cmd.Key)
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

	return "", fmt.Errorf("unhandled command: %s", cmd.Operation)
}
