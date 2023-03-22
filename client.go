package openai

import (
	"context"
	"github.com/uzziahlin/transport/http"
)

type Client struct {
	http.Client
}

func NewClient() *Client {
	return &Client{}
}

func (c *Client) Post(ctx context.Context, url string, body, resp any) error {
	return c.Do(ctx, "POST", url, nil, nil, body, resp)
}

func (c *Client) Get(ctx context.Context, url string, params, resp any) error {
	return c.Do(ctx, "GET", url, nil, params, nil, resp)
}

func (c *Client) Do(ctx context.Context, method, url string, headers map[string]string, params, body, resp any) error {

	req, err := c.NewRequest(method, url, headers, params, body)

	if err != nil {
		return err
	}

	_, err = c.doAndGetHeaders(ctx, req, resp)

	return err
}

func (c *Client) doAndGetHeaders(ctx context.Context, r *http.Request, resp any) (http.Header, error) {

}

func (c *Client) NewRequest(method string, url string, headers map[string]string, params any, body any) (*http.Request, error) {

}
