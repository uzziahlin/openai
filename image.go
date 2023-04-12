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
	ImageCreatePath    = "/images/generations"
	ImageEditPath      = "/images/edits"
	ImageVariationPath = "/images/variations"
)

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
	// Image is the path to the image file
	Image string `json:"image"`
	// Mask is the path to the mask file
	Mask   string `json:"mask,omitempty"`
	Prompt string `json:"prompt"`
	ImageAttributes
}

func (i *ImageEditRequest) getFormFiles() []*FormFile {
	res := []*FormFile{
		{
			fieldName: "image",
			filename:  i.Image,
		},
	}

	if i.Mask != "" {
		res = append(res, &FormFile{
			fieldName: "mask",
			filename:  i.Mask,
		})
	}

	return res
}

func (i *ImageEditRequest) getFormFields() []*FormField {
	res := []*FormField{
		{
			fieldName:  "prompt",
			fieldValue: i.Prompt,
		},
	}

	if i.N != 0 {
		res = append(res, &FormField{
			fieldName:  "n",
			fieldValue: strconv.Itoa(i.N),
		})
	}

	if i.Size != "" {
		res = append(res, &FormField{
			fieldName:  "size",
			fieldValue: i.Size,
		})
	}

	if i.ResponseFormat != "" {
		res = append(res, &FormField{
			fieldName:  "response_format",
			fieldValue: i.ResponseFormat,
		})
	}

	if i.User != "" {
		res = append(res, &FormField{
			fieldName:  "user",
			fieldValue: i.User,
		})
	}

	return res
}

type ImageVariationRequest struct {
	Image string `json:"image"`
	ImageAttributes
}

func (i *ImageVariationRequest) getFormFiles() []*FormFile {
	return []*FormFile{
		{
			fieldName: "image",
			filename:  i.Image,
		},
	}
}

func (i *ImageVariationRequest) getFormFields() []*FormField {
	var res []*FormField

	if i.N != 0 {
		res = append(res, &FormField{
			fieldName:  "n",
			fieldValue: strconv.Itoa(i.N),
		})
	}

	if i.Size != "" {
		res = append(res, &FormField{
			fieldName:  "size",
			fieldValue: i.Size,
		})
	}

	if i.ResponseFormat != "" {
		res = append(res, &FormField{
			fieldName:  "response_format",
			fieldValue: i.ResponseFormat,
		})
	}

	if i.User != "" {
		res = append(res, &FormField{
			fieldName:  "user",
			fieldValue: i.User,
		})
	}

	return res
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
	var resp ImageResponse
	err := i.client.Post(ctx, ImageCreatePath, req, &resp)
	return &resp, err
}

func (i ImageServiceOp) Edit(ctx context.Context, req *ImageEditRequest) (*ImageResponse, error) {
	var resp ImageResponse
	err := i.client.Upload(ctx, ImageEditPath, req.getFormFiles(), &resp, req.getFormFields()...)
	return &resp, err
}

func (i ImageServiceOp) Variation(ctx context.Context, req *ImageVariationRequest) (*ImageResponse, error) {
	var resp ImageResponse
	err := i.client.Upload(ctx, ImageVariationPath, req.getFormFiles(), &resp, req.getFormFields()...)
	return &resp, err
}
