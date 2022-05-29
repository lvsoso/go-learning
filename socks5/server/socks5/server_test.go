package socks5

import (
	"bytes"
	"reflect"
	"testing"
)

func TestAuth(t *testing.T) {

	t.Run("a valid client auth message", func(t *testing.T) {
		var buf bytes.Buffer
		buf.Write([]byte{SOCKS5Version, 2, MethodNoAuth, MethodGSSAPI})
		if err := auth(&buf); err != nil {
			t.Fatalf("should get error nil but got %s", err)
		}

		want := []byte{SOCKS5Version, MethodNoAuth}
		got := buf.Bytes()
		if !reflect.DeepEqual(want, got) {
			t.Fatalf("hsould get message %v but got %v", want, got)
		}
	})

	t.Run("an invalid client auth message", func(t *testing.T) {
		var buf bytes.Buffer
		buf.Write([]byte{SOCKS5Version, 2, MethodNoAuth})
		if err := auth(&buf); err == nil {
			t.Fatalf("should get error EOF but got nil")
		}
	})
}
