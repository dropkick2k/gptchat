package util

import (
	"github.com/google/uuid"
)

// GenerateUUID generates a new UUID string
func GenerateUUID() string {
	return uuid.New().String()
}

// TripleQuote is a helper for triple quotes in strings
const TripleQuote = `"""`
