package gemini

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/chaos-io/chaos/genai"
	"github.com/chaos-io/chaos/logs"
	jsoniter "github.com/json-iterator/go"
)

const (
	defaultEndpoint = "https://generativelanguage.googleapis.com/v1beta/models/"

	// See: https://ai.google.dev/gemini-api/docs/models?hl=zh-cn
	modelGemini25ProVision = "gemini-2.5-pro-vision"
	modelGemini25Pro       = "gemini-2.5-pro"

	modelGemini20FlashLite    = "gemini-2.0-flash-lite"
	modelGemini20FlashPVImage = "gemini-2.0-flash-preview-image-generation"

	modelImagen40 = "imagen-4.0-generate-preview-06-06"
)

func init() {
	genai.Register("gemini", New)
}

// gemini implements the GenAI interface using Google Gemini 2.5 API.
// See: https://github.com/googleapis/go-genai
type gemini struct {
	options *genai.Options
}

func New(opts ...genai.Option) genai.GenAI {
	options := &genai.Options{}
	for _, o := range opts {
		o(options)
	}

	if options.APIKey == "" {
		options.APIKey = os.Getenv("GEMINI_API_KEY")
	}

	return &gemini{options: options}
}

func (g *gemini) Generate(prompt string, opts ...genai.Option) (*genai.Result, error) {
	options := g.options
	for _, opt := range opts {
		opt(options)
	}

	if options.Endpoint == "" {
		options.Endpoint = defaultEndpoint
	}

	switch options.Type {
	case genai.TypeImage:
		return g.getImageResult(prompt)
	case genai.TypeAudio:
		return g.getAudioResult(prompt)
	default:
		return g.getResult(prompt)
	}
}

func (g *gemini) Stream(prompt string, opts ...genai.Option) (*genai.Stream, error) {
	results := make(chan *genai.Result)
	go func() {
		defer close(results)
		res, err := g.Generate(prompt, opts...)
		if err != nil {
			// Send error via Stream.Err, not channel
			return
		}
		results <- res
	}()
	return &genai.Stream{Results: results}, nil
}

func (g *gemini) getImageResult(prompt string) (*genai.Result, error) {
	model := g.options.Model
	if model == "" {
		model = modelGemini20FlashPVImage
	}

	url := fmt.Sprintf("%s%s:generateContent", g.options.Endpoint, model)
	body := map[string]interface{}{
		"contents": []map[string]interface{}{
			{"parts": []map[string]string{{"text": prompt}}},
		},
		"generationConfig": map[string][]string{"responseModalities": {"TEXT", "IMAGE"}},
	}

	buf, err := g.httpDo(url, body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text       string `json:"text"`
					InlineData struct {
						Data []byte `json:"data"`
					} `json:"inlineData"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := jsoniter.Unmarshal(buf, &result); err != nil {
		return nil, err
	}
	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		logs.Warnw("failed to get the valid result", "result", string(buf))
		return nil, fmt.Errorf("no candidates returned")
	}

	// debug response log
	//_ = os.WriteFile("my.log", buf, 0644)

	var text string
	parts := result.Candidates[0].Content.Parts
	for _, part := range parts {
		if part.Text != "" {
			text = part.Text
		} else if part.InlineData.Data != nil {
			imageBytes := part.InlineData.Data
			logs.Debugw("len imageBytes", "len", len(imageBytes))
			outputFile := "gemini-generated.png"
			_ = os.WriteFile(outputFile, imageBytes, 0644)
		}
	}

	return &genai.Result{
		Prompt: prompt,
		Type:   g.options.Type,
		Data:   nil,
		Text:   text,
	}, nil
}

func (g *gemini) getAudioResult(prompt string) (*genai.Result, error) {
	model := g.options.Model
	if model == "" {
		model = modelGemini25Pro
	}

	url := fmt.Sprintf("%s%s:generateContent", g.options.Endpoint, model)
	body := map[string]interface{}{
		"contents": []map[string]interface{}{
			{"parts": []map[string]string{{"text": prompt}}},
		},
		"responseMimeType": "audio/wav",
	}

	buf, err := g.httpDo(url, body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					InlineData struct {
						Data []byte `json:"data"`
					} `json:"inline_data"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := jsoniter.Unmarshal(buf, &result); err != nil {
		return nil, err
	}
	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no audio returned")
	}

	return &genai.Result{
		Prompt: prompt,
		Type:   g.options.Type,
		Data:   result.Candidates[0].Content.Parts[0].InlineData.Data,
	}, nil
}

func (g *gemini) getResult(prompt string) (*genai.Result, error) {
	model := g.options.Model
	if model == "" {
		model = modelGemini25Pro
	}

	url := fmt.Sprintf("%s%s:generateContent", g.options.Endpoint, model)
	body := map[string]interface{}{
		"contents": []map[string]interface{}{
			{"parts": []map[string]string{{"text": prompt}}},
		},
	}

	buf, err := g.httpDo(url, body)
	if err != nil {
		return nil, err
	}

	var result struct {
		Candidates []struct {
			Content struct {
				Parts []struct {
					Text string `json:"text"`
				} `json:"parts"`
			} `json:"content"`
		} `json:"candidates"`
	}

	if err := jsoniter.Unmarshal(buf, &result); err != nil {
		return nil, err
	}
	if len(result.Candidates) == 0 || len(result.Candidates[0].Content.Parts) == 0 {
		return nil, fmt.Errorf("no candidates returned")
	}

	return &genai.Result{
		Prompt: prompt,
		Type:   g.options.Type,
		Data:   nil,
		Text:   result.Candidates[0].Content.Parts[0].Text,
	}, nil
}

func (g *gemini) httpDo(url string, body map[string]interface{}) ([]byte, error) {
	b, _ := jsoniter.Marshal(body)
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost, url, bytes.NewBuffer(b))
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-goog-api-key", g.options.APIKey)

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	return io.ReadAll(resp.Body)
}
