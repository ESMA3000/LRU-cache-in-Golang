package api

import (
	"bufio"
	"fmt"
	"lru/src"
	"os"
)

func Cli(mgr *src.CacheManager) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("LRU Engine CLI")
	for {
		scanner.Scan()
		var input string = scanner.Text()

		cmd, err := Parse([]byte(input))
		if err != nil {
			fmt.Println("Error:", err)
			continue
		}

		if cmd.operation == "EXIT" || cmd.operation == "exit" {
			return
		}

		result, err := Execute(mgr, cmd)
		if err != nil {
			fmt.Println("Error:", err)
		} else {
			fmt.Println(result)
		}
	}
}
