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

import (
	"context"
	"fmt"
)

const (
	ModelListPath     = "/models"
	ModelRetrievePath = "/models/%s"

	GPT35Turbo        = "gpt-3.5-turbo"
	GPT35Turbo0301    = "gpt-3.5-turbo-0301"
	GPT35Turbo0613    = "gpt-3.5-turbo-0613"
	GPT35Turbo16k     = "gpt-3.5-turbo-16k"
	GPT35Turbo16k0613 = "gpt-3.5-turbo-16k-0613"
	GPT4              = "gpt-4"
	GPT40314          = "gpt-4-0314"
	GPT40613          = "gpt-4-0613"

	Ada        = "ada"
	TextAda001 = "text-ada-001"

	Babbage        = "babbage"
	TextBabbage001 = "text-babbage-001"

	Curie        = "curie"
	TextCurie001 = "text-curie-001"

	DaVinci        = "davinci"
	TextDaVinci001 = "text-davinci-001"
	TextDaVinci002 = "text-davinci-002"
	TextDaVinci003 = "text-davinci-003"
)

type ModelService interface {
	List(ctx context.Context) (*ModelResponse, error)
	Retrieve(ctx context.Context, model string) (*Model, error)
}

type ModelResponse struct {
	Data   []*Model `json:"data"`
	Object string   `json:"object"`
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
	var resp ModelResponse
	err := m.client.Get(ctx, ModelListPath, nil, &resp)
	return &resp, err
}

func (m ModelServiceOp) Retrieve(ctx context.Context, model string) (*Model, error) {
	var resp Model
	err := m.client.Get(ctx, fmt.Sprintf(ModelRetrievePath, model), nil, &resp)
	return &resp, err
}
