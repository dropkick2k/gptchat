package memory

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"time"

	"github.com/ian-kent/gptchat/config"
	"github.com/ian-kent/gptchat/util"
	openai "github.com/sashabaranov/go-openai"
)

// Memory struct with additional metadata
type Memory struct {
	ID           string    `json:"id"`            // Unique identifier for each memory
	DateStored   string    `json:"date_stored"`   // Timestamp when the memory was stored
	Memory       string    `json:"memory"`        // The actual memory content
	Context      string    `json:"context"`       // Contextual information
	Tags         []string  `json:"tags"`          // Tags for categorization
	Priority     int       `json:"priority"`      // Priority level for prioritization
	LastAccessed time.Time `json:"last_accessed"` // Timestamp of the last access
}

// Module struct with mutex for concurrency safety
type Module struct {
	cfg      config.Config
	client   *openai.Client
	memories []Memory
	mu       sync.Mutex
}

// ID returns the module identifier
func (m *Module) ID() string {
	return "memory"
}

// Load initializes the module
func (m *Module) Load(cfg config.Config, client *openai.Client) error {
	m.cfg = cfg
	m.client = client
	return m.loadFromFile()
}

// UpdateConfig updates the module's configuration
func (m *Module) UpdateConfig(cfg config.Config) {
	m.cfg = cfg
}

// Execute handles memory commands
func (m *Module) Execute(args, body string) (string, error) {
	switch args {
	case "store":
		return m.Store(body)
	case "recall":
		return m.Recall(body)
	case "review":
		return m.Review(body)
	case "delete":
		return m.Delete(body)
	default:
		return "", fmt.Errorf("command not implemented: /memory %s", args)
	}
}

// Prompt returns the memory prompt
func (m *Module) Prompt() string {
	return memoryPrompt
}

const memoryPrompt = `You have an advanced long-term memory system with capabilities for selective memory, memory tagging, context prioritization, memory review, and deletion.

### Memory Commands:

1. **Store Memory**
   - **Command:** `/memory store { ... }`
   - **Description:** Store a new memory with optional tags and priority.
   - **Example:**
     ```
     /memory store {
         "memory": "I bought cookies yesterday",
         "context": "The user was discussing what they'd eaten",
         "tags": ["personal", "food"],
         "priority": 2
     }
     ```

2. **Recall Memory**
   - **Command:** `/memory recall { ... }`
   - **Description:** Recall memories based on a query and optional tags.
   - **Example:**
     ```
     /memory recall {
         "query": "When did I buy cookies?",
         "tags": ["food"]
     }
     ```

3. **Review Memories**
   - **Command:** `/memory review { ... }`
   - **Description:** Review all stored memories, optionally filtered by tags.
   - **Example:**
     ```
     /memory review {
         "tags": ["personal"]
     }
     ```

4. **Delete Memory**
   - **Command:** `/memory delete { ... }`
   - **Description:** Delete a specific memory by its ID.
   - **Example:**
     ```
     /memory delete {
         "id": "123e4567-e89b-12d3-a456-426614174000"
     }
     ```

### Usage Guidelines:

- **Storing Memories:** Always include useful context and relevant tags to enhance future recall and analysis.
- **Recalling Memories:** Use specific queries and, if necessary, filter by tags to retrieve the most relevant memories.
- **Reviewing and Managing Memories:** Regularly review your memories to maintain relevancy and delete any unnecessary or sensitive information.

### Best Practices:

- **Consistent Tagging:** Use a consistent set of tags to categorize your memories effectively.
- **Prioritize Important Memories:** Assign higher priority levels to critical memories to ensure they are recalled first.
- **Regular Review:** Periodically review and delete unnecessary memories to maintain a clean memory base.
- **Secure Sensitive Information:** Avoid storing highly sensitive data unless necessary and utilize the delete command to remove such information when needed.

You can use these memory commands at any time to manage the agent's long-term memory. The commands must be an entire message, with no conversational text, and follow the specified JSON structure.

You must not remember the current date. The current date changes and is not a useful memory.
`

// loadFromFile loads memories from the JSON file
func (m *Module) loadFromFile() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	_, err := os.Stat("memories.json")
	if os.IsNotExist(err) {
		return nil
	}

	b, err := ioutil.ReadFile("memories.json")
	if err != nil {
		return err
	}

	err = json.Unmarshal(b, &m.memories)
	if err != nil {
		return err
	}

	return nil
}

// writeToFile saves memories to the JSON file
func (m *Module) writeToFile() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	b, err := json.MarshalIndent(m.memories, "", "  ")
	if err != nil {
		return err
	}

	err = ioutil.WriteFile("memories.json", b, 0660)
	if err != nil {
		return err
	}
	return nil
}

// appendMemory adds a new memory and saves to file
func (m *Module) appendMemory(mem Memory) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.memories = append(m.memories, mem)
	return m.writeToFile()
}
