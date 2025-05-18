package api

import (
	"bufio"
	"fmt"
	"lru/src"
	"os"
)

func Cli(cache *src.LRUMap) {
	scanner := bufio.NewScanner(os.Stdin)

	fmt.Println("LRU Cache CLI")
	fmt.Println("Commands: put <key> <value>, get <key>, eject <key>, print, quit")
	for {
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

		result, err := Execute(cache, cmd)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println(result)
		}
	}
}
