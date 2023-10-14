package messages

import (
	"errors"

	"github.com/vmihailenco/msgpack/v5"
)

var ErrUnknownType = errors.New("unknown type")

type Type int

const (
	typeUnknown Type = iota
	TypeError
	TypePing
	TypeOpenTunnel
	TypeTunnelOpened
	TypeConnectSession
	TypeRequestConnection
)

type Message interface {
	T() Type
}

type msgWrapper struct {
	Type Type
	Data msgpack.RawMessage
}

func Pack(m Message) ([]byte, error) {
	raw, err := msgpack.Marshal(m)
	if err != nil {
		return nil, err
	}

	msg := msgWrapper{
		Type: m.T(),
		Data: raw,
	}

	return msgpack.Marshal(msg)
}

func Unpack(data []byte) (Message, error) {
	var m msgWrapper

	err := msgpack.Unmarshal(data, &m)
	if err != nil {
		return nil, err
	}

	var msg Message

	switch m.Type {
	default:
		return nil, ErrUnknownType

	case TypePing:
		msg = &Ping{}
	case TypeError:
		msg = &Error{}
	case TypeOpenTunnel:
		msg = &OpenTunnel{}
	case TypeTunnelOpened:
		msg = &TunnelOpened{}
	case TypeConnectSession:
		msg = &ConnectSession{}
	case TypeRequestConnection:
		msg = &RequestConnection{}
	}

	err = msgpack.Unmarshal(m.Data, msg)
	if err != nil {
		return nil, err
	}

	return msg, nil
}
