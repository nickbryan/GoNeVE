package network

import (
	"fmt"

	"google.golang.org/grpc"
)

// Handler is a hook into the grpc.ClientConn. A consumer who wishes to send information from the Client
// to the Server should register a Handler interact with the network.
type Handler interface {
	Handle(*grpc.ClientConn)
}

// HandlerFunc provides a convenience wrapper to allow a function to be used as a Handler.
type HandlerFunc func(*grpc.ClientConn)

func (hf HandlerFunc) Handle(cc *grpc.ClientConn) {
	hf(cc)
}

type Client struct {
	Addr string

	cc *grpc.ClientConn
}

func NewClient(addr string) *Client {
	return &Client{Addr: addr}
}

func (c *Client) Dial(opts ...grpc.DialOption) error {
	// TODO: grpc.DialContext ?? use sigs to detect cmd+c and cancel context?
	conn, err := grpc.Dial(c.Addr, append(opts, grpc.WithInsecure())...)
	if err != nil {
		return fmt.Errorf("unable to Dial server: %v", err)
	}

	c.cc = conn

	return nil
}

func (c *Client) Handle(h Handler) {
	h.Handle(c.cc)
}
