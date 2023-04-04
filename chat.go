package openai

import (
	"context"
	"encoding/json"
)

const (
	ChatCreatePath = "/v1/chat/completions"
)

type ChatService interface {
	Create(ctx context.Context, req *ChatCreateRequest) (chan *ChatCreateResponse, error)
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
	Delta        *Delta   `json:"delta"`
	Message      *Message `json:"message"`
	FinishReason string   `json:"finish_reason"`
}

type Delta struct {
	Role    string `json:"role"`
	Content string `json:"content"`
}

type ChatServiceOp struct {
	client *Client
}

func (c ChatServiceOp) Create(ctx context.Context, req *ChatCreateRequest) (chan *ChatCreateResponse, error) {
	res := make(chan *ChatCreateResponse)

	if !req.Stream {
		var resp ChatCreateResponse
		err := c.client.Post(ctx, ChatCreatePath, req, &resp)
		if err != nil {
			return nil, err
		}
		go func() {
			select {
			case <-ctx.Done():
			case res <- &resp:
			}
			close(res)
		}()
		return res, nil
	}

	es, err := c.client.Stream(ctx, ChatCreatePath, req)

	if err != nil {
		return nil, err
	}

	go func() {
		for e := range es {
			var resp ChatCreateResponse
			err := json.Unmarshal([]byte(e.Data), &resp)
			if err != nil {
				c.client.logger.Error(err, "failed to unmarshal chat response")
				continue
			}
			select {
			case <-ctx.Done():
				break
			case res <- &resp:
			}
		}
		close(res)
	}()

	return res, nil
}
