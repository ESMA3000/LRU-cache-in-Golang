package main

import (
	"context"
	"flag"
	"fmt"
	"lrue/api"
	"lrue/src"
	"os/signal"
	"strconv"
	"syscall"
)

type Config struct {
	port       string
	bufferSize int
	only       string
}

func main() {
	config := parseFlags()
	if err := validateConfig(config); err != nil {
		src.FatalError("Invalid configuration", err)
	}
	mgr := src.NewCacheManager[uint8, uint64, []byte]()

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	switch config.only {
	case "tcp":
		api.ServerTCP(config.port, uint16(config.bufferSize), mgr)
	case "cli":
		api.Cli(ctx, mgr)
	default:
		go api.ServerTCP(config.port, uint16(config.bufferSize), mgr)
		api.Cli(ctx, mgr)
	}

	<-ctx.Done()
	fmt.Println("Shutting down gracefully...")
	mgr.ClearAllCaches()
}

func parseFlags() Config {
	port := flag.String("port", "7333", "Port to run the server on")
	bufferSize := flag.Int("buffer", 256, "Buffer size for TCP connections")
	only := flag.String("only", "", "Run only either TCP server or CLI")
	flag.Parse()
	return Config{
		port:       *port,
		bufferSize: *bufferSize,
		only:       *only,
	}
}

func validateConfig(config Config) error {
	portNum, err := strconv.Atoi(config.port)
	if err != nil || portNum < 1024 || portNum > int(^uint16(0)) {
		return fmt.Errorf("port must be a number between 1024 and 65535")
	}
	if config.bufferSize <= 16 || config.bufferSize > 1024 {
		return fmt.Errorf("buffer size must be between 16 and 1024 bytes")
	}
	return nil
}
