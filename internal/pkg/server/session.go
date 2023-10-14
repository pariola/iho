package server

import (
	"net"
	"sync/atomic"
	"time"

	"github.com/google/uuid"

	"iho/internal/pkg/messages"
	"iho/internal/pkg/tunnel"
)

const (
	pingInterval         = 5 * time.Second
	maxPingTimeout       = 5 * time.Second
	tunnelConnectionWait = 3 * time.Second
)

type session struct {
	ctrl net.Conn

	listener net.Listener

	registry *registry

	lastPingAt atomic.Int64
}

func newSession(conn net.Conn, listener net.Listener, registry *registry) *session {
	return &session{
		ctrl:     conn,
		listener: listener,
		registry: registry,
	}
}

func (s *session) close() {
	s.listener.Close()
	s.ctrl.Close()
}

func (s *session) controlLoop() {
	for {
		msg, err := tunnel.Receive(s.ctrl)
		if err != nil {
			return
		}

		switch msg.(type) {
		case *messages.Ping:
			s.lastPingAt.Store(time.Now().UnixNano())
		}
	}
}

func (s *session) heartbeat() {
	ping := time.NewTicker(pingInterval)
	defer ping.Stop()

	for range ping.C {
		lastPing := time.Unix(0, s.lastPingAt.Load())
		if time.Since(lastPing) > maxPingTimeout {
			s.close() // close session
			return
		}
	}
}

func (s *session) serve() {
	go s.heartbeat()
	go s.controlLoop()

	var tempDelay time.Duration // how long to sleep on accept failure

	for {
		conn, err := s.listener.Accept()
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
			return
		}

		id := uuid.New().String()
		s.registry.add(id, conn)

		err = tunnel.Send(s.ctrl, &messages.RequestConnection{Id: id})
		if err != nil {
			s.registry.cleanup(id) // remove connection
			return
		}

		// cleanup if not used
		time.AfterFunc(tunnelConnectionWait, func() { s.registry.cleanup(id) })
	}
}
