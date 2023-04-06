package openai

import "context"

const (
	CompletionsCreatePath = "/completions"
)

type CompletionService interface {
	Create(ctx context.Context, req *CompletionCreateRequest) (*CompletionCreateResponse, error)
}

type CompletionCreateRequest struct {
	Model            string           `json:"model"`
	Prompt           string           `json:"prompt,omitempty"`
	Suffix           string           `json:"suffix,omitempty"`
	MaxTokens        int64            `json:"max_tokens,omitempty"`
	Temperature      float64          `json:"temperature,omitempty"`
	TopP             float64          `json:"top_p,omitempty"`
	N                int64            `json:"n,omitempty"`
	Stream           bool             `json:"stream,omitempty"`
	Logprobs         int64            `json:"logprobs,omitempty"`
	Echo             bool             `json:"echo,omitempty"`
	Stop             []string         `json:"stop,omitempty"`
	PresencePenalty  float64          `json:"presence_penalty,omitempty"`
	FrequencyPenalty float64          `json:"frequency_penalty,omitempty"`
	BestOf           int64            `json:"best_of,omitempty"`
	LogitBias        map[string]int64 `json:"logit_bias,omitempty"`
	User             string           `json:"user,omitempty"`
}

type CompletionCreateResponse struct {
	Id      string        `json:"id"`
	Object  string        `json:"object"`
	Created int64         `json:"created"`
	Model   string        `json:"model"`
	Choices []*Completion `json:"choices"`
	Usage   Usage         `json:"usage"`
}

type Completion struct {
	Text         string `json:"text"`
	Delta        *Delta `json:"delta"`
	Index        int64  `json:"index"`
	Logprobs     int64  `json:"logprobs"`
	FinishReason string `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int64 `json:"prompt_tokens"`
	CompletionTokens int64 `json:"completion_tokens"`
	TotalTokens      int64 `json:"total_tokens"`
}

type CompletionServiceOp struct {
	client *Client
}

func (c CompletionServiceOp) Create(ctx context.Context, req *CompletionCreateRequest) (*CompletionCreateResponse, error) {
	var resp CompletionCreateResponse
	err := c.client.Post(ctx, CompletionsCreatePath, req, &resp)
	return &resp, err
}
