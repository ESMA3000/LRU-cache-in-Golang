package api

import (
	"bufio"
	"context"
	"fmt"
	"lru/src"
	"os"
)

func Cli(ctx context.Context, mgr *src.CacheManager) {
	scanner := bufio.NewScanner(os.Stdin)
	fmt.Println("LRU Engine CLI")
	for {
		select {
		case <-ctx.Done():
			return
		default:
			if !scanner.Scan() {
				return
			}
			var input string = scanner.Text()

			cmd, err := Parse([]byte(input))
			if err != nil {
				fmt.Println("Error:", err)
				continue
			}
			result, err := Execute(mgr, cmd)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println(result)
			}
		}
	}
}
