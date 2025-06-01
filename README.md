# LRU Cache Implementation in Go

A Least Recently Used (LRU) cache implementation with both CLI and TCP server interfaces.

## Usage

```bash
go run main.go
```

### TCP Server Example

Start the server:
```bash
go run main.go -only tcp -port 7333 -buffer 256
```

Connect using netcat or telnet:
```bash
nc localhost 7333
```

## Command Line Arguments

- `-port`: TCP server port (default: "7333", range: 1024-65535)
- `-buffer`: TCP buffer size in bytes (default: 256, range: 16-1024)
- `-only`: Run specific interface ("tcp" or "cli")

### Available Commands

Cache Management:
- `CREATE <cache_name> <capacity>`: Create a new cache with specified capacity
- `DESTROY <cache_name>`: Remove a cache instance
- `LIST`: Show all available caches

Cache Operations:
- `SET <cache_name> <key> <value>`: Add or update a key-value pair in specified cache
- `GET <cache_name> <key>`: Retrieve a value by key from specified cache
- `DEL <cache_name> <key>`: Remove a key-value pair from specified cache
- `PRINT <cache_name>`: Display specified cache contents
- `CLEAR <cache_name>`: Remove all entries from specified cache
- `CLEAR_ALL`: Clear all caches
- `HELP`: Show available commands

## Implementation Details

- Multiple independent cache instances
- Each cache uses a doubly-linked list to maintain access order
- Most recent items are moved to the head of the list
- Least recently used items are removed when capacity is reached

## License

This project is licensed under the GNU Affero General Public License v3.0 - see the [LICENSE](LICENSE) file for details.