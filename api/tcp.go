package api

import (
	"fmt"
	"lru/src"
	"net"
	"sync"
	"sync/atomic"
)

const maxConnections = 256

var bufferPool sync.Pool

func handleConnection(conn net.Conn, bufferSize uint16, mgr *src.CacheManager) {
	if bufferPool.New == nil {
		bufferPool.New = func() any {
			b := make([]byte, bufferSize)
			return &b
		}
	}
	buf := bufferPool.Get().(*[]byte)
	defer func() {
		conn.Close()
		bufferPool.Put(buf)
	}()

	conn.Write([]byte("Connected to lru cachemanager\r\n"))
	for {
		input := *buf
		n, err := conn.Read(input)
		if err != nil {
			fmt.Printf("Connection closed by client\n")
			break
		}

		cmd, err := Parse(input[:n])
		if err != nil {
			conn.Write(fmt.Appendf(nil, "ERR %s\r\n", err))
			continue
		}

		if cmd.operation == Cmd_EXIT {
			break
		}

		result, err := Execute(mgr, cmd)
		if err != nil {
			conn.Write(fmt.Appendf(nil, "ERR %s\r\n", err))
		} else {
			conn.Write(fmt.Appendf(nil, "%s\r\n", result))
		}
	}
}

func ServerTCP(port string, bufferSize uint16, mgr *src.CacheManager) {
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
			handleConnection(conn, bufferSize, mgr)
			atomic.AddInt32(&activeConnections, -1)
		}()
	}
}
