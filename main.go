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
	capacity := flag.Int("capacity", 16, "Capacity of the cache")
	bufferSize := flag.Int("buffer", 256, "Buffer size for TCP connections")
	only := flag.String("only", "", "Run only either TCP server or CLI")
	flag.Parse()
	portNum, err := strconv.Atoi(*port)
	if err != nil || portNum < 1024 || portNum > int(^uint16(0)) {
		fmt.Println("Port must be a number between 1024 and 65535.")
		return
	}
	if *capacity <= 0 || *capacity > int(^uint8(0)) {
		fmt.Println("Capacity must be non-negative and can not exceed UINT8_MAX (256)")
		return
	}
	if *bufferSize <= 16 || *bufferSize > 1024 {
		fmt.Println("Buffer size must be between 16 and 1024 bytes.")
		return
	}
	cache := src.InitLRU(uint8(*capacity))
	switch *only {
	case "tcp":
		api.ServerTCP(*port, uint16(*bufferSize), &cache)
		return
	case "cli":
		api.Cli(&cache)
		return
	default:
		go api.ServerTCP(*port, uint16(*bufferSize), &cache)
		api.Cli(&cache)
	}
}
