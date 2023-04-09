package openai

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
	"net/http"
	"strings"
	"testing"
	"time"
)

func TestChatServiceOp_Create(t *testing.T) {
	server := newMockServer(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)

		all, err := io.ReadAll(r.Body)
		require.NoError(t, err)

		var req ChatCreateRequest
		err = json.Unmarshal(all, &req)
		require.NoError(t, err)

		mockData := loadTestdata("chat_completion_create.json")

		if !req.Stream {
			// 模拟网络延迟
			time.Sleep(3 * time.Second)
			_, _ = w.Write(mockData)
			return
		}

		mockOutputWithStream(r.Context(), w, mockData, 5)
	})

	client := newMockClient(server.URL)

	mockReq := ChatCreateRequest{
		Model: "gpt-3.5-turbo",
		Messages: []*Message{
			{Role: "system", Content: "You are a helpful assistant."},
			{Role: "user", Content: "Who won the world series in 2020?"},
			{Role: "assistant", Content: "The Los Angeles Dodgers won the World Series in 2020."},
			{Role: "user", Content: "Where was it played?"},
		},
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
		req          *ChatCreateRequest
		wantRes      *ChatCreateResponse
		wantErr      error
		wantResCount int
	}{
		{
			name: "test chat create not stream",
			ctx:  context.TODO(),
			req: func() *ChatCreateRequest {
				r := mockReq
				r.Stream = false
				return &r
			}(),
			wantRes: func() *ChatCreateResponse {
				var wantRes ChatCreateResponse
				loadMockData("chat_completion_create.json", &wantRes)
				return &wantRes
			}(),
			wantResCount: 1,
		},
		{
			name: "test chat create not stream but timeout",
			ctx: func() context.Context {
				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				cancels = append(cancels, cancel)
				return ctx
			}(),
			req: func() *ChatCreateRequest {
				r := mockReq
				r.Stream = false
				return &r
			}(),
			wantRes: func() *ChatCreateResponse {
				var wantRes ChatCreateResponse
				loadMockData("chat_completion_create.json", &wantRes)
				return &wantRes
			}(),
			wantErr:      context.DeadlineExceeded,
			wantResCount: 0,
		},
		{
			name: "test chat create stream",
			ctx:  context.TODO(),
			req: func() *ChatCreateRequest {
				r := mockReq
				r.Stream = true
				return &r
			}(),
			wantRes: func() *ChatCreateResponse {
				var wantRes ChatCreateResponse
				loadMockData("chat_completion_create.json", &wantRes)
				return &wantRes
			}(),
			wantResCount: 5,
		},
		{
			name: "test chat create stream with timeout",
			ctx: func() context.Context {
				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				cancels = append(cancels, cancel)
				return ctx
			}(),
			req: func() *ChatCreateRequest {
				r := mockReq
				r.Stream = true
				return &r
			}(),
			wantRes: func() *ChatCreateResponse {
				var wantRes ChatCreateResponse
				loadMockData("chat_completion_create.json", &wantRes)
				return &wantRes
			}(),
			wantErr: context.DeadlineExceeded,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			res, err := client.Chat.Create(tc.ctx, tc.req)
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

func mockOutputWithStream(ctx context.Context, w http.ResponseWriter, data []byte, count int) {

	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("Content-Length", "16384")

	w.WriteHeader(http.StatusOK)

	content := strings.ReplaceAll(string(data), "\n", "")

	for i := 0; i < count; i++ {
		select {
		case <-ctx.Done():
			fmt.Println("context done")
			return
		default:
			_, err := fmt.Fprintf(w, "data:%s\n\n", content)
			if err != nil {
				panic(err)
			}
			time.Sleep(600 * time.Millisecond)
		}
	}

	_, _ = w.Write([]byte("data:"))
	_, _ = w.Write([]byte("[DONE]"))
}

func TestSend(t *testing.T) {
	client, err := New(App{
		ApiUrl: "https://api.openai.com",
		ApiKey: "sk-xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx",
	}, WithProxy(&Proxy{
		Url:      "xxx.xxx.xxx.xxx:xxxx",
		Username: "xxxxxx",
		Password: "xxxxxx",
	}))

	if err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	req := ChatCreateRequest{
		Model: "gpt-3.5-turbo",
		Messages: []*Message{
			{Role: "user", Content: "Where was it played?"},
		},
		Stream: true,
	}

	res, err := client.Chat.Create(ctx, &req)

	if err != nil {
		t.Fatal(err)
	}

	for {
		select {
		case <-ctx.Done():
			return
		case r, ok := <-res:
			if !ok {
				return
			}
			fmt.Println(r)
		}
	}
}
