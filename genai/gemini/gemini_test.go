package gemini

import (
	"os"
	"testing"

	"github.com/chaos-io/chaos/genai"
)

func TestGemini_GenerateText(t *testing.T) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		t.Skip("GEMINI_API_KEY not set")
	}

	client := New(genai.WithAPIKey(apiKey))

	prompt := "Explain how AI works in a few words"
	res, err := client.Generate(prompt, genai.Text)
	if err != nil {
		t.Fatalf("failed to generate text: %v", err)
	}
	if res == nil || res.Text == "" {
		t.Fatal("expected a non-empty text result")
	}
	t.Logf("text: %s", res.Text)
}

func TestGemini_GenerateImage(t *testing.T) {
	apiKey := os.Getenv("GEMINI_API_KEY")
	if apiKey == "" {
		t.Skip("GEMINI_API_KEY not set")
	}

	client := New(genai.WithAPIKey(apiKey))
	prompt := "A apple in the hand."
	res, err := client.Generate(prompt, genai.Image)
	if err != nil {
		t.Fatalf("failed to generate image: %v", err)
	}
	if res == nil || res.Text == "" {
		t.Fatal("expected a non-empty image URL")
	}
	t.Logf("image: %s", res.Text)
}
