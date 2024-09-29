package memory

import (
	"encoding/json"
	"fmt"
)

// Delete removes a memory by its ID
func (m *Module) Delete(input string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Define the input structure
	var data struct {
		ID string `json:"id"`
	}

	// Parse the input JSON
	err := json.Unmarshal([]byte(input), &data)
	if err != nil {
		return "", fmt.Errorf("failed to parse input: %v", err)
	}

	// Find and delete the memory with the given ID
	for i, mem := range m.memories {
		if mem.ID == data.ID {
			m.memories = append(m.memories[:i], m.memories[i+1:]...)
			err = m.writeToFile()
			if err != nil {
				return "", fmt.Errorf("failed to write to file: %v", err)
			}
			return fmt.Sprintf("Memory with ID %s deleted successfully.", data.ID), nil
		}
	}

	return "", fmt.Errorf("memory with ID %s not found", data.ID)
}
