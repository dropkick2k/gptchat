package memory

import (
	"encoding/json"
	"fmt"
	"time"

	"github.com/ian-kent/gptchat/util"
)

// Store adds a new memory with metadata
func (m *Module) Store(input string) (string, error) {
	// Define the input structure
	var data struct {
		Memory   string   `json:"memory"`
		Context  string   `json:"context"`
		Tags     []string `json:"tags,omitempty"`
		Priority int      `json:"priority,omitempty"`
	}

	// Parse the input JSON
	err := json.Unmarshal([]byte(input), &data)
	if err != nil {
		return "", fmt.Errorf("failed to parse input: %v", err)
	}

	// Assign default priority if not provided
	if data.Priority == 0 {
		data.Priority = 1
	}

	// Create a new Memory instance
	newMemory := Memory{
		ID:           util.GenerateUUID(), // Ensure util.GenerateUUID() generates a unique UUID
		DateStored:   time.Now().Format(time.RFC3339),
		Memory:       data.Memory,
		Context:      data.Context,
		Tags:         data.Tags,
		Priority:     data.Priority,
		LastAccessed: time.Now(),
	}

	// Append the new memory
	err = m.appendMemory(newMemory)
	if err != nil {
		return "", fmt.Errorf("failed to store memory: %v", err)
	}

	return `You have successfully stored this memory:

` + util.TripleQuote + `
` + newMemory.Memory + `
` + util.TripleQuote, nil
}
