package api

import (
	"fmt"
	"lru/src"
	"net"
)

func handleConnection(conn net.Conn, cache *src.LRUCache) {
	defer conn.Close()
	buf := make([]byte, 1024)
	conn.Write([]byte("Connected to lru cachemanager\r\n"))
	for {
		n, err := conn.Read(buf)
		if err != nil {
			fmt.Printf("Connection closed by client\n")
			break
		}

		cmd, err := Parse(string(buf[:n]))
		if err != nil {
			conn.Write(fmt.Appendf(nil, "ERR %s\r\n", err))
			continue
		}

		result, err := Execute(cache, cmd)
		if err != nil {
			conn.Write(fmt.Appendf(nil, "ERR %s\r\n", err))
		} else {
			conn.Write(fmt.Appendf(nil, "%s\r\n", result))
		}
	}
}

func ServerTCP(port string, capacity uint8) {
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Printf("Error starting listener: %v\n", err)
		return
	}
	defer listener.Close()

	cache := src.InitLRU(capacity)
	fmt.Printf("Server started on port %s\n", port)

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}
		fmt.Printf("New client connected\n")
		go handleConnection(conn, &cache)
	}
}
