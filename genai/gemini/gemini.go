package gemini

import (
	"os"

	"github.com/chaos-io/chaos/genai"
)

func init() {
	genai.Register("gemini", New)
}

// gemini implements the GenAI interface using Google Gemini 2.5 API.
type gemini struct {
	options genai.Options
}

func New(opts ...genai.Option) genai.GenAI {
	var options genai.Options
	for _, o := range opts {
		o(&options)
	}

	if options.APIKey == "" {
		options.APIKey = os.Getenv("GEMINI_API_KEY")
	}

	return &gemini{options: options}
}

func (g *gemini) Generate(prompt string, opts ...genai.Option) (*genai.Result, error) {
	return nil, nil
}

func (g *gemini) Stream(prompt string, opts ...genai.Option) (*genai.Stream, error) {
	return nil, nil
}
