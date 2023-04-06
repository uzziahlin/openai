package openai

import (
	"context"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
)

func TestChatServiceOp_Create(t *testing.T) {
	server := newMockServer(func(w http.ResponseWriter, r *http.Request) {
		_, _ = w.Write(loadTestdata("chat_completion_create.json"))
	})

	client := newMockClient(server.URL)

	resp, err := client.Chat.Create(context.TODO(), &ChatCreateRequest{})

	require.NoError(t, err)

	var wantRes ChatCreateResponse

	loadMockData("chat_completion_create.json", &wantRes)

	require.Equal(t, &wantRes, resp)
}
