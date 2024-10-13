package netconn

import (
	"bufio"
	"fmt"
	"net"
	"net/textproto"
)

type Connector interface {
	Connect() error
	Read() ([]byte, error)
	Write([]byte) error
	Close() error
}

const (
	lineSep = '\n'
)

type Provider struct {
	network string
	addr    string
	conn    net.Conn
}

func NewWithAddr(network, addr string) *Provider {
	return &Provider{
		network: network,
		addr:    addr,
	}
}

func NewWithConn(conn net.Conn) *Provider {
	return &Provider{
		network: conn.RemoteAddr().Network(),
		addr:    conn.RemoteAddr().String(),
		conn:    conn,
	}
}

func (p *Provider) Connect() error {
	conn, err := net.Dial(p.network, p.addr)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %w", err)
	}
	p.conn = conn
	return nil
}

func (p *Provider) Read() ([]byte, error) {
	reader := bufio.NewReader(p.conn)
	tp := textproto.NewReader(reader)
	bytes, err := tp.ReadLineBytes()
	if err != nil {
		return nil, fmt.Errorf("failed to read from connection: %w", err)
	}
	return bytes, nil
}

func (p *Provider) Write(mess []byte) error {
	message := append(mess, lineSep)
	_, err := p.conn.Write(message)
	if err != nil {
		return fmt.Errorf("failed to write to connection: %w", err)
	}
	return nil
}

func (p *Provider) Close() error {
	if p.conn == nil {
		return nil
	}
	return p.conn.Close()
}
