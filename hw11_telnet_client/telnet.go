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
	data, err := client.in.ReadString('\n')
	if err == nil {
		_, err = client.conn.Write([]byte(data))
	}
	return err
}

func (client *TCPClient) Receive() error {
	data := make([]byte, 1024)
	n, err := client.conn.Read(data)
	if err == nil && n > 0 {
		_, err = client.out.Write(data[:n])
	}
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
