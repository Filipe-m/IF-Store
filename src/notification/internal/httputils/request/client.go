package request

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

type Client interface {
	Get(ctx context.Context, path string, opt ...Option) error
	Post(ctx context.Context, path string, opt ...Option) error
	Put(ctx context.Context, path string, opt ...Option) error
	Delete(ctx context.Context, path string, opt ...Option) error
}

type request struct {
	client *http.Client
	url    string
}

type ConfigOption func(c *request)

func WithClient(client *http.Client) ConfigOption {
	return func(c *request) {
		c.client = client
	}
}

func New(url string, opt ...ConfigOption) Client {
	client := request{
		client: http.DefaultClient,
		url:    url,
	}
	for _, o := range opt {
		o(&client)
	}
	return &client
}

type Option func(r *entity)

type entity struct {
	req interface{}
	res interface{}
}

func WithRequest(req interface{}) Option {
	return func(r *entity) {
		r.req = req
	}
}

func WithResponse(res interface{}) Option {
	return func(r *entity) {
		r.res = res
	}
}

func (r *request) Get(ctx context.Context, path string, opt ...Option) error {
	return r.send(ctx, http.MethodGet, path, opt...)
}

func (r *request) Post(ctx context.Context, path string, opt ...Option) error {
	return r.send(ctx, http.MethodPost, path, opt...)
}

func (r *request) Put(ctx context.Context, path string, opt ...Option) error {
	return r.send(ctx, http.MethodPut, path, opt...)
}

func (r *request) Delete(ctx context.Context, path string, opt ...Option) error {
	return r.send(ctx, http.MethodDelete, path, opt...)
}

func (r *request) send(ctx context.Context, method, path string, opt ...Option) error {
	var ety entity
	for _, o := range opt {
		o(&ety)
	}

	var buf bytes.Buffer

	if ety.req != nil {
		err := json.NewEncoder(&buf).Encode(ety.req)
		if err != nil {
			return err
		}
	}

	hreq, err := http.NewRequestWithContext(ctx, method, fmt.Sprintf("%s%s", r.url, path), &buf)
	if err != nil {
		return err
	}

	hreq.Header.Set("Content-Type", "application/json")

	result, err := r.client.Do(hreq)
	if err != nil {
		return err
	}
	defer result.Body.Close()

	if result.StatusCode < 200 || result.StatusCode >= 300 {
		return errors.New("error when request")
	}

	if ety.res != nil {
		err = json.NewDecoder(result.Body).Decode(&ety.res)
		if err != nil {
			return err
		}
		return nil
	}

	return nil
}
