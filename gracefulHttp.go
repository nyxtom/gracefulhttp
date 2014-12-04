package gracefulhttp

import (
	"net"
	"net/http"
	"os"
	"syscall"
	"time"
)

type Server struct {
	http.Server
	listener       *gracefulListener
	FileDescriptor int
}

func ListenAndServe(addr string, fd int) error {
	server := NewServer(addr, fd)
	defer server.Close()
	return server.ListenAndServe()
}

func NewServer(addr string, fd int) *Server {
	server := new(Server)
	server.Addr = addr
	server.ReadTimeout = 10 * time.Second
	server.WriteTimeout = 10 * time.Second
	server.MaxHeaderBytes = 1 << 16
	server.FileDescriptor = fd
	return server
}

func (server *Server) File() *os.File {
	return server.listener.File()
}

func (server *Server) Fd() uintptr {
	fl := server.File()
	fd := fl.Fd()
	noCloseOnExec(fd)
	return fd
}

func (server *Server) ListenAndServe() error {
	var err error
	var l net.Listener

	fd := server.FileDescriptor
	addr := server.Addr
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
	err = server.Serve(server.listener)
	return err
}

func (server *Server) Close() error {
	var err error
	if server.listener != nil {
		err = server.listener.Close()
	}
	return err
}
