package http

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"sync"
)

type Option func(*Client)

type Client struct {
	header map[string][]string
	mutex  sync.Mutex
	host   string
	token  string
	client *http.Client
}

type ClientConn interface {
	ClientInterface
	Info()
}

func WithHeader(k string, vs ...string) Option {
	return func(c *Client) {
		c.mutex.Lock()
		defer c.mutex.Unlock()
		c.header[k] = append(c.header[k], vs...)
	}
}

func WithHost(h string) Option {
	return func(c *Client) {
		c.host = h
	}
}

func WithToke(t string) Option {
	return func(c *Client) {
		c.token = t
	}
}

func NewClient(opts ...Option) ClientConn {
	client := &Client{
		header: make(map[string][]string),
		mutex:  sync.Mutex{},
		client: http.DefaultClient,
	}
	for _, o := range opts {
		o(client)
	}
	return client
}

type ParamsFunc func() (io.Reader, error)
type ResponseFunc func(resp *http.Response) error
type Response struct {
	Code int32       `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

var DefaultParamsFunc = func() (io.Reader, error) {
	return nil, nil
}

func RespFunc[T any](res T) ResponseFunc {
	return func(resp *http.Response) error {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			return err
		}
		if err = json.Unmarshal(body, res); err != nil {
			fmt.Println(err.Error())
		}
		return nil
	}
}

type ClientInterface interface {
	Request(method, path string, paramFn ParamsFunc, respFn ResponseFunc) error
}

func (obj *Client) Request(method, path string, paramFn ParamsFunc, respFn ResponseFunc) error {
	if len(obj.host) == 0 || len(path) == 0 || len(method) == 0 {
		return errors.New("invalid args")
	}
	reader, err := paramFn()
	if err != nil {
		return err
	}
	req, err := http.NewRequest(method, fmt.Sprintf("%s/%s", obj.host, path), reader)
	if err != nil {
		return err
	}
	req.Header.Set("Content-Type", "application/json")
	for k, vs := range obj.header {
		for _, v := range vs {
			req.Header.Set(k, v)
		}
	}
	//内部直接设置 还是通过ctx进行处理
	resp, err := obj.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	return respFn(resp)
}

func (obj *Client) Info() {
	slog.Info("client info", "host", obj.host, "token", obj.token)
}
