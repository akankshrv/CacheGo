package client

import (
	"context"
	"fmt"

	"net"

	"github.com/akanshrv/CacheGo/proto"
)

type Options struct{}

type Client struct {
	conn net.Conn
}

func NewFromConn(conn net.Conn) *Client {
	return &Client{
		conn: conn,
	}
}

func New(endpoint string, opts Options) (*Client, error) {
	conn, err := net.Dial("tcp", endpoint)
	if err != nil {
		return nil, err
	}

	return &Client{
		conn: conn,
	}, nil
}
func (c *Client) Get(ctx context.Context, key []byte) ([]byte, error) {
	cmd := &proto.CommandGet{
		Key: key,
	}
	_, err := c.conn.Write(cmd.Bytes())
	if err != nil {
		return nil, fmt.Errorf("error writing command: %w", err)
	}

	resp, err := proto.ParseGetResponse(c.conn)

	if err != nil {
		return nil, fmt.Errorf("error parsing response: %w", err)
	}

	switch resp.Status {
	case proto.StatusOK:
		return resp.Value, nil
	case proto.StatusKeyNotFound:
		return nil, fmt.Errorf("key not found (%s)", key) // Key not found, return nil, nil
	default:
		return nil, fmt.Errorf("server responded with non-OK status [%s]", resp.Status)
	}

}
func (c *Client) Set(ctx context.Context, key []byte, value []byte, ttl int) error {
	cmd := &proto.CommandSet{
		Key:   key,
		Value: value,
		TTL:   ttl,
	}
	_, err := c.conn.Write(cmd.Bytes())
	if err != nil {
		return err
	}

	resp, err := proto.ParseSetResponse(c.conn)
	if err != nil {
		return err
	}

	if resp.Status != proto.StatusOK {
		return fmt.Errorf("server responsed with non OK status [%s]", resp.Status)
	}
	return nil
}
func (c *Client) Close() error {

	return c.conn.Close()
}
