package src

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func Cmdl() {
	var capacity int
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("LRU Cache Terminal")
	for {
		fmt.Print("Enter the size of the cache: ")
		scanner.Scan()
		input := scanner.Text()
		parsedCapacity, err := strconv.Atoi(input)
		if err != nil {
			fmt.Println("Invalid capacity. Please enter a number.")
			continue
		}
		if parsedCapacity <= 0 {
			fmt.Println("Capacity must be positive. Please try again.")
			continue
		}
		capacity = parsedCapacity
		break
	}
	cache := InitLRU(capacity)

	fmt.Println("Commands: put <key> <value>, get <key>, print, quit")
	for {
		fmt.Print("> ")
		scanner.Scan()
		input := scanner.Text()
		args := strings.Fields(input)

		if len(args) == 0 {
			continue
		}

		switch args[0] {
		case "put":
			if len(args) != 3 {
				fmt.Println("Usage: put <key> <value>")
				continue
			}
			cache.Put(args[1], args[2])
			fmt.Println("Added to cache")

		case "get":
			if len(args) != 2 {
				fmt.Println("Usage: get <key>")
				continue
			}
			if value := cache.Get(args[1]); value != nil {
				fmt.Printf("Value: %v\n", value)
			} else {
				fmt.Println("Key not found")
			}

		case "print":
			cache.Print()

		case "clear":
			cache.Clear()
			fmt.Println("Cache cleared")

		case "quit":
			return

		default:
			fmt.Println("Unknown command. Available commands: put, eject, get, print, clear, quit")
		}
	}
}
