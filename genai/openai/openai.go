package openai

import (
	"bytes"
	"fmt"
	"io"
	"net/http"
	"os"

	jsoniter "github.com/json-iterator/go"

	"github.com/chaos-io/chaos/genai"
)

const (
	modelDALLE3     = "dall-e-3"
	modelTTS1       = "tts-1"
	modelGPT35Turbo = "gpt-3.5-turbo"

	urlImage   = "https://api.openai.com/v1/images/generations"
	urlAudio   = "https://api.openai.com/v1/audio/speech"
	urlDefault = "https://api.openai.com/v1/chat/completions"
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
	options := o.options
	for _, opt := range opts {
		opt(&options)
	}

	switch options.Type {
	case genai.TypeImage:
		return o.getImageResult(prompt, options)
	case genai.TypeAudio:
		return o.getAudioResult(prompt, options)
	default:
		return o.getResult(prompt, options)
	}
}

func (o *openAI) Stream(prompt string, opts ...genai.Option) (*genai.Stream, error) {
	results := make(chan *genai.Result)
	go func() {
		defer close(results)
		res, err := o.Generate(prompt, opts...)
		if err != nil {
			// send error via Stream.Err, not channel
			return
		}
		results <- res
	}()
	return &genai.Stream{Results: results}, nil
}

func (o *openAI) getImageResult(prompt string, options genai.Options) (*genai.Result, error) {
	model := options.Model
	if model == "" {
		model = modelDALLE3
	}
	url := urlImage
	body := map[string]interface{}{
		"prompt": prompt,
		"n":      1,
		"size":   "1024x1024",
		"model":  model,
	}

	buf, err := o.httpDo(url, body, options)
	if err != nil {
		return nil, err
	}

	var result struct {
		Data []struct {
			URL string `json:"url"`
		} `json:"data"`
	}

	if err := jsoniter.Unmarshal(buf, &result); err != nil {
		return nil, err
	}
	if len(result.Data) == 0 {
		return nil, fmt.Errorf("no image returned")
	}

	return &genai.Result{
		Prompt: prompt,
		Type:   options.Type,
		Data:   nil,
		Text:   result.Data[0].URL,
	}, nil
}

func (o *openAI) getAudioResult(prompt string, options genai.Options) (*genai.Result, error) {
	model := options.Model
	if model == "" {
		model = modelTTS1
	}
	url := urlAudio
	body := map[string]interface{}{
		"model": model,
		"input": prompt,
		"voice": "alloy", // or another supported voice
	}

	buf, err := o.httpDo(url, body, options)
	if err != nil {
		return nil, err
	}

	return &genai.Result{
		Prompt: prompt,
		Type:   options.Type,
		Data:   buf,
	}, nil
}

func (o *openAI) getResult(prompt string, options genai.Options) (*genai.Result, error) {
	model := options.Model
	if model == "" {
		model = modelGPT35Turbo
	}
	url := urlDefault
	body := map[string]interface{}{
		"model":    model,
		"messages": []map[string]string{{"role": "user", "content": prompt}},
	}

	buf, err := o.httpDo(url, body, options)
	if err != nil {
		return nil, err
	}

	var result struct {
		Choices []struct {
			Message struct {
				Content string `json:"content"`
			} `json:"message"`
		} `json:"choices"`
	}
	if err := jsoniter.Unmarshal(buf, &result); err != nil {
		return nil, err
	}
	if len(result.Choices) == 0 {
		return nil, fmt.Errorf("no choices returned")
	}

	return &genai.Result{
		Prompt: prompt,
		Type:   options.Type,
		Data:   nil,
		Text:   result.Choices[0].Message.Content,
	}, nil
}

func (o *openAI) httpDo(url string, body map[string]interface{}, options genai.Options) ([]byte, error) {
	b, _ := jsoniter.Marshal(body)
	req, err := http.NewRequest(http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", options.APIKey))

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer func() {
		_ = resp.Body.Close()
	}()

	return io.ReadAll(resp.Body)
}
