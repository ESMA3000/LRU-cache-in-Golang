package src

import (
	"bytes"
	"encoding/gob"
	"fmt"
	"os"
)

func (m *LRUMap) SaveToDisk(filePath string) error {
	nodes := m.Iterator()
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(nodes); err != nil {
		return fmt.Errorf("error marshaling binary data: %v", err)
	}

	return os.WriteFile(filePath, buf.Bytes(), 0644)
}

func (m *LRUMap) LoadFromDisk(filePath string) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return fmt.Errorf("error reading file: %v", err)
	}

	var nodes []*Node
	if err := gob.NewDecoder(bytes.NewReader(content)).Decode(&nodes); err != nil {
		return fmt.Errorf("error decoding binary data: %v", err)
	}

	/* for _, node := range nodes {
		m.nodes[node.key] = m.getNodeFromPool()
		m.nodes[node.key].key = node.key
		m.nodes[node.key].value = node.value
		m.setHead(m.nodes[node.key])
	} */

	return nil
}
