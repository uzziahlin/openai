package openai

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http"
	"testing"
	"time"
)

func TestFineTuneServiceOp_Create(t *testing.T) {
	server := newMockServer(newMockHandler(t, "POST", "fine_tune_create_response.json"))
	client := newMockClient(server.URL)
	defer server.Close()

	cancels := make([]context.CancelFunc, 0)
	defer func() {
		for _, cancel := range cancels {
			cancel()
		}
	}()

	mockReq := &FineTuneCreateRequest{
		TrainingFile: "file-XGinujblHPwGLSztz8cPS8XY",
		Model:        "davinci",
	}

	testCase := []struct {
		name    string
		ctx     context.Context
		wantRes *FineTune
		wantErr error
	}{
		{
			name: "test fine-tune create success",
			ctx:  context.TODO(),
			wantRes: func() *FineTune {
				var wantRes FineTune
				loadMockData("fine_tune_create_response.json", &wantRes)
				return &wantRes
			}(),
		},
		{
			name: "test fine-tune create timeout",
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
			res, err := client.FineTunes.Create(tc.ctx, mockReq)
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

func TestFineTuneServiceOp_List(t *testing.T) {
	server := newMockServer(newMockHandler(t, "GET", "fine_tune_list_response.json"))
	client := newMockClient(server.URL)
	defer server.Close()

	cancels := make([]context.CancelFunc, 0)
	defer func() {
		for _, cancel := range cancels {
			cancel()
		}
	}()

	testCase := []struct {
		name    string
		ctx     context.Context
		wantRes *FineTuneListResponse
		wantErr error
	}{
		{
			name: "test fine-tune list success",
			ctx:  context.TODO(),
			wantRes: func() *FineTuneListResponse {
				var wantRes FineTuneListResponse
				loadMockData("fine_tune_list_response.json", &wantRes)
				return &wantRes
			}(),
		},
		{
			name: "test fine-tune list timeout",
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
			res, err := client.FineTunes.List(tc.ctx)
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

func TestFineTuneServiceOp_Retrieve(t *testing.T) {
	server := newMockServer(newMockHandler(t, "GET", "fine_tune_retrieve_response.json"))
	client := newMockClient(server.URL)
	defer server.Close()

	cancels := make([]context.CancelFunc, 0)
	defer func() {
		for _, cancel := range cancels {
			cancel()
		}
	}()

	testCase := []struct {
		name    string
		ctx     context.Context
		id      string
		wantRes *FineTune
		wantErr error
	}{
		{
			name: "test fine-tune retrieve success",
			ctx:  context.TODO(),
			id:   "ft-AF1WoRqd3aJAHsqc9NY7iL8F",
			wantRes: func() *FineTune {
				var wantRes FineTune
				loadMockData("fine_tune_retrieve_response.json", &wantRes)
				return &wantRes
			}(),
		},
		{
			name: "test fine-tune retrieve timeout",
			ctx: func() context.Context {
				ctx, cancel := context.WithTimeout(context.Background(), 2)
				cancels = append(cancels, cancel)
				return ctx
			}(),
			id:      "ft-AF1WoRqd3aJAHsqc9NY7iL8F",
			wantErr: context.DeadlineExceeded,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			res, err := client.FineTunes.Retrieve(tc.ctx, tc.id)
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

func TestFineTuneServiceOp_Cancel(t *testing.T) {
	server := newMockServer(newMockHandler(t, "POST", "fine_tune_cancel_response.json"))
	client := newMockClient(server.URL)
	defer server.Close()

	cancels := make([]context.CancelFunc, 0)
	defer func() {
		for _, cancel := range cancels {
			cancel()
		}
	}()

	testCase := []struct {
		name    string
		ctx     context.Context
		id      string
		wantRes *FineTune
		wantErr error
	}{
		{
			name: "test fine-tune cancel success",
			ctx:  context.TODO(),
			id:   "ft-AF1WoRqd3aJAHsqc9NY7iL8F",
			wantRes: func() *FineTune {
				var wantRes FineTune
				loadMockData("fine_tune_cancel_response.json", &wantRes)
				return &wantRes
			}(),
		},
		{
			name: "test fine-tune cancel timeout",
			ctx: func() context.Context {
				ctx, cancel := context.WithTimeout(context.Background(), 2)
				cancels = append(cancels, cancel)
				return ctx
			}(),
			id:      "ft-AF1WoRqd3aJAHsqc9NY7iL8F",
			wantErr: context.DeadlineExceeded,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			res, err := client.FineTunes.Cancel(tc.ctx, tc.id)
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

func TestFineTuneServiceOp_ListEvents(t *testing.T) {
	server := newMockServer(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodGet, r.Method)

		query := r.URL.Query()

		stream := query.Get("stream")

		mockData := loadTestdata("fine_tune_events_response.json")

		if stream == "false" {
			// 模拟网络延迟
			time.Sleep(3 * time.Second)
			_, _ = w.Write(mockData)
			return
		}

		mockOutputWithStream(r.Context(), w, mockData, 5)
	})
	defer server.Close()

	client := newMockClient(server.URL)

	cancels := make([]context.CancelFunc, 0)
	defer func() {
		for _, cancel := range cancels {
			cancel()
		}
	}()

	testCase := []struct {
		name         string
		ctx          context.Context
		id           string
		stream       bool
		wantRes      *EventListResponse
		wantErr      error
		wantResCount int
	}{
		{
			name: "test list events not stream",
			ctx:  context.TODO(),
			id:   "ft-AF1WoRqd3aJAHsqc9NY7iL8F",
			wantRes: func() *EventListResponse {
				var wantRes EventListResponse
				loadMockData("fine_tune_events_response.json", &wantRes)
				return &wantRes
			}(),
			wantResCount: 1,
		},
		{
			name: "test list events not stream but timeout",
			ctx: func() context.Context {
				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				cancels = append(cancels, cancel)
				return ctx
			}(),
			id: "ft-AF1WoRqd3aJAHsqc9NY7iL8F",
			wantRes: func() *EventListResponse {
				var wantRes EventListResponse
				loadMockData("fine_tune_events_response.json", &wantRes)
				return &wantRes
			}(),
			wantErr:      context.DeadlineExceeded,
			wantResCount: 0,
		},
		{
			name:   "test list events stream",
			ctx:    context.TODO(),
			id:     "ft-AF1WoRqd3aJAHsqc9NY7iL8F",
			stream: true,
			wantRes: func() *EventListResponse {
				var wantRes EventListResponse
				loadMockData("fine_tune_events_response.json", &wantRes)
				return &wantRes
			}(),
			wantResCount: 5,
		},
		{
			name: "test list events stream with timeout",
			ctx: func() context.Context {
				ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
				cancels = append(cancels, cancel)
				return ctx
			}(),
			id:     "ft-AF1WoRqd3aJAHsqc9NY7iL8F",
			stream: true,
			wantRes: func() *EventListResponse {
				var wantRes EventListResponse
				loadMockData("fine_tune_events_response.json", &wantRes)
				return &wantRes
			}(),
			wantErr: context.DeadlineExceeded,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			res, err := client.FineTunes.ListEvents(tc.ctx, tc.id, tc.stream)
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

func TestFineTuneServiceOp_DeleteModel(t *testing.T) {
	server := newMockServer(newMockHandler(t, "DELETE", "model_delete_response.json"))
	client := newMockClient(server.URL)
	defer server.Close()

	cancels := make([]context.CancelFunc, 0)
	defer func() {
		for _, cancel := range cancels {
			cancel()
		}
	}()

	testCase := []struct {
		name    string
		ctx     context.Context
		id      string
		wantRes *ModelDeleteResponse
		wantErr error
	}{
		{
			name: "test model delete success",
			ctx:  context.TODO(),
			id:   "ft-AF1WoRqd3aJAHsqc9NY7iL8F",
			wantRes: func() *ModelDeleteResponse {
				var wantRes ModelDeleteResponse
				loadMockData("model_delete_response.json", &wantRes)
				return &wantRes
			}(),
		},
		{
			name: "test model delete timeout",
			ctx: func() context.Context {
				ctx, cancel := context.WithTimeout(context.Background(), 2)
				cancels = append(cancels, cancel)
				return ctx
			}(),
			id:      "ft-AF1WoRqd3aJAHsqc9NY7iL8F",
			wantErr: context.DeadlineExceeded,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			res, err := client.FineTunes.DeleteModel(tc.ctx, tc.id)
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
