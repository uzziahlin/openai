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
	EditCreatePath = "/edits"
)

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

func (e EditServiceOp) Create(ctx context.Context, req *EditCreateRequest) (*EditCreateResponse, error) {
	var res EditCreateResponse
	err := e.client.Post(ctx, EditCreatePath, req, &res)
	return &res, err
}
