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
	"net/http"
	"os"
	"testing"
)

func TestFileServiceOp_List(t *testing.T) {
	server := newMockServer(newMockHandler(t, "GET", "files_list_response.json"))
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
		wantRes *FileListResponse
		wantErr error
	}{
		{
			name: "test files list success",
			ctx:  context.TODO(),
			wantRes: func() *FileListResponse {
				var wantRes FileListResponse
				loadMockData("files_list_response.json", &wantRes)
				return &wantRes
			}(),
		},
		{
			name: "test files list timeout",
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
			res, err := client.Files.List(tc.ctx)
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

func TestFileServiceOp_Upload(t *testing.T) {
	server := newMockServer(newMockHandler(t, "POST", "file_upload_response.json"))
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
		req     *FileUploadRequest
		wantRes *File
		wantErr error
		before  func()
		after   func()
	}{
		{
			name:   "test file upload success",
			client: newMockClient(server.URL),
			ctx:    context.TODO(),
			req: &FileUploadRequest{
				File:    "testdata/file_upload.json",
				Purpose: "fine-tune",
			},
			wantRes: func() *File {
				var wantRes File
				loadMockData("file_upload_response.json", &wantRes)
				return &wantRes
			}(),
			before: func() {
				file, err := os.Create("testdata/file_upload.json")
				require.NoError(t, err)
				file.Close()
			},
			after: func() {
				err := os.Remove("testdata/file_upload.json")
				require.NoError(t, err)
			},
		},
		{
			name:   "test file upload timeout",
			client: newMockClient(server.URL),
			ctx: func() context.Context {
				ctx, cancel := context.WithTimeout(context.Background(), 2)
				cancels = append(cancels, cancel)
				return ctx
			}(),
			req: &FileUploadRequest{
				File:    "testdata/file_upload.json",
				Purpose: "fine-tune",
			},
			wantErr: context.DeadlineExceeded,
			before: func() {
				file, err := os.Create("testdata/file_upload.json")
				require.NoError(t, err)
				file.Close()
			},
			after: func() {
				err := os.Remove("testdata/file_upload.json")
				require.NoError(t, err)
			},
		},
		{
			name: "test file upload create form file error",
			client: newMockClient(server.URL, WithFormBuilder(func(w io.Writer) FormBuilder {
				return &mockFormBuilder{
					createFormFile: func(name string, filename string) error {
						return fmt.Errorf("%w, %s", os.ErrNotExist, "testdata/file_upload.json")
					},
				}
			})),
			ctx: func() context.Context {
				ctx, cancel := context.WithTimeout(context.Background(), 2)
				cancels = append(cancels, cancel)
				return ctx
			}(),
			req: &FileUploadRequest{
				File:    "testdata/file_upload.json",
				Purpose: "fine-tune",
			},
			wantErr: os.ErrNotExist,
			before: func() {

			},
			after: func() {

			},
		},
		{
			name: "test file upload create form field error",
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
			req: &FileUploadRequest{
				File:    "testdata/file_upload.json",
				Purpose: "fine-tune",
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
			res, err := tc.client.Files.Upload(tc.ctx, tc.req)
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

func TestFileServiceOp_Delete(t *testing.T) {
	server := newMockServer(newMockHandler(t, "DELETE", "file_delete_response.json"))
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
		fileId  string
		wantRes *FileDeleteResponse
		wantErr error
	}{
		{
			name:   "test file delete success",
			ctx:    context.TODO(),
			fileId: "file-XjGxS3KTG0uNmNOK362iJua3",
			wantRes: func() *FileDeleteResponse {
				var wantRes FileDeleteResponse
				loadMockData("file_delete_response.json", &wantRes)
				return &wantRes
			}(),
		},
		{
			name: "test file delete timeout",
			ctx: func() context.Context {
				ctx, cancel := context.WithTimeout(context.Background(), 2)
				cancels = append(cancels, cancel)
				return ctx
			}(),
			fileId:  "file-XjGxS3KTG0uNmNOK362iJua3",
			wantErr: context.DeadlineExceeded,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			res, err := client.Files.Delete(tc.ctx, tc.fileId)
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

func TestFileServiceOp_Retrieve(t *testing.T) {
	server := newMockServer(newMockHandler(t, "GET", "file_retrieve_response.json"))
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
		fileId  string
		wantRes *File
		wantErr error
	}{
		{
			name:   "test file retrieve success",
			ctx:    context.TODO(),
			fileId: "file-XjGxS3KTG0uNmNOK362iJua3",
			wantRes: func() *File {
				var wantRes File
				loadMockData("file_retrieve_response.json", &wantRes)
				return &wantRes
			}(),
		},
		{
			name: "test file retrieve timeout",
			ctx: func() context.Context {
				ctx, cancel := context.WithTimeout(context.Background(), 2)
				cancels = append(cancels, cancel)
				return ctx
			}(),
			fileId:  "file-XjGxS3KTG0uNmNOK362iJua3",
			wantErr: context.DeadlineExceeded,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			res, err := client.Files.Retrieve(tc.ctx, tc.fileId)
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

func TestFileServiceOp_RetrieveContent(t *testing.T) {

	cancels := make([]context.CancelFunc, 0)
	defer func() {
		for _, cancel := range cancels {
			cancel()
		}
	}()

	testCase := []struct {
		name    string
		ctx     context.Context
		fileId  string
		wantRes []byte
		wantErr error
		handler http.HandlerFunc
	}{
		{
			name:   "test file retrieve content success with json",
			ctx:    context.TODO(),
			fileId: "file-XjGxS3KTG0uNmNOK362iJua3",
			wantRes: func() []byte {
				return loadTestdata("mock_file_content.json")
			}(),
			handler: func(w http.ResponseWriter, r *http.Request) {

				require.Equal(t, "GET", r.Method)

				w.Write(loadTestdata("mock_file_content.json"))
			},
		},
		{
			name:   "test file retrieve content success not json",
			ctx:    context.TODO(),
			fileId: "file-XjGxS3KTG0uNmNOK362iJua3",
			wantRes: func() []byte {
				return loadTestdata("mock_file_content.txt")
			}(),
			handler: func(w http.ResponseWriter, r *http.Request) {

				require.Equal(t, "GET", r.Method)

				w.Write(loadTestdata("mock_file_content.txt"))
			},
		},
		{
			name: "test file retrieve timeout",
			ctx: func() context.Context {
				ctx, cancel := context.WithTimeout(context.Background(), 2)
				cancels = append(cancels, cancel)
				return ctx
			}(),
			fileId:  "file-XjGxS3KTG0uNmNOK362iJua3",
			wantErr: context.DeadlineExceeded,
		},
	}

	for _, tc := range testCase {
		t.Run(tc.name, func(t *testing.T) {
			server := newMockServer(tc.handler)
			client := newMockClient(server.URL)
			defer server.Close()

			res, err := client.Files.RetrieveContent(tc.ctx, tc.fileId)
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
