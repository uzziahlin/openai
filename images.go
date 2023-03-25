package openai

import "context"

type ImagesService interface {
	Create(ctx context.Context, req *ImagesCreateRequest) (*ImagesResponse, error)
	Edit(ctx context.Context, req *ImagesEditRequest) (*ImagesResponse, error)
	Variation(ctx context.Context, req *ImagesVariationRequest) (*ImagesResponse, error)
}

type ImagesCreateRequest struct {
	Prompt string `json:"prompt"`
	ImagesAttributes
}

type ImagesResponse struct {
	Created int64   `json:"created"`
	Data    []Image `json:"data"`
}

type Image struct {
	Url string `json:"url"`
}

type ImagesEditRequest struct {
	Image  string `json:"image"`
	Mask   string `json:"mask,omitempty"`
	Prompt string `json:"prompt"`
	ImagesAttributes
}

type ImagesVariationRequest struct {
	Image string `json:"image"`
	ImagesAttributes
}

type ImagesAttributes struct {
	N              int    `json:"n,omitempty"`
	Size           string `json:"size,omitempty"`
	ResponseFormat string `json:"response_format,omitempty"`
	User           string `json:"user,omitempty"`
}
