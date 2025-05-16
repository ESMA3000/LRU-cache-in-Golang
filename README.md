# LRU Cache Implementation in Go

A flexible Least Recently Used (LRU) cache implementation with both CLI and TCP server interfaces.

## Usage

### Command Line Interface

```bash
go run main.go -api cli -capacity 16
```

### TCP Server

Start the server:
```bash
go run main.go -api tcp -port 7333 -capacity 16
```

Connect using netcat or telnet:
```bash
nc localhost 7333
```

## Command Line Arguments

- `-api`: Interface to use (`tcp` or `cli`, default: `tcp`)
- `-port`: Port for TCP server (default: `7333`)
- `-capacity`: Cache capacity (default: `16`, max: `256`)

### Available CLI commands:

- `PUT/SET <key> <value>`: Add or update a key-value pair
- `GET <key>`: Retrieve a value by key
- `EJECT/DEL <key>`: Remove a key-value pair
- `PRINT`: Display all cache contents
- `CLEAR`: Remove all entries
- `QUIT`: Exit the program
- `HELP`: Show available commands

## Implementation Details

- The cache uses a doubly-linked list to maintain access order
- Most recent items are moved to the head of the list
- Least recently used items are removed when capacity is reached
- O(1) time complexity for all operations

## License

This project is licensed under the ISC License - see the [LICENSE](LICENSE) file for details.
