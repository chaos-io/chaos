package genai

// Option is functional option for configuring providers.
type Option func(*Options)

// Options holds configuration for providers.
type Options struct {
	APIKey   string
	Endpoint string
	Type     string // text, image, audio, etc.
	Model    string // model name, e.g. gemini-2.5-pro
	// Add more fields as needed
}
