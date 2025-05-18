package api

import (
	"bufio"
	"fmt"
	"lru/src"
	"os"
)

func Cli(capacity uint8) {
	scanner := bufio.NewScanner(os.Stdin)
	cache := src.InitLRU(capacity)

	fmt.Println("LRU Cache CLI")
	fmt.Println("Commands: put <key> <value>, get <key>, eject <key>, print, quit")
	for {
		fmt.Print("> ")
		scanner.Scan()
		var input string = scanner.Text()

		cmd, err := Parse(input)
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		if cmd.Operation == "QUIT" || cmd.Operation == "quit" {
			return
		}

		result, err := Execute(&cache, cmd)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println(result)
		}
	}
}
