package openai

import "context"

type EmbeddingService interface {
	Create(ctx context.Context, req *EmbeddingCreateRequest) (*EmbeddingCreateResponse, error)
}

type EmbeddingCreateRequest struct {
	Model string   `json:"model"`
	Input []string `json:"input"`
	User  string   `json:"user,omitempty"`
}

type EmbeddingCreateResponse struct {
	Object string          `json:"object"`
	Data   []*Embedding    `json:"data"`
	Model  string          `json:"model"`
	Usage  *EmbeddingUsage `json:"usage"`
}

type Embedding struct {
	Object    string    `json:"object"`
	Embedding []float64 `json:"embedding"`
	Index     int64     `json:"index"`
}

type EmbeddingUsage struct {
	PromptTokens int64 `json:"prompt_tokens"`
	TotalTokens  int64 `json:"total_tokens"`
}

type EmbeddingServiceOp struct {
	client *Client
}

func (e EmbeddingServiceOp) Create(ctx context.Context, req *EmbeddingCreateRequest) (*EmbeddingCreateResponse, error) {
	//TODO implement me
	panic("implement me")
}
