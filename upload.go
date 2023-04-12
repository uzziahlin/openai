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
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
)

// Upload file to openai
// files: file to upload
// v: response data
// fields: other fields
func (c *Client) Upload(ctx context.Context, relPath string, files []*FormFile, v any, fields ...*FormField) error {

	rel, err := url.Parse(relPath)
	if err != nil {
		return err
	}

	u := c.baseURL.ResolveReference(rel)

	form := &bytes.Buffer{}

	builder := c.formBuilder(form)

	if files != nil && len(files) > 0 {
		for _, file := range files {
			err := builder.CreateFormFile(file.fieldName, file.filename)
			if err != nil {
				return err
			}
		}
	}

	if fields != nil && len(fields) > 0 {
		for _, field := range fields {
			err := builder.CreateFormField(field.fieldName, field.fieldValue)
			if err != nil {
				return err
			}
		}
	}

	err = builder.Close()
	if err != nil {
		return err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", u.String(), form)
	if err != nil {
		return err
	}

	req.Header.Set("Content-Type", builder.FormDataContentType())
	req.Header.Set("Authorization", "Bearer "+c.apiKey)
	req.Header.Set("Accept", "application/json")

	resp, err := c.do(ctx, req, true, false)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(&v)
	}

	return nil
}

// FormBuilder is a builder for multipart/form-data requests.
type FormBuilder interface {
	CreateFormFile(name string, filename string) error
	CreateFormField(name string, value string) error
	FormDataContentType() string
	io.Closer
}

func NewMultiPartFormBuilder(w io.Writer) FormBuilder {
	writer := multipart.NewWriter(w)
	return &multipartFormBuilder{
		w: writer,
	}
}

type multipartFormBuilder struct {
	w *multipart.Writer
}

func (m multipartFormBuilder) CreateFormFile(name string, filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}
	defer f.Close()

	formFile, err := m.w.CreateFormFile(name, filepath.Base(filename))
	if err != nil {
		return err
	}

	_, err = io.Copy(formFile, f)
	if err != nil {
		return err
	}

	return nil
}

func (m multipartFormBuilder) CreateFormField(name string, value string) error {
	return m.w.WriteField(name, value)
}

func (m multipartFormBuilder) FormDataContentType() string {
	return m.w.FormDataContentType()
}

func (m multipartFormBuilder) Close() error {
	return m.w.Close()
}

// MultipartRequestBuilder is a builder for multipart/form-data requests.
// 这里抽象没有做好，后期考虑重构
/*type MultipartRequestBuilder struct {
	baseUrl *url.URL
	relPath string
	Files   []*FormFile
	Fields  []*FormField
}

func (m MultipartRequestBuilder) Build(ctx context.Context) (*http.Request, error) {
	rel, err := url.Parse(m.relPath)
	if err != nil {
		return nil, err
	}

	u := m.baseUrl.ResolveReference(rel)

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	if files := m.Files; files != nil && len(files) > 0 {
		for _, file := range files {
			f, err := os.Open(file.filename)
			if err != nil {
				return nil, err
			}

			formFile, err := writer.CreateFormFile(file.fieldName, filepath.Base(file.filename))
			if err != nil {
				_ = f.Close()
				return nil, err
			}

			_, err = io.Copy(formFile, f)
			if err != nil {
				_ = f.Close()
				return nil, err
			}

			err = f.Close()

			if err != nil {
				return nil, err
			}
		}
	}

	if fields := m.Fields; fields != nil && len(fields) > 0 {
		for _, field := range fields {
			formField, err := writer.CreateFormField(field.fieldName)
			if err != nil {
				return nil, err
			}
			// todo 是否需要将value定义为io.Reader？
			_, err = formField.Write([]byte(field.fieldValue))
			if err != nil {
				return nil, err
			}
		}
	}

	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, "POST", u.String(), body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())
	return req, nil
}*/

type FormFile struct {
	fieldName string
	filename  string
}

func NewFormFile(fieldName, filename string) *FormFile {
	return &FormFile{
		fieldName: fieldName,
		filename:  filename,
	}
}

type FormField struct {
	fieldName  string
	fieldValue string
}

func NewFormField(fieldName, fieldValue string) *FormField {
	return &FormField{
		fieldName:  fieldName,
		fieldValue: fieldValue,
	}
}
