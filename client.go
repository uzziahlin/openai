package openai

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"github.com/google/go-querystring/query"
	"github.com/uzziahlin/transport/http"
	"io"
	"io/ioutil"
	"net/url"
	"strconv"
)

type App struct {
	ApiUrl string
	ApiKey string
}

type Option func(*Client)

func NewClient(app App, opts ...Option) (*Client, error) {

	u, err := url.Parse(app.ApiUrl)

	if err != nil {
		return nil, err
	}

	c := &Client{
		baseURL: u,
		apiKey:  app.ApiKey,
		client:  http.NewDefaultClient(),
		logger: &LeveledLogger{
			Level: LevelDebug,
		},
	}

	c.Chat = &chatServiceOp{
		client: c,
	}

	for _, opt := range opts {
		opt(c)
	}

	return c, nil
}

func WithRetries(retries int) Option {
	return func(c *Client) {
		c.retries = retries
	}
}

type Client struct {
	client  http.Client
	baseURL *url.URL
	apiKey  string
	retries int

	logger Logger
	Chat   ChatService
}

func (c *Client) Post(ctx context.Context, relPath string, body, resp any) error {
	return c.Do(ctx, "POST", relPath, nil, nil, body, resp)
}

func (c *Client) Get(ctx context.Context, relPath string, params, resp any) error {
	return c.Do(ctx, "GET", relPath, nil, params, nil, resp)
}

func (c *Client) Do(ctx context.Context, method, relPath string, headers map[string]string, params, body, resp any) error {

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

	_, err = c.doAndGetHeaders(ctx, req, resp, false)

	return err
}

func (c *Client) doAndGetHeaders(ctx context.Context, r *http.Request, v any, skipBody bool) (http.Header, error) {

	var (
		resp     *http.Response
		err      error
		attempts int
	)

	c.logRequest(r, skipBody)

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

	c.logResponse(resp, skipBody)
	defer resp.Body.Close()

	if v != nil {
		err = json.NewDecoder(resp.Body).Decode(&v)
		if err != nil {
			return nil, err
		}
	}

	return resp.Header, nil
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
