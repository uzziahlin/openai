package openai

import (
	"context"
)

const (
	ChatCreatePath = "/v1/chat/completions"
)

type Chat interface {
	Create(ctx context.Context, req *ChatReq) (*ChatReap, error)
}

type ChatReq struct {
	Model            string           `json:"model"`
	Messages         []*Message       `json:"messages"`
	Temperature      float64          `json:"temperature"`
	TopP             float64          `json:"top_p"`
	N                int64            `json:"n"`
	Stream           bool             `json:"stream"`
	Stop             []string         `json:"stop"`
	MaxTokens        int64            `json:"max_tokens"`
	PresencePenalty  float64          `json:"presence_penalty"`
	FrequencyPenalty float64          `json:"frequency_penalty"`
	LogitBias        map[string]int64 `json:"logit_bias"`
	User             string           `json:"user"`
}

type Message struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatReap struct {
	Id      string    `json:"id"`
	Object  string    `json:"object"`
	Created int64     `json:"created"`
	Choices []*Choice `json:"choices"`
	Usage   Usage     `json:"usage"`
}

type Choice struct {
	Index        int64    `json:"index"`
	Message      *Message `json:"message"`
	FinishReason string   `json:"finish_reason"`
}

type Usage struct {
	PromptTokens     int64 `json:"prompt_tokens"`
	CompletionTokens int64 `json:"completion_tokens"`
	TotalTokens      int64 `json:"total_tokens"`
}

type chat struct {
	client *Client
}

func (c chat) Create(ctx context.Context, req *ChatReq) (*ChatReap, error) {

}
