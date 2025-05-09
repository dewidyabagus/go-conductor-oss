package http

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"sync"
)

var (
	once   sync.Once
	client *http.Client
)

type Response struct {
	status     string // e.g. "200 OK"
	statusCode int    // e.g. 200
	body       *bytes.Buffer
}

func (r *Response) StatusCode() int {
	return r.statusCode
}

func (r *Response) BodyBytes() []byte {
	return r.body.Bytes()
}

type httpClient struct {
	baseUrl string
}

type HttpClient interface {
	Get(ctx context.Context, path string, headers map[string]string) (*Response, error)
}

func NewHttpClient(baseURL string) HttpClient {
	once.Do(func() {
		client = http.DefaultClient
	})
	return &httpClient{baseURL}
}

func (c *httpClient) Get(ctx context.Context, path string, headers map[string]string) (*Response, error) {
	return c.do(ctx, http.MethodGet, c.baseUrl+path, headers, nil)
}

func (c *httpClient) do(ctx context.Context, method, url string, headers map[string]string, body any) (*Response, error) {
	raw, err := json.Marshal(body)
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequestWithContext(ctx, method, url, bytes.NewBuffer(raw))
	if err != nil {
		return nil, err
	}
	for k, v := range headers {
		req.Header.Set(k, v)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	result := &Response{
		status:     resp.Status,
		statusCode: resp.StatusCode,
		body:       new(bytes.Buffer),
	}
	if _, err = io.Copy(result.body, resp.Body); err != nil {
		return nil, err
	}
	return result, nil
}
