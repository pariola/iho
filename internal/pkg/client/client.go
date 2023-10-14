package client

import (
	"io"
	"log"
	"net"
	"time"

	"iho/internal/pkg/messages"
	"iho/internal/pkg/tunnel"
)

const pingInterval = 3 * time.Second

type client struct {
	toAddr, remoteAddr string
}

func logf(format string, v ...any) {
	log.Printf(format, v...)
}

func Connect(toAddr, remoteAddr string) error {
	conn, err := net.Dial("tcp", remoteAddr)
	if err != nil {
		return err
	}
	defer conn.Close()

	err = tunnel.Send(conn, &messages.OpenTunnel{})
	if err != nil {
		return err
	}

	msg, err := tunnel.Receive(conn)
	if err != nil {
		return err
	}

	c := &client{
		toAddr:     toAddr,
		remoteAddr: remoteAddr,
	}

	switch m := msg.(type) {
	default:

	case *messages.Error:
		logf("Server Error: %s", m.Message)

	case *messages.TunnelOpened:
		logf("Tunnel listening on %s", m.Port)
		c.handleTunnel(conn)
		logf("Closing tunnel")
	}

	return nil
}

func (c *client) heartbeat(conn net.Conn) {
	t := time.NewTicker(pingInterval)
	defer t.Stop()

	logf("ping message every %v", pingInterval)

	for range t.C {
		err := tunnel.Send(conn, &messages.Ping{})
		if err != nil {
			logf("ping failed: %+v", err)
			return
		}
	}
}

func (c *client) handleTunnel(conn net.Conn) {
	go c.heartbeat(conn)

	for {
		msg, err := tunnel.Receive(conn)
		if err != nil {
			if err != io.EOF {
				logf("read from remote failed, err: %v", err)
			}
			return
		}

		switch m := msg.(type) {
		default:

		case *messages.Error:
			logf("Server Error: %s", m.Message)
			return

		case *messages.RequestConnection:
			go c.handleConnection(m.Id)
		}
	}

}

func (c *client) handleConnection(id string) error {
	remoteConn, err := net.Dial("tcp", c.remoteAddr)
	if err != nil {
		return err
	}
	defer remoteConn.Close()

	err = tunnel.Send(remoteConn, &messages.ConnectSession{Id: id})
	if err != nil {
		return err
	}

	toConn, err := net.Dial("tcp", c.toAddr)
	if err != nil {
		logf("Error dailing local addr: %+v", err)
		return err
	}
	defer toConn.Close()

	tunnel.Proxy(toConn, remoteConn)
	return nil
}
