package socks5

import (
	"errors"
)

var (
	ErrVersionNotSupported       = errors.New("protocol version not supported")
	ErrMethodVersionNotSupported = errors.New("sub-negotiation method version not supported")
	ErrCommandNotSupported       = errors.New("requst command not supported")
	ErrInvalidReservedField      = errors.New("invalid reserved field")
	ErrAddressTypeNotSupported   = errors.New("address type not supported")

	ErrPasswordCheckerNotSet = errors.New("error password checker not set")
	ErrPasswordAuthFailure   = errors.New("error authenticating username/password")
)
