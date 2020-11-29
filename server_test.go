package socks6

import (
	"bytes"
	"net"
	"testing"
	"time"
)

func testWriteTCP(t *testing.T, conn net.Conn, buf []byte) {
	n, err := conn.Write(buf)
	if err != nil {
		t.Fatal(err)
	}
	if n != len(buf) {
		t.Fatalf("data was truncated")
	}
}

func testReadTCP(t *testing.T, conn net.Conn, buf []byte) {
	n, err := conn.Read(buf)
	if err != nil {
		t.Fatal(err)
	}
	if n != len(buf) {
		t.Fatalf("data was truncated")
	}
}

func TestSocks4Server(t *testing.T) {
	hostname := "localhost:8888"
	s := NewSocksServer(1*time.Millisecond, 1*time.Millisecond, 10)
	go func() {
		if err := s.ListenAndServe(hostname); err != nil {
			// t.Fatal(err)
		}
	}()
	conn, err := net.Dial("tcp", hostname)
	if err != nil {
		t.Fatal(err)
	}
	// connect to google.com
	testWriteTCP(t, conn, []byte{4, 1, 0, 80, 216, 58, 200, 46})

	res := make([]byte, 8)
	testReadTCP(t, conn, res)
	if !bytes.Equal(res, []byte{0, 90, 0, 80, 216, 58, 200, 46}) {
		t.Fatalf("failed to connect")
	}
	conn.Close()
	s.Close()
}
