package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"time"
)

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

type simpleClient struct {
	conn    net.Conn
	address string
	timeout time.Duration
	in      io.ReadCloser
	out     io.Writer
}

func (c *simpleClient) Connect() error {
	conn, err := net.DialTimeout("tcp", c.address, c.timeout)
	if err != nil {
		return err
	}
	c.conn = conn

	return nil
}

func (c *simpleClient) Close() error {
	return c.conn.Close()
}

func (c *simpleClient) Send() error {
	if _, err := io.Copy(c.conn, c.in); err != nil {
		return fmt.Errorf("error sending: %w", err)
	}

	byteLine := []byte("...EOF\n")
	n, err := os.Stderr.Write(byteLine)
	if err != nil {
		return err
	}
	if n != len(byteLine) {
		return fmt.Errorf("the app could not send ...EOF to Stderr")
	}
	return nil
}

func (c *simpleClient) Receive() error {
	if _, err := io.Copy(c.out, c.conn); err != nil {
		return fmt.Errorf("error reading: %w", err)
	}

	byteLine := []byte("...Connection was closed by peer\n")
	n, err := os.Stderr.Write(byteLine)
	if err != nil {
		return err
	}
	if n != len(byteLine) {
		return fmt.Errorf("the app could not send ...EOF to Stderr")
	}
	return nil
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &simpleClient{
		in:      in,
		out:     out,
		address: address,
		timeout: timeout,
	}
}
