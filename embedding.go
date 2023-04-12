// Copyright 2023 Ken Lin
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package openai

import "context"

const (
	EmbeddingCreatePath = "/embeddings"
)

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
	var resp EmbeddingCreateResponse
	err := e.client.Post(ctx, EmbeddingCreatePath, req, &resp)
	return &resp, err
}
