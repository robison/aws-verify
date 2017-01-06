package main

import "log"
import "net"
import "net/http"
import "os"
import "os/signal"
import "syscall"
import "time"

/**
 * Create a new instance of an HTTP server
 */
func CreateServer(socket string, handler http.Handler) *Server {
	return &Server{
		server: &http.Server{
			Handler:      handler,
			ReadTimeout:  time.Second,
			WriteTimeout: time.Second,
		},
		socket: socket,
	}
}

/**
 * Configuration and handles for an HTTP server
 */
type Server struct {
	listener net.Listener
	server   *http.Server
	socket   string
	shutdown chan os.Signal
}

/**
 * Create listener and attach HTTP server to it
 */
func (s *Server) Listen() {
	listener, err := net.Listen("unix", s.socket)
	fatal(err)

	s.listener = listener
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
}

/**
 * Close the listener and clean up it's socket handle
 */
func (s *Server) Close() {
	log.Printf("Closing server and cleaning up socket %s", s.socket)

	s.listener.Close()
	os.Remove(s.socket)
}
