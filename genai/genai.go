// Package genai provides a generic interface for generative AI providers.
package genai

import (
	"sync"
)

var (
	providers     map[string]GenAI
	providersOnce sync.Once
)

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

func Register(name string, gen GenAI) {
	providersOnce.Do(func() {
		providers = make(map[string]GenAI)
	})
	providers[name] = gen
}

func Get(name string) GenAI {
	if gen, ok := providers[name]; ok {
		return gen
	}
	return nil
}
