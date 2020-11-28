package socks6

import (
	"testing"
	"time"
)

func TestListenAndServer(t *testing.T) {
	hostname := "localhost:8888"
	s := NewSocksServer(1*time.Millisecond, 1*time.Millisecond, 1024)
	if err := s.ListenAndServe(hostname); err != nil {
		t.Fatal(err)
	}
}
