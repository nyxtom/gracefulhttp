package main

import (
	"net"
	"net/http"
	"os"
	"syscall"
	"time"
)

type Server struct {
	http.Server
	listener *gracefulListener
}

func ListenAndServe(addr string, fd int) {
	server := NewServer(addr)
	defer server.Close()
	server.ListenAndServe(fd)
}

func NewServer(addr string) *Server {
	server := &Server{
		Addr:           addr,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 16,
	}
	return server
}

func (server *Server) ListenAndServe(fd int) {
	var err error
	var l net.Listener

	if fd != 0 {
		// listen on the existing file descriptor
		f := os.NewFile(uintptr(fd), "listen socket")
		l, err = net.FileListener(f)
	} else {
		l, err = net.Listen("tcp", addr)
	}

	if err != nil {
		return err
	}

	if fd != 0 {
		parent := syscall.Getppid()
		syscall.Kill(parent, syscall.SIGTERM)
	}

	server.listener = newGracefulListener(l)
	server.Serve(server.listener)
}

func (server *Server) Close() error {
	return server.listener.Close()
}
