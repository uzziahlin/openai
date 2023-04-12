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

func TestModerationServiceOp_Create(t *testing.T) {
	server := newMockServer(newMockHandler(t, "POST", "moderation_create_response.json"))
	client := newMockClient(server.URL)
	defer server.Close()

	cancels := make([]context.CancelFunc, 0)
	defer func() {
		for _, cancel := range cancels {
			cancel()
		}
	}()

	mockReq := &ModerationCreateRequest{
		Input: "test input",
		Model: "text-moderation-stable",
	}

	testCase := []struct {
		name    string
		ctx     context.Context
		wantRes *ModerationCreateResponse
		wantErr error
	}{
		{
			name: "test moderation create success",
			ctx:  context.TODO(),
			wantRes: func() *ModerationCreateResponse {
				var wantRes ModerationCreateResponse
				loadMockData("moderation_create_response.json", &wantRes)
				return &wantRes
			}(),
		},
		{
			name: "test moderation create timeout",
			ctx: func() context.Context {
				ctx, cancel := context.WithTimeout(context.Background(), 2)
				cancels = append(cancels, cancel)
				return ctx
			}(),
			wantErr: context.DeadlineExceeded,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			res, err := client.Moderations.Create(tc.ctx, mockReq)
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
