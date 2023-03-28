package openai

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/google/go-querystring/query"
	"github.com/uzziahlin/transport/http"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"
)

type App struct {
	ApiUrl string
	ApiKey string
}

type Option func(*Client)

func NewClient(app App, opts ...Option) *Client {

	u, err := url.Parse(app.ApiUrl)

	if err != nil {
		panic(err)
	}

	c := &Client{
		baseURL: u,
		apiKey:  app.ApiKey,
		logger: &LeveledLogger{
			Level: LevelDebug,
		},
	}

	c.Models = &ModelServiceOp{
		client: c,
	}

	c.Chat = &ChatServiceOp{
		client: c,
	}

	c.Images = &ImageServiceOp{
		client: c,
	}

	for _, opt := range opts {
		opt(c)
	}

	if c.proxyUrl != nil {
		c.client = http.NewDefaultClient(http.WithProxy(c.proxyUrl))
	} else {
		c.client = http.NewDefaultClient()
	}

	return c
}

func WithProxy(proxyUrl string) Option {
	return func(c *Client) {
		u, err := url.Parse(proxyUrl)
		if err != nil {
			c.logError(err, "failed to parse proxy url: %s")
			return
		}
		c.proxyUrl = u
	}
}

func WithRetries(retries int) Option {
	return func(c *Client) {
		c.retries = retries
	}
}

func WithLogger(logger Logger) Option {
	return func(c *Client) {
		c.logger = logger
	}
}

type Client struct {
	client   http.Client
	baseURL  *url.URL
	proxyUrl *url.URL
	apiKey   string
	retries  int

	logger Logger

	Models ModelService
	Chat   ChatService
	Images ImageService
}

func (c *Client) Stream(ctx context.Context, relPath string, body any) (EventSource, error) {

	headers := make(map[string]string, 3)
	headers["Content-Type"] = "application/json"
	headers["Accept"] = "text/event-stream"
	headers["Authorization"] = "Bearer " + c.apiKey

	req, err := c.NewRequest("POST", relPath, headers, nil, body)

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
	return c.Do(ctx, "POST", relPath, nil, nil, body, resp)
}

func (c *Client) Get(ctx context.Context, relPath string, params, resp any) error {
	return c.Do(ctx, "GET", relPath, nil, params, nil, resp)
}

func (c *Client) Do(ctx context.Context, method, relPath string, headers map[string]string, params, body, v any) error {

	if headers == nil {
		headers = make(map[string]string, 3)
	}

	headers["Content-Type"] = "application/json"
	headers["Accept"] = "application/json"
	headers["Authorization"] = "Bearer " + c.apiKey

	req, err := c.NewRequest(method, relPath, headers, params, body)

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

func (c *Client) do(ctx context.Context, r *http.Request, skipReqBody, skipRespBody bool) (*http.Response, error) {

	var (
		resp     *http.Response
		err      error
		attempts int
	)

	c.logRequest(r, skipReqBody)

	for {
		attempts++

		resp, err = c.client.Send(ctx, r)

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

func (c *Client) NewRequest(method string, relPath string, headers map[string]string, params any, body any) (*http.Request, error) {
	rel, err := url.Parse(relPath)
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

	req := &http.Request{
		Method: method,
		Url:    u.String(),
		Body:   ioutil.NopCloser(bytes.NewBuffer(data)),
	}

	if headers != nil {
		for k, v := range headers {
			req.SetHeader(k, v)
		}
	}

	return req, nil

}

// skipBody: if upload image, skip log its binary
func (c *Client) logRequest(req *http.Request, skipBody bool) {
	if req == nil {
		return
	}
	if req.Url != "" {
		c.logger.Debugf("%s: %s", req.Method, req.Url)
	}
	if !skipBody {
		c.logBody(&req.Body, "SENT: %s")
	}
}

func (c *Client) logResponse(res *http.Response, skipBody bool) {
	if res == nil {
		return
	}
	c.logger.Debugf("RECV %d: %s", res.StatusCode, res.StatusCode)
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
		c.logger.Debugf(format, string(b))
	}
	*body = ioutil.NopCloser(bytes.NewBuffer(b))
}

func (c *Client) logError(err error, format string) {
	if err == nil {
		return
	}

	c.logger.Errorf(format, err.Error())

}

func (c *Client) Upload(ctx context.Context, relPath string, files []*FormFile, v any, fields ...*FormField) error {
	builder := MultipartRequestBuilder{
		baseUrl: c.baseURL,
		relPath: relPath,
		Files:   files,
		Fields:  fields,
	}

	request, err := builder.Build()
	if err != nil {
		return err
	}

	resp, err := c.do(ctx, request, true, false)

	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(&v)
	}

	return nil
}

type MultipartRequestBuilder struct {
	baseUrl *url.URL
	relPath string
	Files   []*FormFile
	Fields  []*FormField
}

func (m MultipartRequestBuilder) Build() (*http.Request, error) {
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
			defer f.Close()

			formFile, err := writer.CreateFormFile(file.fieldName, filepath.Base(file.filename))
			if err != nil {
				return nil, err
			}

			_, err = io.Copy(formFile, f)
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
			// todo 是否需要讲value定义为io.Reader？
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

	req := &http.Request{
		Method: "POST",
		Url:    u.String(),
		Body:   ioutil.NopCloser(body),
	}

	return req, nil
}

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
