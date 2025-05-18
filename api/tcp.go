package api

import (
	"fmt"
	"lru/src"
	"net"
	"sync/atomic"
)

const maxConnections = 256

func handleConnection(conn net.Conn, bufferSize uint16, cache *src.LRUCache) {
	defer conn.Close()
	buf := make([]byte, bufferSize)
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

func ServerTCP(port string, bufferSize uint16, cache *src.LRUCache) {
	var activeConnections int32 = 0
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		fmt.Printf("Error starting listener: %v\n", err)
		return
	}
	defer listener.Close()

	fmt.Printf("Server started on port %s\n", port)

	for {
		if atomic.LoadInt32(&activeConnections) >= maxConnections {
			continue
		}

		conn, err := listener.Accept()
		if err != nil {
			fmt.Printf("Error accepting connection: %v\n", err)
			continue
		}

		atomic.AddInt32(&activeConnections, 1)
		fmt.Printf("New client connected (active: %d)\n", atomic.LoadInt32(&activeConnections))
		go func() {
			handleConnection(conn, bufferSize, cache)
			atomic.AddInt32(&activeConnections, -1)
		}()
	}
}
