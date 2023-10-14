package server

import (
	"net"
	"sync"
)

type registry struct {
	m     sync.Mutex
	conns map[string]net.Conn
}

func newRegistry() *registry {
	return &registry{conns: make(map[string]net.Conn)}
}

func (r *registry) add(k string, conn net.Conn) {
	r.m.Lock()
	defer r.m.Unlock()

	r.conns[k] = conn
}

func (r *registry) get(k string) (conn net.Conn, ok bool) {
	r.m.Lock()
	defer r.m.Unlock()

	conn, ok = r.conns[k]

	if ok {
		delete(r.conns, k)
	}

	return
}

func (r *registry) cleanup(k string) error {
	r.m.Lock()
	defer r.m.Unlock()

	conn, ok := r.conns[k]
	if !ok {
		return nil
	}

	delete(r.conns, k)

	return conn.Close()
}
