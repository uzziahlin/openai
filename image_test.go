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
	"fmt"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"io"
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
			name: "test image edit create form file error",
			client: newMockClient(server.URL, WithFormBuilder(func(w io.Writer) FormBuilder {
				return &mockFormBuilder{
					createFormFile: func(name string, filename string) error {
						return fmt.Errorf("%w, %s", os.ErrNotExist, "testdata/image.png")
					},
				}
			})),
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
			wantErr: os.ErrNotExist,
			before: func() {

			},
			after: func() {

			},
		},
		{
			name: "test image edit create form field error",
			client: newMockClient(server.URL, WithFormBuilder(func(w io.Writer) FormBuilder {
				return &mockFormBuilder{
					createFormFile: func(name string, filename string) error {
						return fmt.Errorf("%w, %s", os.ErrInvalid, "write form field error")
					},
				}
			})),
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
			wantErr: os.ErrInvalid,
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

func TestImageServiceOp_Variation(t *testing.T) {
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
		req     *ImageVariationRequest
		wantRes *ImageResponse
		wantErr error
		before  func()
		after   func()
	}{
		{
			name:   "test image variation success",
			client: newMockClient(server.URL),
			ctx:    context.TODO(),
			req: &ImageVariationRequest{
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
			name:   "test image variation timeout",
			client: newMockClient(server.URL),
			ctx: func() context.Context {
				ctx, cancel := context.WithTimeout(context.Background(), 2)
				cancels = append(cancels, cancel)
				return ctx
			}(),
			req: &ImageVariationRequest{
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
			name: "test image variation create form field error",
			client: newMockClient(server.URL, WithFormBuilder(func(w io.Writer) FormBuilder {
				return &mockFormBuilder{
					createFormFile: func(name string, filename string) error {
						return fmt.Errorf("%w, %s", os.ErrNotExist, "testdata/image.png")
					},
				}
			})),
			ctx: func() context.Context {
				ctx, cancel := context.WithTimeout(context.Background(), 2)
				cancels = append(cancels, cancel)
				return ctx
			}(),
			req: &ImageVariationRequest{
				Image: "testdata/image.png",
				ImageAttributes: ImageAttributes{
					N:    1,
					Size: "1024*1024",
				},
			},
			wantErr: os.ErrNotExist,
			before: func() {

			},
			after: func() {

			},
		},
		{
			name: "test image variation write form field error",
			client: newMockClient(server.URL, WithFormBuilder(func(w io.Writer) FormBuilder {
				return &mockFormBuilder{
					createFormFile: func(name string, filename string) error {
						return fmt.Errorf("%w, %s", os.ErrInvalid, "write form field error")
					},
				}
			})),
			ctx: func() context.Context {
				ctx, cancel := context.WithTimeout(context.Background(), 2)
				cancels = append(cancels, cancel)
				return ctx
			}(),
			req: &ImageVariationRequest{
				Image: "testdata/image.png",
				ImageAttributes: ImageAttributes{
					N:    1,
					Size: "1024*1024",
				},
			},
			wantErr: os.ErrInvalid,
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
			res, err := tc.client.Images.Variation(tc.ctx, tc.req)
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
	createFormFile  func(name string, filename string) error
	createFormField func(name string, value string) error
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
