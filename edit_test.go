package openai

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestEditServiceOp_Create(t *testing.T) {
	server := newMockServer(newMockHandler(t, "POST", "edit_create.json"))
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
				loadMockData("edit_create.json", &wantRes)
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
