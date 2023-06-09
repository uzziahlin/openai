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
	"encoding/json"
)

const (
	CompletionsCreatePath = "/completions"
)

type CompletionService interface {
	Create(ctx context.Context, req *CompletionCreateRequest) (chan *CompletionCreateResponse, error)
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

// Create 创建一个新的聊天，为了兼容 stream 模式，返回一个 channel，如果不是 stream 模式，返回的 channel 会在第一次返回后关闭
// 如果是 stream 模式，返回的 channel 会在 ctx.Done() 或者 stream 关闭后关闭
// 这里其实也可以考虑拆分为两个方法，一个是 Create，一个是 CreateStream，但是这样会导致 API 不一致，所以这里就不拆分了
func (c CompletionServiceOp) Create(ctx context.Context, req *CompletionCreateRequest) (chan *CompletionCreateResponse, error) {

	// 如果不是 stream 模式，返回一个 channel，并将结果通过 channel 返回
	if !req.Stream {
		var resp CompletionCreateResponse
		err := c.client.Post(ctx, CompletionsCreatePath, req, &resp)
		if err != nil {
			return nil, err
		}
		res := make(chan *CompletionCreateResponse, 1)
		res <- &resp
		close(res)
		return res, nil
	}

	// 如果是 stream 模式，返回一个 channel，这个 channel 会在 ctx.Done() 或者 stream 关闭后关闭
	es, err := c.client.PostByStream(ctx, CompletionsCreatePath, req)

	if err != nil {
		return nil, err
	}

	res := make(chan *CompletionCreateResponse)

	go func() {
		defer close(res)

		for {
			select {
			case <-ctx.Done():
				return
			case e, ok := <-es:
				if !ok {
					return
				}
				var resp CompletionCreateResponse
				err := json.Unmarshal([]byte(e.Data), &resp)
				if err != nil {
					c.client.logger.Error(err, "failed to unmarshal chat response")
					continue
				}
				select {
				case <-ctx.Done():
					return
				case res <- &resp:
				}
			}

		}

	}()

	return res, nil
}
