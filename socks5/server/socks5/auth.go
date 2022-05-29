package socks5

import (
	"errors"
	"io"
)

type Method = byte

const (
	MethodNoAuth       Method = 0x00
	MethodGSSAPI       Method = 0x01
	MethodPassword     Method = 0x02
	MethodNoAcceptable Method = 0xff
)

type ClientAuthMessage struct {
	Version  byte
	NMethods byte
	Methods  []Method
}

func NewClientAuthMessage(conn io.Reader) (*ClientAuthMessage, error) {
	// Read version, NMethods
	buf := make([]byte, 2)
	_, err := io.ReadFull(conn, buf)
	if err != nil {
		return nil, err
	}

	// Validate version
	if buf[0] != SOCKS5Version {
		return nil, errors.New("Protocol version not supported")
	}

	// Read Methods
	nmethods := buf[1]
	buf = make([]byte, nmethods)
	_, err = io.ReadFull(conn, buf)
	if err != nil {
		return nil, err
	}

	return &ClientAuthMessage{
		Version:  SOCKS5Version,
		NMethods: nmethods,
		Methods:  buf,
	}, nil
}

func NewServerAuthMessage(conn io.Writer, method Method) error {
	buf := []byte{SOCKS5Version, method}
	_, err := conn.Write(buf)
	return err
}
