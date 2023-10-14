package server

import (
	"net"
	"time"

	"iho/internal/pkg/messages"
	"iho/internal/pkg/tunnel"
)

type handler struct {
	addr     string
	registry *registry
}

func ListenAndServe(addr string) error {
	h := handler{
		addr:     addr,
		registry: newRegistry(),
	}

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	defer l.Close()

	// go func() {
	// 	for range time.NewTicker(2 * time.Second).C {
	// 		logf("Goroutines: %d", runtime.NumGoroutine())
	// 		// logf("Connections: %d", h.sessions.len())
	// 	}
	// }()

	var tempDelay time.Duration // how long to sleep on accept failure

	for {
		conn, err := l.Accept()
		if err != nil {
			if ne, ok := err.(net.Error); ok && ne.Timeout() {
				if tempDelay == 0 {
					tempDelay = 5 * time.Millisecond
				} else {
					tempDelay *= 2
				}
				if max := 1 * time.Second; tempDelay > max {
					tempDelay = max
				}
				logf("tcp: 'accept' error: %v; retrying in %v", err, tempDelay)
				time.Sleep(tempDelay)
				continue
			}
			return err
		}

		go h.handleConn(conn)
	}
}

func (h handler) handleConn(conn net.Conn) {
	defer conn.Close()

	msg, err := tunnel.Receive(conn)
	if err != nil {
		logf("failed to decode data %v", err)
		return
	}

	switch m := msg.(type) {
	case *messages.OpenTunnel:
		h.startSession(conn)

	case *messages.ConnectSession:
		h.tunnelSession(conn, m.Id)

	default:
		logf("client says something unknown")
	}
}

func (h handler) startSession(conn net.Conn) {
	l, err := net.Listen("tcp", "")
	if err != nil {
		tunnel.Send(conn, &messages.Error{Message: err.Error()})
		return
	}

	defer l.Close()

	tunnel.Send(conn, &messages.TunnelOpened{Port: l.Addr().String()})

	newSession(conn, l, h.registry).serve()
}

func (h handler) tunnelSession(conn net.Conn, connId string) {
	c, ok := h.registry.get(connId)
	if !ok {
		logf("unknown connection")
		return
	}
	tunnel.Proxy(c, conn)
}
