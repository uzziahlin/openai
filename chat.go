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
	// ChatCreatePath 聊天创建路径
	ChatCreatePath = "/chat/completions"

	FunctionCallNone = FunctionCallString("none")
	FunctionCallAuto = FunctionCallString("auto")

	FinishReasonFunctionCall = "function_call"
	FinishReasonStop         = "stop"
)

type ChatService interface {
	Create(ctx context.Context, req *ChatCreateRequest) (chan *ChatCreateResponse, error)
}

type ChatCreateRequest struct {
	Model            string           `json:"model"`
	Messages         []*Message       `json:"messages,omitempty"`
	Functions        []*Function      `json:"functions,omitempty"`
	FunctionCall     IFunctionCall    `json:"function_call,omitempty"`
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

type IFunctionCall interface {
	Call()
}

type Function struct {
	Name        string     `json:"name"`
	Description string     `json:"description,omitempty"`
	Parameters  *Parameter `json:"parameters,omitempty"`
}

type Parameter struct {
	Type       string               `json:"type"` // object only
	Properties map[string]*Property `json:"properties,omitempty"`
	Required   []string             `json:"required,omitempty"`
}

type Property struct {
	Type        string               `json:"type"`
	Description string               `json:"description,omitempty"`
	Properties  map[string]*Property `json:"properties,omitempty"` // if type is object, this field will be set
}

type Message struct {
	Role         string       `json:"role"` // user,assistant,system,function
	Content      string       `json:"content,omitempty"`
	Name         string       `json:"name,omitempty"`
	FunctionCall FunctionCall `json:"function_call,omitempty"`
}

type FunctionCall struct {
	Name      string `json:"name"`
	Arguments string `json:"arguments,omitempty"`
}

func (f FunctionCall) Call() {
	//TODO implement me
	panic("implement me")
}

type FunctionCallString string

func (f FunctionCallString) Call() {
	//TODO implement me
	panic("implement me")
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

// Create 创建一个新的聊天，为了兼容 stream 模式，返回一个 channel，如果不是 stream 模式，返回的 channel 会在第一次返回后关闭
// 如果是 stream 模式，返回的 channel 会在 ctx.Done() 或者 stream 关闭后关闭
// 这里其实也可以考虑拆分为两个方法，一个是 Create，一个是 CreateStream，但是这样会导致 API 不一致，所以这里就不拆分了
func (c ChatServiceOp) Create(ctx context.Context, req *ChatCreateRequest) (chan *ChatCreateResponse, error) {

	// 如果不是 stream 模式，返回一个 channel，并将结果通过 channel 返回
	if !req.Stream {
		var resp ChatCreateResponse
		err := c.client.Post(ctx, ChatCreatePath, req, &resp)
		if err != nil {
			return nil, err
		}
		res := make(chan *ChatCreateResponse, 1)
		res <- &resp
		close(res)
		return res, nil
	}

	// 如果是 stream 模式，返回一个 channel，这个 channel 会在 ctx.Done() 或者 stream 关闭后关闭
	es, err := c.client.PostByStream(ctx, ChatCreatePath, req)

	if err != nil {
		return nil, err
	}

	res := make(chan *ChatCreateResponse)

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
				var resp ChatCreateResponse
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
