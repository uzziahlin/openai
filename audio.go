package openai

import "context"

type AudioService interface {
	Transcriptions(ctx context.Context, req TranscriptionsRequest) (*TranscriptionsResponse, error)
	Translations(ctx context.Context, req TranslationsRequest) (*TranslationsResponse, error)
}

type TranscriptionsRequest struct {
	File           string  `json:"file"`
	Model          string  `json:"model"`
	Prompt         string  `json:"prompt,omitempty"`
	ResponseFormat string  `json:"response_format,omitempty"`
	Temperature    float64 `json:"temperature,omitempty"`
	Language       string  `json:"language,omitempty"`
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

type TranslationsResponse struct {
	Text string `json:"text"`
}

type AudioServiceOp struct {
	client *Client
}

func (a AudioServiceOp) Transcriptions(ctx context.Context, req TranscriptionsRequest) (*TranscriptionsResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (a AudioServiceOp) Translations(ctx context.Context, req TranslationsRequest) (*TranslationsResponse, error) {
	//TODO implement me
	panic("implement me")
}
