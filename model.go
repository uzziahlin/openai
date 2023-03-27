package openai

import "context"

type ModelService interface {
	List(ctx context.Context) (*ModelResponse, error)
	Retrieve(ctx context.Context, model string) (*Model, error)
}

type ModelResponse struct {
	Data []*Model `json:"data"`
}

type Model struct {
	Id         string   `json:"id"`
	Object     string   `json:"object"`
	OwnedBy    string   `json:"owned_by"`
	Permission []string `json:"permission"`
}

type ModelServiceOp struct {
	client *Client
}

func (m ModelServiceOp) List(ctx context.Context) (*ModelResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (m ModelServiceOp) Retrieve(ctx context.Context, model string) (*Model, error) {
	//TODO implement me
	panic("implement me")
}
