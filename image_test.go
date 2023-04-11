package openai

import (
	"context"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"os"
	"testing"
)

func TestImageServiceOp_Create(t *testing.T) {
	server := newMockServer(newMockHandler(t, "POST", "image_response.json"))
	client := newMockClient(server.URL)
	defer server.Close()

	cancels := make([]context.CancelFunc, 0)
	defer func() {
		for _, cancel := range cancels {
			cancel()
		}
	}()

	mockReq := &ImageCreateRequest{
		Prompt: "A cute baby sea otter",
		ImageAttributes: ImageAttributes{
			N:    1,
			Size: "1024*1024",
		},
	}

	testCase := []struct {
		name    string
		ctx     context.Context
		req     *ImageCreateRequest
		wantRes *ImageResponse
		wantErr error
	}{
		{
			name: "test image create success",
			ctx:  context.TODO(),
			req:  mockReq,
			wantRes: func() *ImageResponse {
				var wantRes ImageResponse
				loadMockData("image_response.json", &wantRes)
				return &wantRes
			}(),
		},
		{
			name: "test image create timeout",
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
			res, err := client.Images.Create(tc.ctx, tc.req)
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

func TestImageServiceOp_Edit(t *testing.T) {
	server := newMockServer(newMockHandler(t, "POST", "image_response.json"))
	// client := newMockClient(server.URL)
	defer server.Close()

	cancels := make([]context.CancelFunc, 0)
	defer func() {
		for _, cancel := range cancels {
			cancel()
		}
	}()

	testCase := []struct {
		name    string
		client  *Client
		ctx     context.Context
		req     *ImageEditRequest
		wantRes *ImageResponse
		wantErr error
		before  func()
		after   func()
	}{
		{
			name:   "test image edit success",
			client: newMockClient(server.URL),
			ctx:    context.TODO(),
			req: &ImageEditRequest{
				Image: "testdata/image.png",
				ImageAttributes: ImageAttributes{
					N:    1,
					Size: "1024*1024",
				},
			},
			wantRes: func() *ImageResponse {
				var wantRes ImageResponse
				loadMockData("image_response.json", &wantRes)
				return &wantRes
			}(),
			before: func() {
				file, err := os.Create("testdata/image.png")
				require.NoError(t, err)
				file.Close()
			},
			after: func() {
				err := os.Remove("testdata/image.png")
				require.NoError(t, err)
			},
		},
		{
			name:   "test image edit timeout",
			client: newMockClient(server.URL),
			ctx: func() context.Context {
				ctx, cancel := context.WithTimeout(context.Background(), 2)
				cancels = append(cancels, cancel)
				return ctx
			}(),
			req: &ImageEditRequest{
				Image: "testdata/image.png",
				ImageAttributes: ImageAttributes{
					N:    1,
					Size: "1024*1024",
				},
			},
			wantErr: context.DeadlineExceeded,
			before: func() {
				file, err := os.Create("testdata/image.png")
				require.NoError(t, err)
				file.Close()
			},
			after: func() {
				err := os.Remove("testdata/image.png")
				require.NoError(t, err)
			},
		},
		{
			name:   "test image edit cannot open file",
			client: newMockClient(server.URL),
			ctx: func() context.Context {
				ctx, cancel := context.WithTimeout(context.Background(), 2)
				cancels = append(cancels, cancel)
				return ctx
			}(),
			req: &ImageEditRequest{
				Image: "testdata/image.png",
				ImageAttributes: ImageAttributes{
					N:    1,
					Size: "1024*1024",
				},
			},
			wantErr: nil,
			before: func() {

			},
			after: func() {

			},
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			tc.before()
			defer tc.after()
			res, err := tc.client.Images.Edit(tc.ctx, tc.req)
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

type mockFormBuilder struct {
	createFormFile      func(name string, filename string) error
	createFormField     func(name string, value string) error
	formDataContentType func() string
	close               func() error
}

func (m mockFormBuilder) CreateFormFile(name string, filename string) error {
	return m.createFormFile(name, filename)
}

func (m mockFormBuilder) CreateFormField(name string, value string) error {
	return m.createFormField(name, value)
}

func (m mockFormBuilder) FormDataContentType() string {
	return ""
}

func (m mockFormBuilder) Close() error {
	return nil
}
