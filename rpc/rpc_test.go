package rpc_test

import (
	"cdecl-lsp/rpc"
	"testing"
)

type EncodingExample struct {
	Foo string
}

var example = EncodingExample{
	Foo: "bar",
}

func TestEncodeMessage(t *testing.T) {
	expected := "Content-Length: 13\r\n\r\n{\"Foo\":\"bar\"}"
	actual := rpc.EncodeMessage(example)
	if actual != expected {
		t.Fatalf("expected %s, got %s", expected, actual)
	}
}

func TestDecodeMessage(t *testing.T) {
	msg := "Content-Length: 16\r\n\r\n{\"method\":\"bar\"}"
	method, content, err := rpc.DecodeMessage([]byte(msg))
	if err != nil {
		t.Fatal(err)
	}

	if len(content) != 16 {
		t.Fatalf("expected %d, got %d", 16, len(content))
	}

	if method != "bar" {
		t.Fatalf("expected %s, got %s", "bar", method)
	}
}
