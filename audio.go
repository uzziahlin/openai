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
	"strconv"
)

const (
	AudioTranscriptionsPath = "/audio/transcriptions"
	AudioTranslationsPath   = "/audio/translations"
)

type AudioService interface {
	Transcriptions(ctx context.Context, req *TranscriptionsRequest) (*TranscriptionsResponse, error)
	Translations(ctx context.Context, req *TranslationsRequest) (*TranslationsResponse, error)
}

type TranscriptionsRequest struct {
	// File is the path to the audio file
	File           string  `json:"file"`
	Model          string  `json:"model"`
	Prompt         string  `json:"prompt,omitempty"`
	ResponseFormat string  `json:"response_format,omitempty"`
	Temperature    float64 `json:"temperature,omitempty"`
	Language       string  `json:"language,omitempty"`
}

func (t TranscriptionsRequest) getFormFiles() []*FormFile {
	return []*FormFile{
		{
			fieldName: "file",
			filename:  t.File,
		},
	}
}

func (t TranscriptionsRequest) getFormFields() []*FormField {
	formFields := []*FormField{
		{
			fieldName:  "model",
			fieldValue: t.Model,
		},
	}

	if t.Prompt != "" {
		formFields = append(formFields, &FormField{
			fieldName:  "prompt",
			fieldValue: t.Prompt,
		})
	}

	if t.ResponseFormat != "" {
		formFields = append(formFields, &FormField{
			fieldName:  "response_format",
			fieldValue: t.ResponseFormat,
		})
	}

	if t.Temperature != 0 {
		formFields = append(formFields, &FormField{
			fieldName:  "temperature",
			fieldValue: strconv.FormatFloat(t.Temperature, 'f', -1, 64),
		})
	}

	if t.Language != "" {
		formFields = append(formFields, &FormField{
			fieldName:  "language",
			fieldValue: t.Language,
		})
	}

	return formFields
}

type TranscriptionsResponse struct {
	Text string `json:"text"`
}

type TranslationsRequest struct {
	File           string  `json:"file"`
	Model          string  `json:"model"`
	Prompt         string  `json:"prompt,omitempty"`
	ResponseFormat string  `json:"response_format,omitempty"`
	Temperature    float64 `json:"temperature,omitempty"`
}

func (t TranslationsRequest) getFormFiles() []*FormFile {
	return []*FormFile{
		{
			fieldName: "file",
			filename:  t.File,
		},
	}
}

func (t TranslationsRequest) getFormFields() []*FormField {
	formFields := []*FormField{
		{
			fieldName:  "model",
			fieldValue: t.Model,
		},
	}

	if t.Prompt != "" {
		formFields = append(formFields, &FormField{
			fieldName:  "prompt",
			fieldValue: t.Prompt,
		})
	}

	if t.ResponseFormat != "" {
		formFields = append(formFields, &FormField{
			fieldName:  "response_format",
			fieldValue: t.ResponseFormat,
		})
	}

	if t.Temperature != 0 {
		formFields = append(formFields, &FormField{
			fieldName:  "temperature",
			fieldValue: strconv.FormatFloat(t.Temperature, 'f', -1, 64),
		})
	}

	return formFields
}

type TranslationsResponse struct {
	Text string `json:"text"`
}

type AudioServiceOp struct {
	client *Client
}

func (a AudioServiceOp) Transcriptions(ctx context.Context, req *TranscriptionsRequest) (*TranscriptionsResponse, error) {
	var resp TranscriptionsResponse
	err := a.client.Upload(ctx, AudioTranscriptionsPath, req.getFormFiles(), &resp, req.getFormFields()...)
	return &resp, err
}

func (a AudioServiceOp) Translations(ctx context.Context, req *TranslationsRequest) (*TranslationsResponse, error) {
	var resp TranslationsResponse
	err := a.client.Upload(ctx, AudioTranslationsPath, req.getFormFiles(), &resp, req.getFormFields()...)
	return &resp, err
}
