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
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEditServiceOp_Create(t *testing.T) {
	server := newMockServer(newMockHandler(t, "POST", "edit_create_response.json"))
	client := newMockClient(server.URL)
	defer server.Close()

	cancels := make([]context.CancelFunc, 0)
	defer func() {
		for _, cancel := range cancels {
			cancel()
		}
	}()

	mockReq := &EditCreateRequest{
		Model:       "text-davinci-edit-001",
		Input:       "What day of the wek is it?",
		Instruction: "Fix the spelling mistakes",
	}

	testCase := []struct {
		name    string
		ctx     context.Context
		req     *EditCreateRequest
		wantRes *EditCreateResponse
		wantErr error
	}{
		{
			name: "test edit create success",
			ctx:  context.TODO(),
			req:  mockReq,
			wantRes: func() *EditCreateResponse {
				var wantRes EditCreateResponse
				loadMockData("edit_create_response.json", &wantRes)
				return &wantRes
			}(),
		},
		{
			name: "test edit create timeout",
			ctx: func() context.Context {
				ctx, cancel := context.WithTimeout(context.Background(), 2)
				cancels = append(cancels, cancel)
				return ctx
			}(),
			req:     mockReq,
			wantErr: context.DeadlineExceeded,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			res, err := client.Edits.Create(tc.ctx, tc.req)
			b := assert.ErrorIs(t, err, tc.wantErr)
			if !b {
				t.Fatalf("wantErr: %v, gotErr: %v", tc.wantErr, err)
			}
			if err != nil {
				return
			}
			require.Equal(t, tc.wantRes, res)
		})
	}
}
