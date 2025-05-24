package main

import (
	"flag"
	"fmt"
	"lru/api"
	"lru/src"
	"strconv"
)

func main() {
	port := flag.String("port", "7333", "Port to run the server on")
	bufferSize := flag.Int("buffer", 256, "Buffer size for TCP connections")
	only := flag.String("only", "", "Run only either TCP server or CLI")
	flag.Parse()
	portNum, err := strconv.Atoi(*port)
	if err != nil || portNum < 1024 || portNum > int(^uint16(0)) {
		fmt.Println("Port must be a number between 1024 and 65535.")
		return
	}
	if *bufferSize <= 16 || *bufferSize > 1024 {
		fmt.Println("Buffer size must be between 16 and 1024 bytes.")
		return
	}
	mgr := src.NewCacheManager()
	switch *only {
	case "tcp":
		api.ServerTCP(*port, uint16(*bufferSize), mgr)
		return
	case "cli":
		api.Cli(mgr)
		return
	default:
		go api.ServerTCP(*port, uint16(*bufferSize), mgr)
		api.Cli(mgr)
	}
}
