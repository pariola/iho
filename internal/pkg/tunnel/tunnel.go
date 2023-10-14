package tunnel

import (
	"io"
	"net"

	"iho/internal/pkg/messages"
)

const frameSize = 128

func Send(conn net.Conn, msg messages.Message) error {
	data, err := messages.Pack(msg)
	if err != nil {
		return err
	}

	// pad message
	padded := make([]byte, frameSize)
	_ = copy(padded, data)

	_, err = conn.Write(padded)
	if err != nil {
		return err
	}

	return nil
}

func Receive(conn net.Conn) (messages.Message, error) {
	data := make([]byte, frameSize)

	_, err := conn.Read(data)
	if err != nil {
		return nil, err
	}

	msg, err := messages.Unpack(data)
	if err != nil {
		return nil, err
	}

	return msg, nil
}

func Proxy(s1, s2 io.ReadWriter) {
	done := make(chan struct{}, 1)

	go pipe(s1, s2, done)
	go pipe(s2, s1, done)

	<-done
}

func pipe(dst io.Writer, src io.Reader, done chan struct{}) {
	io.Copy(dst, src)
	done <- struct{}{}
}
