package memory

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/ian-kent/gptchat/util"
	openai "github.com/sashabaranov/go-openai"
)

// Recall retrieves memories based on user input
func (m *Module) Recall(input string) (string, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	// Define the input structure
	var data struct {
		Query string   `json:"query"`
		Tags  []string `json:"tags,omitempty"`
	}

	// Parse the input JSON
	err := json.Unmarshal([]byte(input), &data)
	if err != nil {
		return "", fmt.Errorf("failed to parse input: %v", err)
	}

	// Prepare the memories in JSON format
	b, err := json.Marshal(m.memories)
	if err != nil {
		return "", err
	}

	resp, err := m.client.CreateChatCompletion(
		context.Background(),
		openai.ChatCompletionRequest{
			Model: m.cfg.OpenAIAPIModel(),
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleSystem,
					Content: `You are a helpful assistant.

I'll give you a list of existing memories, and a prompt which asks you to identify the memory I'm looking for.

You should review the listed memories and suggest which memories might match the request.`,
				},
				{
					Role: openai.ChatMessageRoleSystem,
					Content: `Here are your memories in JSON format:

` + util.TripleQuote + `
` + string(b) + `
` + util.TripleQuote,
				},
				{
					Role: openai.ChatMessageRoleSystem,
					Content: `Help me find any memories which may match this request:

` + util.TripleQuote + `
` + data.Query + `
` + util.TripleQuote,
				},
			},
		},
	)
	if err != nil {
		return "", err
	}

	response := resp.Choices[0].Message.Content

	return `You have successfully recalled this memory:

` + util.TripleQuote + `
` + response + `
` + util.TripleQuote, nil
}
