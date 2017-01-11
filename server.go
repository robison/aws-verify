package main

import "log"
import "net"
import "net/http"
import "os"
import "os/signal"
import "syscall"
import "time"

// CreateServer creates a new instance of an HTTP server
func CreateServer(socket string, mode os.FileMode, handler http.Handler) *Server {
	return &Server{
		server: &http.Server{
			Handler:      handler,
			ReadTimeout:  time.Second,
			WriteTimeout: time.Second,
		},

		socket: socket,
		mode: mode,
	}
}

// Server configuration and handle for an HTTP server
type Server struct {
	listener net.Listener
	server   *http.Server
	shutdown chan os.Signal

	socket   string
	mode		os.FileMode
}

// Listen creates a UNIX listener and attaches an HTTP server to it
func (s *Server) Listen() error {
	listener, err := net.Listen("unix", s.socket)
	if err != nil {
		return err
	}

	s.listener = listener

	err = os.Chmod(s.socket, s.mode)
	if err != nil {
		s.Close()
		return err
	}

	s.shutdown = make(chan os.Signal, 1)

	// Set up a deferred cleanup routine
	signal.Notify(s.shutdown, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)

	go func() {
		sig := <-s.shutdown
		log.Printf("Received signal '%s', shutting down.", sig)

		s.Close()
	}()

	log.Printf("Listen for requests on unix:%s", s.socket)
	s.server.Serve(listener)

	log.Printf("Goodbye")

	return nil
}

// Close closes the UNIX listener and cleans up it's socket handle
func (s *Server) Close() {
	log.Printf("Closing server and cleaning up socket %s", s.socket)

	s.listener.Close()
	os.Remove(s.socket)
}
