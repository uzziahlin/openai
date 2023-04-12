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
	"net/http"
)

const (
	FilesListPath           = "/files"
	FileUploadPath          = "/files"
	FileDeletePath          = "/files/%s"
	FileRetrievePath        = "/files/%s"
	FileContentRetrievePath = "/files/%s/content"
)

type FileService interface {
	List(ctx context.Context) (*FileListResponse, error)
	Upload(ctx context.Context, req *FileUploadRequest) (*File, error)
	Delete(ctx context.Context, fileId string) (*FileDeleteResponse, error)
	Retrieve(ctx context.Context, fileId string) (*File, error)
	RetrieveContent(ctx context.Context, fileId string) ([]byte, error)
}

type FileListResponse struct {
	Data []*File `json:"data"`
}

type FileUploadRequest struct {
	File    string `json:"file"`
	Purpose string `json:"purpose"`
}

func (f FileUploadRequest) getFormFiles() []*FormFile {
	return []*FormFile{
		{
			fieldName: "file",
			filename:  f.File,
		},
	}
}

func (f FileUploadRequest) getFormFields() []*FormField {
	return []*FormField{
		{
			fieldName:  "purpose",
			fieldValue: f.Purpose,
		},
	}
}

type FileDeleteResponse struct {
	Id      string `json:"id"`
	Object  string `json:"object"`
	Deleted bool   `json:"deleted"`
}

type File struct {
	Id        string `json:"id"`
	Object    string `json:"object"`
	Bytes     int64  `json:"bytes"`
	CreatedAt int64  `json:"created_at"`
	Filename  string `json:"filename"`
	Purpose   string `json:"purpose"`
}

type FileServiceOp struct {
	client *Client
}

func (f FileServiceOp) List(ctx context.Context) (*FileListResponse, error) {
	var resp FileListResponse
	err := f.client.Get(ctx, FilesListPath, nil, &resp)
	return &resp, err
}

func (f FileServiceOp) Upload(ctx context.Context, req *FileUploadRequest) (*File, error) {
	var resp File
	err := f.client.Upload(ctx, FileUploadPath, req.getFormFiles(), &resp, req.getFormFields()...)
	return &resp, err
}

func (f FileServiceOp) Delete(ctx context.Context, fileId string) (*FileDeleteResponse, error) {
	var resp FileDeleteResponse
	err := f.client.Delete(ctx, fmt.Sprintf(FileDeletePath, fileId), nil, &resp)
	return &resp, err
}

func (f FileServiceOp) Retrieve(ctx context.Context, fileId string) (*File, error) {
	var resp File
	err := f.client.Get(ctx, fmt.Sprintf(FileRetrievePath, fileId), nil, &resp)
	return &resp, err
}

func (f FileServiceOp) RetrieveContent(ctx context.Context, fileId string) ([]byte, error) {
	return f.client.GetBytes(ctx, http.MethodGet, fmt.Sprintf(FileContentRetrievePath, fileId), nil, nil, nil)
}
