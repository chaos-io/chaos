package openai

import (
	"os"

	"github.com/chaos-io/chaos/genai"
)

func init() {
	genai.Register("openai", New)
}

type openAI struct {
	options genai.Options
}

func New(opts ...genai.Option) genai.GenAI {
	var options genai.Options
	for _, o := range opts {
		o(&options)
	}

	if options.APIKey == "" {
		options.APIKey = os.Getenv("OPENAI_API_KEY")
	}

	return &openAI{options: options}
}

func (o *openAI) Generate(prompt string, opts ...genai.Option) (*genai.Result, error) {
	return nil, nil
}

func (o *openAI) Stream(prompt string, opts ...genai.Option) (*genai.Stream, error) {
	return nil, nil
}
