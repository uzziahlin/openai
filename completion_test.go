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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"testing"
	"time"
)

func TestCompletionServiceOp_Create(t *testing.T) {
	server := newMockServer(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		all, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		var req ChatCreateRequest
		err = json.Unmarshal(all, &req)
		require.NoError(t, err)

		mockData := loadTestdata("completion_create_response.json")

		if !req.Stream {
			// 模拟网络延迟
			time.Sleep(3 * time.Second)
			_, _ = w.Write(mockData)
			return
		}

		mockOutputWithStream(r.Context(), w, mockData, 5)
	})
	defer server.Close()

	client := newMockClient(server.URL)

	mockReq := CompletionCreateRequest{
		Model:       "text-davinci-003",
		Prompt:      "Say this is a test",
		MaxTokens:   7,
		Temperature: 0,
		TopP:        1,
		N:           1,
		Stop:        []string{"\n"},
	}

	cancels := make([]context.CancelFunc, 0)
	defer func() {
		for _, cancel := range cancels {
			cancel()
		}
	}()

	testCase := []struct {
		name         string
		ctx          context.Context
		req          *CompletionCreateRequest
		wantRes      *CompletionCreateResponse
		wantErr      error
		wantResCount int
	}{
		{
			name: "test completion create not stream",
			ctx:  context.TODO(),
			req: func() *CompletionCreateRequest {
				r := mockReq
				r.Stream = false
				return &r
			}(),
			wantRes: func() *CompletionCreateResponse {
				var wantRes CompletionCreateResponse
				loadMockData("completion_create_response.json", &wantRes)
				return &wantRes
			}(),
			wantResCount: 1,
		},
		{
			name: "test completion create not stream but timeout",
			ctx: func() context.Context {
				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				cancels = append(cancels, cancel)
				return ctx
			}(),
			req: func() *CompletionCreateRequest {
				r := mockReq
				r.Stream = false
				return &r
			}(),
			wantRes: func() *CompletionCreateResponse {
				var wantRes CompletionCreateResponse
				loadMockData("completion_create_response.json", &wantRes)
				return &wantRes
			}(),
			wantErr:      context.DeadlineExceeded,
			wantResCount: 0,
		},
		{
			name: "test completion create stream",
			ctx:  context.TODO(),
			req: func() *CompletionCreateRequest {
				r := mockReq
				r.Stream = true
				return &r
			}(),
			wantRes: func() *CompletionCreateResponse {
				var wantRes CompletionCreateResponse
				loadMockData("completion_create_response.json", &wantRes)
				return &wantRes
			}(),
			wantResCount: 5,
		},
		{
			name: "test completion create stream with timeout",
			ctx: func() context.Context {
				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				cancels = append(cancels, cancel)
				return ctx
			}(),
			req: func() *CompletionCreateRequest {
				r := mockReq
				r.Stream = true
				return &r
			}(),
			wantRes: func() *CompletionCreateResponse {
				var wantRes CompletionCreateResponse
				loadMockData("chat_completion_create_response.json", &wantRes)
				return &wantRes
			}(),
			wantErr: context.DeadlineExceeded,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			res, err := client.Completions.Create(tc.ctx, tc.req)
			b := assert.ErrorIs(t, err, tc.wantErr)
			if !b {
				t.Fatalf("wantErr: %v, gotErr: %v", tc.wantErr, err)
			}
			if err != nil {
				return
			}
			count := 0
		LOOP:
			for {
				select {
				case <-tc.ctx.Done():
					break LOOP
				case r, ok := <-res:
					if !ok {
						break LOOP
					}
					require.Equal(t, tc.wantRes, r)
					count++
				}
			}
			require.Equal(t, tc.wantResCount, count)

		})
	}
}
