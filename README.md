# LRU Cache Implementation in Go

A Least Recently Used (LRU) cache implementation with both CLI and TCP server interfaces.

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

- `-port`: TCP server port (default: "7333", range: 1024-65535)
- `-capacity`: Cache capacity (default: 16, max: 256)
- `-buffer`: TCP buffer size in bytes (default: 256, range: 16-1024)
- `-only`: Run specific interface ("tcp" or "cli")

### Available commands:

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

## License

This project is licensed under the ISC License - see the [LICENSE](LICENSE) file for details.
