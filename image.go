package openai

import "context"

type ImageService interface {
	Create(ctx context.Context, req *ImageCreateRequest) (*ImageResponse, error)
	Edit(ctx context.Context, req *ImageEditRequest) (*ImageResponse, error)
	Variation(ctx context.Context, req *ImageVariationRequest) (*ImageResponse, error)
}

type ImageCreateRequest struct {
	Prompt string `json:"prompt"`
	ImageAttributes
}

type ImageResponse struct {
	Created int64   `json:"created"`
	Data    []Image `json:"data"`
}

type Image struct {
	Url string `json:"url"`
}

type ImageEditRequest struct {
	Image  string `json:"image"`
	Mask   string `json:"mask,omitempty"`
	Prompt string `json:"prompt"`
	ImageAttributes
}

type ImageVariationRequest struct {
	Image string `json:"image"`
	ImageAttributes
}

type ImageAttributes struct {
	N              int    `json:"n,omitempty"`
	Size           string `json:"size,omitempty"`
	ResponseFormat string `json:"response_format,omitempty"`
	User           string `json:"user,omitempty"`
}

type ImageServiceOp struct {
	client *Client
}

func (i ImageServiceOp) Create(ctx context.Context, req *ImageCreateRequest) (*ImageResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (i ImageServiceOp) Edit(ctx context.Context, req *ImageEditRequest) (*ImageResponse, error) {
	//TODO implement me
	panic("implement me")
}

func (i ImageServiceOp) Variation(ctx context.Context, req *ImageVariationRequest) (*ImageResponse, error) {
	//TODO implement me
	panic("implement me")
}
