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
