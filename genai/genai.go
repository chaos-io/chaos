// Package genai provides a generic interface for generative AI providers.
package genai

import (
	"sync"
)

var providers map[string]func(...Option) GenAI
var providersOnce sync.Once

// GenAI is the generic interface for generative for AI providers.
type GenAI interface {
	Generate(prompt string, opts ...Option) (*Result, error)
	Stream(prompt string, opts ...Option) (*Stream, error)
}

// Result is the unified response from GenAI providers.
type Result struct {
	Prompt string
	Type   string
	Data   []byte // for audio/image binary data
	Text   string // for text or image URL
}

// Stream represents a streaming responses from a GenAI provider.
type Stream struct {
	Results <-chan *Result
	Err     error
	// You can add fields for cancellation, errors, etc. if needed
}

// Register a GenAI provider by name.
func Register(name string, gen func(...Option) GenAI) {
	providersOnce.Do(func() {
		providers = make(map[string]func(...Option) GenAI)
	})
	providers[name] = gen
}

// Get a GenAI provider by name.
func Get(name string) func(...Option) GenAI {
	if gen, ok := providers[name]; ok {
		return gen
	}
	return nil
}
