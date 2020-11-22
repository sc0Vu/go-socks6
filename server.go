package socks6

import (
	"fmt"
	"net"
	"sync"
	"time"
	// "io"
	// "encoding/binary"
)

var (
	ErrCopyEmptyBuffer = fmt.Errorf("copy empty buffer")
	ErrClosedListener  = fmt.Errorf("closed listener")
)

// Server is socks 6 server
type Server struct {
	dialTimeout time.Duration
	sleepTime   time.Duration
	ln          net.Listener
	od          sync.Once
	ch          chan struct{}
}

func NewSocksServer(dialTimeout, sleepTime time.Duration) Server {
	return Server{
		dialTimeout: dialTimeout,
		sleepTime:   sleepTime,
		od:          sync.Once{},
		ch:          make(chan struct{}),
	}
}

func (srv *Server) Closed() bool {
	select {
	case <-srv.ch:
		return true
	default:
		return false
	}
}

// handshake
func (srv *Server) handshake(conn net.Conn) (err error) {
	return
}

// Serve net.Listener
func (srv *Server) Serve(ln net.Listener) error {
	if srv.Closed() {
		return fmt.Errorf("socks server is closed")
	}
	srv.ln = ln
	for {
		conn, err := ln.Accept()
		if err != nil {
			if srv.Closed() {
				err = ErrClosedListener
				return err
			}
			if ne, ok := err.(net.Error); ok && ne.Temporary() {
				// sleep for a while
				time.Sleep(srv.sleepTime)
				continue
			}
		}
		if err := srv.handshake(conn); err == nil {
			// dial the other connection and copy each other
		}
	}
}

// ListenAndServe listen and serve socks server to hostname
// hostname should be [ip]:[port]
func (srv *Server) ListenAndServe(hostname string) error {
	if srv.Closed() {
		return fmt.Errorf("socks server is closed")
	}
	ln, err := net.Listen("tcp", hostname)
	if err != nil {
		return err
	}
	return srv.Serve(ln)
}

// Close socks server
func (srv *Server) Close() {
	srv.od.Do(func() {
		srv.ln.Close()
		close(srv.ch)
	})
}
