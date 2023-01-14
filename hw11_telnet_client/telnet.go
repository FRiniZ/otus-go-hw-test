package main

import (
	"bufio"
	"io"
	"net"
	"time"
)

type TCPClient struct {
	address string
	timeout time.Duration
	conn    net.Conn
	in      *bufio.Reader
	out     io.Writer
}

func (client *TCPClient) Connect() error {
	var err error
	client.conn, err = net.DialTimeout("tcp", client.address, client.timeout)
	return err
}

func (client *TCPClient) Send() error {
	_, err := io.Copy(client.conn, client.in)
	return err
}

func (client *TCPClient) Receive() error {
	_, err := io.Copy(client.out, client.conn)
	return err
}

func (client *TCPClient) Close() error {
	return client.conn.Close()
}

type TelnetClient interface {
	Connect() error
	io.Closer
	Send() error
	Receive() error
}

func NewTelnetClient(address string, timeout time.Duration, in io.ReadCloser, out io.Writer) TelnetClient {
	return &TCPClient{address: address, timeout: timeout, in: bufio.NewReader(in), out: out}
}
