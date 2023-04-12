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
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/go-logr/logr"
	"github.com/go-logr/zapr"
	"github.com/google/go-querystring/query"
	"go.uber.org/zap"
	"io"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"
)

type App struct {
	ApiUrl string
	ApiKey string
}

type Option func(*Client)

func New(app App, opts ...Option) (*Client, error) {

	u, err := url.Parse(app.ApiUrl)

	if err != nil {
		return nil, err
	}

	// 默认使用zap，环境默认为开发环境，如果有特殊要求，使用者可以自行注入日志实现
	logger, err := zap.NewDevelopment()
	if err != nil {
		return nil, err
	}

	c := &Client{
		client:      &http.Client{},
		baseURL:     u,
		version:     "v1",
		apiKey:      app.ApiKey,
		retries:     3,
		formBuilder: NewMultiPartFormBuilder,
		logger:      zapr.NewLogger(logger),
	}

	c.Models = &ModelServiceOp{
		client: c,
	}

	c.Completions = &CompletionServiceOp{
		client: c,
	}

	c.Chat = &ChatServiceOp{
		client: c,
	}

	c.Edits = &EditServiceOp{
		client: c,
	}

	c.Images = &ImageServiceOp{
		client: c,
	}

	c.Embeddings = &EmbeddingServiceOp{
		client: c,
	}

	c.Audio = &AudioServiceOp{
		client: c,
	}

	c.Files = &FileServiceOp{
		client: c,
	}

	c.FineTunes = &FineTuneServiceOp{
		client: c,
	}

	c.Moderations = &ModerationServiceOp{
		client: c,
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.proxy != nil && c.proxy.Url != "" {
		proxyUrl, err := url.Parse(c.proxy.Url)
		if err != nil {
			c.logger.Error(err, "parse proxy url error")
			return nil, err
		}

		if c.proxy.Username != "" {
			var user *url.Userinfo
			if c.proxy.Password != "" {
				user = url.UserPassword(c.proxy.Username, c.proxy.Password)
			} else {
				user = url.User(c.proxy.Username)
			}
			proxyUrl.User = user
		}

		transport := &http.Transport{
			Proxy: http.ProxyURL(proxyUrl),
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		}

		if c.proxy.TlsCfg != nil {
			transport.TLSClientConfig = c.proxy.TlsCfg
		}

		// 手动设置代理认证信息
		/*if c.proxy.Username != "" {
			proxyAuth := fmt.Sprintf("%s:%s", c.proxy.Username, c.proxy.Password)
			basicAuth := "Basic " + base64.StdEncoding.EncodeToString([]byte(proxyAuth))
			transport.ProxyConnectHeader = http.Header{
				"Proxy-Authorization": []string{basicAuth},
			}
		}*/

		c.client.Transport = transport
	}

	return c, nil
}

func WithProxy(proxy *Proxy) Option {
	return func(c *Client) {
		c.proxy = proxy
	}
}

func WithRetries(retries int) Option {
	return func(c *Client) {
		c.retries = retries
	}
}

func WithFormBuilder(builder func(w io.Writer) FormBuilder) Option {
	return func(c *Client) {
		c.formBuilder = builder
	}
}

func WithLogger(logger logr.Logger) Option {
	return func(c *Client) {
		c.logger = logger
	}
}

// WithVersion 设置默认版本，如果不设置，默认为v1
func WithVersion(version string) Option {
	return func(c *Client) {
		c.version = version
	}
}

type Proxy struct {
	Url      string
	Username string
	Password string
	TlsCfg   *tls.Config
}

type Client struct {
	client  *http.Client
	baseURL *url.URL
	version string
	proxy   *Proxy
	apiKey  string
	retries int

	formBuilder func(w io.Writer) FormBuilder

	logger logr.Logger

	Models      ModelService
	Completions CompletionService
	Chat        ChatService
	Edits       EditService
	Images      ImageService
	Embeddings  EmbeddingService
	Audio       AudioService
	Files       FileService
	FineTunes   FineTuneService
	Moderations ModerationService
}

// V 设置版本,返回一个新的Client实例，不会修改原有实例
func (c *Client) V(version string) Client {
	newClient := *c
	newClient.version = version
	return newClient
}

func (c *Client) Close() error {
	c.client.CloseIdleConnections()
	return nil
}

func (c *Client) GetByStream(ctx context.Context, relPath string, params any) (EventSource, error) {
	return c.Stream(ctx, http.MethodGet, relPath, nil, params, nil)
}

func (c *Client) PostByStream(ctx context.Context, relPath string, body any) (EventSource, error) {
	return c.Stream(ctx, http.MethodPost, relPath, nil, nil, body)
}

// Stream 为请求提供流式处理
func (c *Client) Stream(ctx context.Context, method, relPath string, headers map[string]string, params, body any) (EventSource, error) {

	if headers == nil {
		headers = make(map[string]string, 3)
	}

	headers["Content-Type"] = "application/json"
	headers["Accept"] = "text/event-stream"
	headers["Authorization"] = "Bearer " + c.apiKey

	req, err := c.NewRequest(ctx, method, relPath, headers, params, body)

	if err != nil {
		return nil, err
	}

	resp, err := c.do(ctx, req, false, true)

	if err != nil {
		return nil, err
	}

	es := NewEventSource(ctx, resp.Body, "[DONE]")

	return es, nil
}

// NewEventSource 处理SSE
func NewEventSource(ctx context.Context, r io.ReadCloser, doneStr string) EventSource {
	es := make(EventSource)

	go func() {
		defer func() {
			_ = r.Close()
			close(es)
		}()
		scanner := bufio.NewScanner(r)
		var event Event
		for scanner.Scan() {
			line := scanner.Text()
			if line == "" {
				select {
				case <-ctx.Done():
					return
				case es <- event:
					event = Event{}
				}
			} else if strings.HasPrefix(line, "event:") {
				event.Event = strings.TrimSpace(line[len("event:"):])
			} else if strings.HasPrefix(line, "data:") {
				event.Data = strings.TrimSpace(line[len("data:"):])
				if event.Data == doneStr {
					return
				}
			} else if strings.HasPrefix(line, "id:") {
				event.Id = strings.TrimSpace(line[len("id:"):])
			} else if strings.HasPrefix(line, "retry:") {
				duration, err := time.ParseDuration(strings.TrimSpace(line[len("retry:"):]))
				if err != nil {
					event.Err = err
				} else {
					event.Retry = duration
				}
			}
			if ctx.Err() != nil {
				return
			}
		}
	}()

	return es
}

type EventSource chan Event

type Event struct {
	Id    string
	Event string
	Data  string
	Retry time.Duration
	Err   error
}

func (c *Client) Post(ctx context.Context, relPath string, body, resp any) error {
	return c.Do(ctx, http.MethodPost, relPath, nil, nil, body, resp)
}

func (c *Client) Get(ctx context.Context, relPath string, params, resp any) error {
	return c.Do(ctx, http.MethodGet, relPath, nil, params, nil, resp)
}

func (c *Client) Delete(ctx context.Context, relPath string, params, resp any) error {
	return c.Do(ctx, http.MethodDelete, relPath, nil, params, nil, resp)
}

func (c *Client) Do(ctx context.Context, method, relPath string, headers map[string]string, params, body, v any) error {

	if headers == nil {
		headers = make(map[string]string, 3)
	}

	headers["Content-Type"] = "application/json"
	headers["Accept"] = "application/json"
	headers["Authorization"] = "Bearer " + c.apiKey

	req, err := c.NewRequest(ctx, method, relPath, headers, params, body)

	if err != nil {
		return err
	}

	resp, err := c.do(ctx, req, false, false)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(&v)
	}

	return err
}

// GetBytes 获取字节流, 也可以考虑合并到Do中
// 但是由于api中大部分都是json, 所以这里单独提取出来
func (c *Client) GetBytes(ctx context.Context, method, relPath string, headers map[string]string, params, body any) ([]byte, error) {

	if headers == nil {
		headers = make(map[string]string, 3)
	}

	headers["Content-Type"] = "application/json"
	headers["Authorization"] = "Bearer " + c.apiKey

	req, err := c.NewRequest(ctx, method, relPath, headers, params, body)

	if err != nil {
		return nil, err
	}

	resp, err := c.do(ctx, req, false, false)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	return ioutil.ReadAll(resp.Body)
}

func (c *Client) do(ctx context.Context, r *http.Request, skipReqBody, skipRespBody bool) (*http.Response, error) {

	var (
		resp     *http.Response
		err      error
		attempts int
	)

	c.logRequest(r, skipReqBody)

	for {
		attempts++

		resp, err = c.client.Do(r)

		// 由客户端引起的错误，不需要重试
		if err != nil {
			return nil, err
		}

		// 检查是否有错误，目前只检查状态码
		err = checkErr(resp)

		if err == nil {
			break
		}

		// 重试次数超过限制
		if attempts > c.retries {
			return nil, err
		}
	}

	c.logResponse(resp, skipRespBody)

	return resp, nil
}

func checkErr(resp *http.Response) error {
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	return errors.New(strconv.Itoa(resp.StatusCode))
}

func (c *Client) NewRequest(ctx context.Context, method string, relPath string, headers map[string]string, params any, body any) (*http.Request, error) {
	rel, err := url.Parse(path.Join(c.version, relPath))
	if err != nil {
		return nil, err
	}

	u := c.baseURL.ResolveReference(rel)

	if params != nil {
		optionsQuery, err := query.Values(params)
		if err != nil {
			return nil, err
		}

		for k, values := range u.Query() {
			for _, v := range values {
				optionsQuery.Add(k, v)
			}
		}
		u.RawQuery = optionsQuery.Encode()
	}

	var data []byte
	if body != nil {
		data, err = json.Marshal(body)
		if err != nil {
			return nil, err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, u.String(), bytes.NewBuffer(data))

	if err != nil {
		return nil, err
	}

	if headers != nil {
		for k, v := range headers {
			req.Header.Set(k, v)
		}
	}

	return req, nil

}

// skipBody: if upload image, skip log its binary
func (c *Client) logRequest(req *http.Request, skipBody bool) {
	if req == nil {
		return
	}
	if req.URL != nil {
		// debug level
		c.logger.V(1).Info(fmt.Sprintf("%s: %s", req.Method, req.URL.String()))
	}
	if !skipBody {
		c.logBody(&req.Body, "SENT: %s")
	}
}

func (c *Client) logResponse(res *http.Response, skipBody bool) {
	if res == nil {
		return
	}
	// debug level
	c.logger.V(1).Info(fmt.Sprintf("RECV %d:", res.StatusCode))
	if !skipBody {
		c.logBody(&res.Body, "RESP: %s")
	}

}

func (c *Client) logBody(body *io.ReadCloser, format string) {
	if body == nil || *body == nil {
		return
	}
	b, _ := ioutil.ReadAll(*body)
	if len(b) > 0 {
		// debug level
		c.logger.V(1).Info(fmt.Sprintf(format, string(b)))
	}
	*body = ioutil.NopCloser(bytes.NewBuffer(b))
}
