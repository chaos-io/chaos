package genai

type DummyGenAI struct{}

func (g *DummyGenAI) Generate(prompt string, opts ...Option) (*Result, error) {
	return &Result{Prompt: prompt, Type: "dummy", Text: "dummy response"}, nil
}

func (g *DummyGenAI) Stream(prompt string, opts ...Option) (*Stream, error) {
	results := make(chan *Result, 1)
	results <- &Result{Prompt: prompt, Type: "dummy", Text: "dummy response"}
	close(results)
	return &Stream{Results: results}, nil
}
