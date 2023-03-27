package openai

import (
	"context"
)

const (
	ChatCreatePath = "/v1/chat/completions"
)

type ChatService interface {
	Create(ctx context.Context, req *ChatCreateRequest) (*ChatCreateResponse, error)
}

type ChatCreateRequest struct {
	Model            string           `json:"model"`
	Messages         []*Message       `json:"messages,omitempty"`
	Temperature      float64          `json:"temperature,omitempty"`
	TopP             float64          `json:"top_p,omitempty"`
	N                int64            `json:"n,omitempty"`
	Stream           bool             `json:"stream,omitempty"`
	Stop             []string         `json:"stop,omitempty"`
	MaxTokens        int64            `json:"max_tokens,omitempty"`
	PresencePenalty  float64          `json:"presence_penalty,omitempty"`
	FrequencyPenalty float64          `json:"frequency_penalty,omitempty"`
	LogitBias        map[string]int64 `json:"logit_bias,omitempty"`
	User             string           `json:"user,omitempty"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatCreateResponse struct {
	Id      string            `json:"id"`
	Object  string            `json:"object"`
	Created int64             `json:"created"`
	Choices []*ChatCompletion `json:"choices"`
	Usage   Usage             `json:"usage"`
}

type ChatCompletion struct {
	Index        int64    `json:"index"`
	Message      *Message `json:"message"`
	FinishReason string   `json:"finish_reason"`
}

type ChatServiceOp struct {
	client *Client
}

func (c ChatServiceOp) Create(ctx context.Context, req *ChatCreateRequest) (*ChatCreateResponse, error) {
	var resp ChatCreateResponse
	err := c.client.Post(ctx, ChatCreatePath, req, &resp)
	return &resp, err
}
