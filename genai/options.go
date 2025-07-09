package genai

const (
	TypeText  = "text"
	TypeImage = "image"
	TypeAudio = "audio"
)

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

func WithAPIKey(apiKey string) Option {
	return func(o *Options) {
		o.APIKey = apiKey
	}
}

func WithEndpoint(endpoint string) Option {
	return func(o *Options) {
		o.Endpoint = endpoint
	}
}

func WithModel(model string) Option {
	return func(o *Options) {
		o.Model = model
	}
}

func Text(o *Options)  { o.Type = TypeText }
func Image(o *Options) { o.Type = TypeImage }
func Audio(o *Options) { o.Type = TypeAudio }
