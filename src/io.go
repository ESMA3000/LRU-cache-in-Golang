package src

import (
	"fmt"
	"log"
	"os"
	"strings"
)

func FatalError(message string, err error) {
	log.Fatalf("%s %v\n", message, err)
}

func LogErrorConsole(e error) {
	log.Printf("Error: %v\n", e)
}

func LogError(e error) {
	file, err := os.OpenFile("../.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Printf("error opening log file: %v\n", err)
		return
	}
	defer file.Close()
	logger := log.New(file, "", log.LstdFlags|log.Lmicroseconds)
	logger.Println(e)
}

func (m *LRUMap) Print() string {
	m.PrintNodes()
	var builder strings.Builder
	for _, node := range m.nodes {
		builder.WriteString(fmt.Sprintf("Key: %d, Value: %v\n",
			node.key, node.value))
	}
	return builder.String()
}

func (m *LRUMap) PrintNodes() {
	fmt.Println(m.head, m.tail)
	for _, node := range m.nodes {
		fmt.Println(node)
	}
}
