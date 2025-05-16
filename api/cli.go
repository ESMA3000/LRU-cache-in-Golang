package api

import (
	"bufio"
	"fmt"
	"lru/src"
	"os"
	"strings"
)

func Cli(capacity int) {
	scanner := bufio.NewScanner(os.Stdin)
	cache := src.InitLRU(capacity)

	fmt.Println("LRU Cache CLI")
	fmt.Println("Commands: put <key> <value>, get <key>, eject <key>, print, quit")
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
		case "eject":
			if len(args) != 2 {
				fmt.Println("Usage: eject <key>")
				continue
			}
			cache.Eject(args[1])
			fmt.Println("Ejected from cache")

		case "print":
			cache.Print()

		case "clear":
			cache.Clear()
			fmt.Println("Cache cleared")

		case "help":
			fmt.Println("Available commands: put <key> <value>, get <key>, eject <key>, print, clear, quit")

		case "quit":
			return

		default:
			fmt.Println("Unknown command. Available commands: put, eject, get, print, clear, quit")
		}
	}
}
