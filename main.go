package main

import (
	"flag"
	"fmt"
	"lru/api"
)

func main() {
	port := flag.String("port", "7333", "Port to run the server on")
	apiMode := flag.String("api", "tcp", "API to use (tcp or cli)")
	capacity := flag.Int("capacity", 16, "Capacity of the cache (default: 16)")
	flag.Parse()
	if *capacity <= 0 {
		fmt.Println("Capacity must be non-negative.")
		return
	}
	if *capacity > 256 {
		fmt.Println("Capacity must be less than or equal to 256.")
		return
	}
	if *apiMode == "tcp" {
		api.ServerTCP(*port, *capacity)
	} else if *apiMode == "cli" {
		api.Cli(*capacity)
	} else {
		fmt.Println("Invalid API. Use 'tcp' or 'cli'.")
	}

}
