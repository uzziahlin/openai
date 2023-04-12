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

func TestAudioServiceOp_Transcriptions(t *testing.T) {
	server := newMockServer(newMockHandler(t, "POST", "audio_transcriptions_response.json"))
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
		req     *TranscriptionsRequest
		wantRes *TranscriptionsResponse
		wantErr error
		before  func()
		after   func()
	}{
		{
			name:   "test audio transcriptions success",
			client: newMockClient(server.URL),
			ctx:    context.TODO(),
			req: &TranscriptionsRequest{
				File:  "testdata/audio_transcriptions.mp3",
				Model: "whisper-1",
			},
			wantRes: func() *TranscriptionsResponse {
				var wantRes TranscriptionsResponse
				loadMockData("audio_transcriptions_response.json", &wantRes)
				return &wantRes
			}(),
			before: func() {
				file, err := os.Create("testdata/audio_transcriptions.mp3")
				require.NoError(t, err)
				file.Close()
			},
			after: func() {
				err := os.Remove("testdata/audio_transcriptions.mp3")
				require.NoError(t, err)
			},
		},
		{
			name:   "test audio transcriptions timeout",
			client: newMockClient(server.URL),
			ctx: func() context.Context {
				ctx, cancel := context.WithTimeout(context.Background(), 2)
				cancels = append(cancels, cancel)
				return ctx
			}(),
			req: &TranscriptionsRequest{
				File:  "testdata/audio_transcriptions.mp3",
				Model: "whisper-1",
			},
			wantErr: context.DeadlineExceeded,
			before: func() {
				file, err := os.Create("testdata/audio_transcriptions.mp3")
				require.NoError(t, err)
				file.Close()
			},
			after: func() {
				err := os.Remove("testdata/audio_transcriptions.mp3")
				require.NoError(t, err)
			},
		},
		{
			name: "test audio transcriptions create form file error",
			client: newMockClient(server.URL, WithFormBuilder(func(w io.Writer) FormBuilder {
				return &mockFormBuilder{
					createFormFile: func(name string, filename string) error {
						return fmt.Errorf("%w, %s", os.ErrNotExist, "testdata/audio_transcriptions.mp3")
					},
				}
			})),
			ctx: func() context.Context {
				ctx, cancel := context.WithTimeout(context.Background(), 2)
				cancels = append(cancels, cancel)
				return ctx
			}(),
			req: &TranscriptionsRequest{
				File:  "testdata/audio_transcriptions.mp3",
				Model: "whisper-1",
			},
			wantErr: os.ErrNotExist,
			before: func() {

			},
			after: func() {

			},
		},
		{
			name: "test audio transcriptions create form field error",
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
			req: &TranscriptionsRequest{
				File:  "testdata/audio_transcriptions.mp3",
				Model: "whisper-1",
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
			res, err := tc.client.Audio.Transcriptions(tc.ctx, tc.req)
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

func TestAudioServiceOp_Translations(t *testing.T) {
	server := newMockServer(newMockHandler(t, "POST", "audio_translations_response.json"))
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
		req     *TranslationsRequest
		wantRes *TranslationsResponse
		wantErr error
		before  func()
		after   func()
	}{
		{
			name:   "test audio translations success",
			client: newMockClient(server.URL),
			ctx:    context.TODO(),
			req: &TranslationsRequest{
				File:  "testdata/audio_translations.mp3",
				Model: "whisper-1",
			},
			wantRes: func() *TranslationsResponse {
				var wantRes TranslationsResponse
				loadMockData("audio_translations_response.json", &wantRes)
				return &wantRes
			}(),
			before: func() {
				file, err := os.Create("testdata/audio_translations.mp3")
				require.NoError(t, err)
				file.Close()
			},
			after: func() {
				err := os.Remove("testdata/audio_translations.mp3")
				require.NoError(t, err)
			},
		},
		{
			name:   "test audio translations timeout",
			client: newMockClient(server.URL),
			ctx: func() context.Context {
				ctx, cancel := context.WithTimeout(context.Background(), 2)
				cancels = append(cancels, cancel)
				return ctx
			}(),
			req: &TranslationsRequest{
				File:  "testdata/audio_translations.mp3",
				Model: "whisper-1",
			},
			wantErr: context.DeadlineExceeded,
			before: func() {
				file, err := os.Create("testdata/audio_translations.mp3")
				require.NoError(t, err)
				file.Close()
			},
			after: func() {
				err := os.Remove("testdata/audio_translations.mp3")
				require.NoError(t, err)
			},
		},
		{
			name: "test audio translations create form file error",
			client: newMockClient(server.URL, WithFormBuilder(func(w io.Writer) FormBuilder {
				return &mockFormBuilder{
					createFormFile: func(name string, filename string) error {
						return fmt.Errorf("%w, %s", os.ErrNotExist, "testdata/audio_translations.mp3")
					},
				}
			})),
			ctx: func() context.Context {
				ctx, cancel := context.WithTimeout(context.Background(), 2)
				cancels = append(cancels, cancel)
				return ctx
			}(),
			req: &TranslationsRequest{
				File:  "testdata/audio_translations.mp3",
				Model: "whisper-1",
			},
			wantErr: os.ErrNotExist,
			before: func() {

			},
			after: func() {

			},
		},
		{
			name: "test audio transcriptions create form field error",
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
			req: &TranslationsRequest{
				File:  "testdata/audio_translations.mp3",
				Model: "whisper-1",
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
			res, err := tc.client.Audio.Translations(tc.ctx, tc.req)
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
