package main

import "fmt"

type HttpClient struct {
	Timeout int
	MaxIdle int
}

type ClientOption func(client *HttpClient)
type ClientOptions []ClientOption

func (c ClientOptions) apply(h *HttpClient) {
	for _, opt := range c {
		opt(h)
	}
}
func NewHttpClient(opts ...ClientOption) *HttpClient {
	c := &HttpClient{}
	ClientOptions(opts).apply(c)
	return c
}
func WithTimeOut(timeout int) ClientOption {
	return func(client *HttpClient) {
		client.Timeout = timeout
	}
}

func main() {
	a := NewHttpClient(WithTimeOut(100))
	fmt.Println(a)
}
