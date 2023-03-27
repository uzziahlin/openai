package openai

import "context"

type EditService interface {
	Create(ctx context.Context, request *EditCreateRequest) (*EditCreateResponse, error)
}

type EditCreateRequest struct {
	Model       string  `json:"model"`
	Input       string  `json:"input,omitempty"`
	Instruction string  `json:"instruction"`
	N           int64   `json:"n,omitempty"`
	Temperature float64 `json:"temperature,omitempty"`
	TopP        float64 `json:"top_p,omitempty"`
}

type EditCreateResponse struct {
	Object  string  `json:"object"`
	Created int64   `json:"created"`
	Choices []*Edit `json:"choices"`
	Usage   Usage   `json:"usage"`
}

type Edit struct {
	Text  string `json:"text"`
	Index int64  `json:"index"`
}

type EditServiceOp struct {
	client *Client
}

func (e EditServiceOp) Create(ctx context.Context, request *EditCreateRequest) (*EditCreateResponse, error) {
	//TODO implement me
	panic("implement me")
}
