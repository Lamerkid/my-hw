package main

import (
	"io"
	"net"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &Client{
		Address: address,
		Timeout: timeout,
		In:      in,
		Out:     out,
	}
}

type Client struct {
	Address    string
	Timeout    time.Duration
	In         io.ReadCloser
	Out        io.Writer
	connection net.Conn
}

func (c *Client) Connect() error {
	connection, err := net.DialTimeout("tcp", c.Address, c.Timeout)
	if err != nil {
		return err
	}
	c.connection = connection
	return nil
}

func (c *Client) Send() error {
	_, err := io.Copy(c.connection, c.In)
	return err
}

func (c *Client) Receive() error {
	_, err := io.Copy(c.Out, c.connection)
	return err
}

func (c *Client) Close() error {
	return c.connection.Close()
}
